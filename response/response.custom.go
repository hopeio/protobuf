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

func (x *HttpResponse) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	x.Respond(r.Context(), w)
}

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
	return w.Write(x.Body)
}

func (x *ErrResp) ErrResp() *mix.ErrResp {
	return &mix.ErrResp{
		Code: mix.ErrCode(x.Code),
		Msg:  x.Msg,
	}
}

func (x *ErrResp) Error() string {
	return x.ErrResp().Error()
}
