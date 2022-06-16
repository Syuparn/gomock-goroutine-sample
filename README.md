# gomock-goroutine-sample
sample test for gomock with concurrent dependencies

# prepare

## protoc

```bash
# install tools
$ apt install -y protobuf-compiler
$ go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
$ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
# generate go soruce code from protobuf
$ protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/service.proto
```

## gomock

```bash
# install
$ go install github.com/golang/mock/mockgen@latest
# generate mock
$ go generate ./...
```

