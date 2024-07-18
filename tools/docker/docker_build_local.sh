gopath=/mnt/d/SDK/gopath
if [ -z "$1" ]; then
  echo "gopath参数为空"
  exit 1
else
  gopath=$1
  echo "gopath: $gopath"
fi
protoc=/mnt/d/tools/protoc-22.3-linux-x86_64
if [ -z "$2" ]; then
  echo "protoc参数为空"
  exit 1
else
  protoc=$2
  echo "protoc: $protoc"
fi
prorobuf_dir=$(cd $(dirname $0)/../..;pwd)
echo "prorobuf_dir: $prorobuf_dir"

goproxy=https://goproxy.io,https://goproxy.cn,direct
goimage=golang:latest

if [ -n "$3" ]; then
  goimage= $1
fi

# install tools
docker run --rm -e GOPROXY=$goproxy -e GOFLAGS=-buildvcs=false -v $gopath:/go -v $protoc:/protoc -v $prorobuf_dir:/work -w /work/tools --name install $goimage bash ./install_tools.sh /protoc
# docker rm -f install
echo "docker build"
dockerTmpDir=$prorobuf_dir/tools/docker_tmp
mkdir $dockerTmpDir
cp $gopath/bin/protoc-gen-enum $dockerTmpDir/
cp $gopath/bin/protoc-gen-go $dockerTmpDir/
cp $gopath/bin/protoc-gen-go-grpc $dockerTmpDir/
cp $gopath/bin/protoc-gen-go-patch $dockerTmpDir/
cp $gopath/bin/protoc-gen-validator $dockerTmpDir/
cp $gopath/bin/protoc-gen-grpc-gateway $dockerTmpDir/
cp $gopath/bin/protoc-gen-grpc-gin $dockerTmpDir/
cp $gopath/bin/protoc-gen-openapiv2 $dockerTmpDir/
cp $gopath/bin/protoc-gen-gql $dockerTmpDir/
cp $gopath/bin/protoc-gen-gogql $dockerTmpDir/
cp $gopath/bin/gqlgen $dockerTmpDir/
cp $gopath/bin/protogen $dockerTmpDir/
cp -r $prorobuf_dir/_proto $dockerTmpDir/_proto
cp -r $protoc $dockerTmpDir/protoc


docker build -t jybl/protogen --build-arg IMAGE=$goimage -f $prorobuf_dir/tools/docker/Dockerfile_local $dockerTmpDir
rm -rf $dockerTmpDir
docker push jybl/protogen