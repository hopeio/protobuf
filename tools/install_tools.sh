cd $(dirname $0) && pwd

if ! command -v go &> /dev/null
then
    echo "go command not found, please install go"
    exit 1
fi

protobuf=$(go list -m -f {{.Dir}}  github.com/hopeio/protobuf)
protobuf=${protobuf//\\/\/}
protoDir=$protobuf/_proto

if [ -n "$1" ]; then
  export PATH=$1/bin:$PATH
  echo $PATH
fi

if ! command -v protoc &> /dev/null
then
    echo "protoc command not found, will download"
    # 在这里执行其他操作
    source ./install_protoc.sh
fi


# 安装
cd $protobuf/tools
echo "Start Installation"
go install google.golang.org/protobuf/cmd/protoc-gen-go
protoc -I$protoDir --go_out=paths=source_relative:$protobuf $protoDir/hopeio/utils/**/*.proto
cp -r $protobuf/hopeio/* $protobuf/ # 与install_tools.go不一致，mv的问题，无法覆盖非空目录
rm -rf $protobuf/hopeio
go install ./protoc-gen-enum
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
#go install github.com/alta/protopatch/cmd/protoc-gen-go-patch
go install ./protoc-gen-grpc-gin
#go install github.com/envoyproxy/protoc-gen-validate
go install ./protoc-gen-validator
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
go install ./protoc-gen-go-patch
# go install github.com/danielvladco/go-proto-gql/protoc-gen-gql
# go install github.com/danielvladco/go-proto-gql/protoc-gen-gogql
go install ./protoc-gen-gql
go install ./protoc-gen-gogql
go install github.com/99designs/gqlgen
go install ./protogen
echo "Installation Finished"
