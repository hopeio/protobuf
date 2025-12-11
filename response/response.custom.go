/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package response

import (
	"context"
	"net/http"

	httpx "github.com/hopeio/gox/net/http"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var RespOk = &ErrResp{}

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

func (x *HttpResponse) Respond(ctx context.Context, w http.ResponseWriter) (int, error) {
	return x.CommonRespond(ctx, httpx.ResponseWriterWrapper{ResponseWriter: w})
}

func (x *HttpResponse) CommonRespond(ctx context.Context, w httpx.CommonResponseWriter) (int, error) {
	header := w.Header()
	for k, v := range x.Headers {
		header.Set(k, v)
	}
	return w.Write(x.Body)
}

func (x *HttpResponse) ServerHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	return x.CommonRespond(r.Context(), httpx.ResponseWriterWrapper{ResponseWriter: w})
}
