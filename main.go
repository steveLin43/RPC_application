package main

import (
	"RPC_application/server"
	"context"
	"encoding/json"
	"flag"
	"log"
	"net"
	"net/http"
	"strings"

	pb "RPC_application/proto"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/internal/status"
	"google.golang.org/grpc/reflection"
)

var port string

func init() {
	flag.StringVar(&port, "port", "8003", "啟動通訊埠編號")
	flag.Parse()
}

func main() {
	err := RunServer(port)
	if err != nil {
		log.Fatalf("Run Server err: %v", err)
	}
}

func RunServer(port string) error {
	httpMux := RunHttpServer()
	grpcS := RunGrpcServer()
	gatewayMux := runGrpcGatewayServer()

	httpMux.Handler("/", gatewayMux)
	return http.ListenAndServe(":"+port, grpcHandleFunc(grpcS, httpMux))
}

func RunTCPServer(port string) (net.Listener, error) {
	return net.Listen("tcp", ":"+port)
}

func RunHttpServer() *http.ServeMux { //針對http1.0
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/ping",
		func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`pong`))
		},
	)

	return serveMux
}

func RunGrpcServer() *grpc.Server { //針對grpc
	s := grpc.NewServer()
	pb.RegisterTagServiceServer(s, server.NewTagServer())
	reflection.Register(s)

	return s
}

func grpcHandleFunc(grpcServer *grpc.Server, otherHandle http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandle.ServeHTTP(w, r)
		}
	}), &http2.Server{})
}

func runGrpcGatewayServer() *runtime.ServeMux {
	endpoint := "0.0.0.0:" + port
	runtime.HTTPError = grpcGatewayError
	gwmux := runtime.NewServeMux()
	dopts := []grpc.DialOption{grpc.WithInsecure()}
	_ = pb.RegisterTagServiceHandlerFromEndpoint(context.Background(), gwmux, endpoint, dopts)
	return gwmux
}

type httpError struct {
	Code    int32  `json:"code,omitempty`
	Message string `json:"message,omitempty`
}

func grpcGatewayError(ctx context.Context, _ *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, _ *http.Request, err error) {
	s, ok := status.FromError(err)
	if !ok {
		s = status.New(codes.Unknown, err.Error())
	}

	httpError := httpError{Code: int32(s.Code()), Message: s.Message()}
	details := s.Details()
	for _, detail := range details {
		if v, ok := detail.(*pb.Error); ok {
			httpError.Code = v.Code
			httpError.Message = v.Message
		}
	}

	resp, _ := json.Marshal(httpError)
	w.Header().Set("Content-type", marshaler.ContentType())
	w.WriteHeader(runtime.HTTPStatusFromCode(s.Code()))
	_, _ = w.Write(resp)
}
