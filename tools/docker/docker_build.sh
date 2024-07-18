prorobuf_dir=$(cd $(dirname $0)/../..;pwd)

goimage=golang:latest
if [ -n "$1" ]; then
  goimage= $1
fi

dockerTmpDir=$prorobuf_dir/tools/docker_tmp
mkdir $dockerTmpDir && cd $dockerTmpDir
docker build -t jybl/protogen --build-arg IMAGE=$goimage -f $prorobuf_dir/tools/docker/Dockerfile $dockerTmpDir
rm -rf $dockerTmpDir
docker push jybl/protogen