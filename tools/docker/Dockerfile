ARG IMAGE=golang:latest

FROM ${IMAGE}

ENV PATH /protoc/bin:$PATH
ENV GOFLAGS -buildvcs=false


RUN cd / && git clone https://github.com/hopeio/protobuf
RUN cd /protobuf && go mod download && go mod tidy && copy _proto/ /proto/
RUN cd tools  && ./install_protoc.sh && mv protoc-26.1-linux-x86_64 /protoc && rm protoc-26.1-linux-x86_64.zip
RUN go run /protobuf/tools/install_tools.go
RUN rm /protobuf

