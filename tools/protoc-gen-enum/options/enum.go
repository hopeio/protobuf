/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package options

import (
	protogenx "github.com/hopeio/gox/encoding/protobuf/protogen"
	"github.com/hopeio/protobuf/utils/enum"
	gopb "github.com/hopeio/protobuf/utils/patch"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

func FileOptions(o *protogen.File) *gopb.FileOptions {
	return proto.GetExtension(o.Desc.Options(), gopb.E_File).(*gopb.FileOptions)
}

func ValueOptions(v *protogen.Enum) *gopb.Options {
	return proto.GetExtension(v.Desc.Options(), gopb.E_Enum).(*gopb.Options)
}

func EnumValueOptions(v *protogen.EnumValue) *gopb.Options {
	return proto.GetExtension(v.Desc.Options(), gopb.E_Value).(*gopb.Options)
}

func NoExtGenAll(f *protogen.File) bool {
	return protogenx.GetOptionWithDefault[bool](f.Desc, enum.E_NoExtgenAll, false)
}

func NoExtGen(e *protogen.Enum) bool {
	return protogenx.GetOptionWithDefault[bool](e.Desc, enum.E_NoExtgen, false)
}

func GetEnumComment(ev *protogen.EnumValue) string {
	return protogenx.GetOptionWithDefault[string](ev.Desc, enum.E_Comment, "")
}

func GetEnumType(e *protogen.Enum) string {
	return protogenx.GetOptionWithDefault[string](e.Desc, enum.E_Customtype, "int32")
}

func NoEnumValueMap(e *protogen.Enum) bool {
	return protogenx.GetOptionWithDefault[bool](e.Desc, enum.E_NoGenvaluemap, false)
}

func EnabledEnumNumOrder(e *protogen.Enum) bool {
	return protogenx.GetOptionWithDefault[bool](e.Desc, enum.E_Numorder, false)
}

func EnabledEnumJsonMarshal(f *protogen.File, e *protogen.Enum) bool {
	if v, ok := protogenx.GetOption[bool](e.Desc, enum.E_Jsonmarshal); ok {
		return v
	}
	return protogenx.GetOptionWithDefault[bool](f.Desc, enum.E_JsonmarshalAll, false)
}

func EnabledEnumErrCode(e *protogen.Enum) bool {
	return protogenx.GetOptionWithDefault[bool](e.Desc, enum.E_Errcode, false)
}

func EnabledEnumGqlGen(f *protogen.File, e *protogen.Enum) bool {
	if v, ok := protogenx.GetOption[bool](e.Desc, enum.E_Gqlgen); ok {
		return v
	}

	return protogenx.GetOptionWithDefault[bool](f.Desc, enum.E_GqlgenAll, true)
}

func NoEnumPrefix(f *protogen.File, e *protogen.Enum) bool {
	if v, ok := protogenx.GetOption[bool](e.Desc, enum.E_NoPrefix); ok {
		return v
	}

	return protogenx.GetOptionWithDefault[bool](f.Desc, enum.E_NoPrefixAll, false)
}
