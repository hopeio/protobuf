/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package response

import (
	"errors"
	httpi "github.com/hopeio/utils/net/http"
	"github.com/hopeio/utils/net/http/consts"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"io"
	"net/http"

	"google.golang.org/protobuf/proto"
)

type Reply struct {
	Code uint32
	Msg  string
	Data proto.Message
}

func (x *HttpResponse) GetContentType() string {
	return x.Headers[consts.HeaderContentType]
}

func (x HttpResponse) MarshalProto(w io.Writer) {
	w.Write(x.Body)
}

func (x *HttpResponse) MarshalGQL(w io.Writer) {
	w.Write(x.Body)
}

func (x *HttpResponse) UnmarshalGQL(v interface{}) error {
	if i, ok := v.([]byte); ok {
		x.Body = i
		return nil
	}
	return errors.New("error type")
}

var ResponseOk = &TinyRep{}

type StringValue = wrapperspb.StringValue

// graphql compatible
type StringValueInput = wrapperspb.StringValue

/*// graphql compatible
type HeaderEntry struct {
	Key   string
	Value string
}

type HttpResponseResolver struct {
}

// graphql compatible
func (receiver *HttpResponseResolver) Header(ctx context.Context, obj *HttpResponse) ([]*HeaderEntry, error) {
	var header []*HeaderEntry
	for k, v := range obj.Header {
		header = append(header, &HeaderEntry{Key: k, Value: v})
	}
	return header, nil
}
*/

func (res *HttpResponse) Response(w http.ResponseWriter) (int, error) {
	return res.CommonResponse(httpi.CommonResponseWriter{w})
}

func (x *HttpResponse) CommonResponse(w httpi.ICommonResponseWriter) (int, error) {
	header := w.Header()
	for k, v := range x.Headers {
		header.Set(k, v)
	}
	return w.Write(x.Body)
}
