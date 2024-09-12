# RPC_application

### 安裝 protoc
參考安裝[網址](https://github.com/protocolbuffers/protobuf/releases/tag/v3.11.2)
[參考文章](https://ithelp.ithome.com.tw/articles/10250131)
或是直接下安裝指令( wins 系統需要先安裝 choco 套件，並且用系統管理員執行下載)
```
Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://chocolatey.org/install.ps1'))
choco install protoc
protoc --version
```

### 安裝 protoc 外掛程式
```
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

### 編譯 proto 檔案後產生相應的 pb.go 檔案
```
protoc --go_out=. --go-grpc_out=. ./proto/*.proto
```

### 安裝 grpc 套件
```
go get -u google.golang.org/grpc
```

### Rpc 種類
Unary RPC: 一元RPC。會限制資料都接收成功且正確後才會進行下一步處理。
Server-side  streaming RPC: 服務端流式RPC
Client-side streaming RPC: 用戶端流式RPC
Bidirectional streaming RPC: 雙向流式 RPC

### 用於驗證的 grpc 套件
安裝
```
go get -u github.com/fullstorydev/grpcurl
go install github.com/fullstorydev/grpcurl/cmd/grpcurl
```

使用(需重新啟動服務)
```
grpcurl -plaintext -d '{"name":"Go"}' localhost:8001 proto.TagService.GetTagList
```

### 多協定服務(gRPC、SSH、HTTPS等)套件 & RESTful 轉 gRPC 請求
用同一個通訊埠去支援多協定
```
go get -u github.com/soheilhy/cmux@v0.1.4
```
下載套件並移動
```
go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@v1.14.5
go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@v1.14.5
mv $GOPATH/bin/protoc-gen-grpc-gateway /usr/local/go/bin
```
import編譯
```
protoc -I$GOPATH -I. -I$GOPATH/pkg/mod -I$GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.14.5/third_party/googleapis --grpc-gateway_out=logtostderr=true:. ./proto/*.proto
```