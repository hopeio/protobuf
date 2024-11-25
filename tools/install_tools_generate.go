//go:build tools

/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package main

import (
	_ "github.com/99designs/gqlgen"
	_ "github.com/alta/protopatch/cmd/protoc-gen-go-patch"
	//_ "github.com/bufbuild/protovalidate-go"
	_ "github.com/danielvladco/go-proto-gql/pkg/graphqlpb"
	//_ "github.com/envoyproxy/protoc-gen-validate"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)

//go:generate go run ./install_tools.go
