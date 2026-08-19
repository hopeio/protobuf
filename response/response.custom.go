/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package response

import (
	"context"
	"net/http"

	"github.com/hopeio/mix"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var RespOk = &ErrResp{}

type StringValue = wrapperspb.StringValue

// ServeHTTP implements http.Handler by delegating to Respond.
func (x *HttpResponse) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	x.Respond(r.Context(), w)
}

// Respond writes the response headers, status code and body to the provided ResponseWriter.
func (x *HttpResponse) Respond(ctx context.Context, w http.ResponseWriter) (int, error) {
	if wx, ok := w.(mix.ResponseWriter); ok {
		header := wx.HeaderX()
		for k, v := range x.Headers {
			header.Add(k, v)
		}
	} else {
		header := w.Header()
		for k, v := range x.Headers {
			header.Add(k, v)
		}
	}
	// Status 字段此前被忽略，响应恒为隐式 200
	if x.Status != 0 {
		w.WriteHeader(int(x.Status))
	}
	return w.Write(x.Body)
}

// ErrResp converts the protobuf ErrResp to the mix.ErrResp type.
func (x *ErrResp) ErrResp() *mix.ErrResp {
	return &mix.ErrResp{
		Code: mix.ErrCode(x.Code),
		Msg:  x.Msg,
		Data: x.Data,
	}
}

// Error implements the error interface.
func (x *ErrResp) Error() string {
	return x.ErrResp().Error()
}
