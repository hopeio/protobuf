prorobuf_dir=$(cd $(dirname $0)/../..;pwd)

goimage=golang:latest

if [ -n "$1" ]; then
  goimage= $1
fi

docker build -t jybl/protogen --build-arg IMAGE=$goimage -f $prorobuf_dir/tools/docker/Dockerfile_upgrade .