cd $(dirname $0) && pwd
prorobuf_dir=$(cd ../../..;pwd)

goimage=golang:latest

if [ -n "$1" ]; then
  goimage= $1
fi

docker build -t jybl/protogen --build-arg IMAGE=$goimage -f $prorobuf_dir/tools/Dockerfile_upgrade .