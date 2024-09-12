package main

import (
	"RPC_application/server"
	"flag"
	"log"
	"net"
	"net/http"

	pb "RPC_application/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var grpcPort string //對接grpc、http2.0
var httpPort string //對接http1.0

func init() {
	flag.StringVar(&grpcPort, "grpc_port", "8001", "grpc啟動通訊埠編號")
	flag.StringVar(&httpPort, "http_port", "9001", "HTTP啟動通訊埠編號")
	flag.Parse()
}

func main() {
	errs := make(chan error) //創建通道以及時接收錯誤
	go func() {
		err := RunHttpServer(httpPort)
		if err != nil {
			errs <- err
		}
	}()

	go func() { //啟動輕量級的執行緒
		err := RunGrpcServer(grpcPort)
		if err != nil {
			errs <- err
		}
	}()

	select { //等待通道來東西
	case err := <-errs:
		log.Fatalf("Run Server err: %v", err)
	}
}

func RunHttpServer(port string) error { //針對http1.0
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/ping",
		func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`pong`))
		},
	)
	return http.ListenAndServe(":"+port, serveMux)
}

func RunGrpcServer(port string) error { //針對grpc
	s := grpc.NewServer()
	pb.RegisterTagServiceServer(s, server.NewTagServer())
	reflection.Register(s)

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}

	return s.Serve(lis)
}
