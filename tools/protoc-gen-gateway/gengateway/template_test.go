package gengateway

import (
	descriptor2 "github.com/hopeio/protobuf/tools/protoc-gen-gateway/descriptor"
	"github.com/hopeio/protobuf/tools/protoc-gen-gateway/httprule"
	"strings"
	"testing"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

func crossLinkFixture(f *descriptor2.File) *descriptor2.File {
	for _, m := range f.Messages {
		m.File = f
	}
	for _, svc := range f.Services {
		svc.File = f
		for _, m := range svc.Methods {
			m.Service = svc
			for _, b := range m.Bindings {
				b.Method = m
				for _, param := range b.PathParams {
					param.Method = m
				}
			}
		}
	}
	return f
}

func TestApplyTemplateHeader(t *testing.T) {
	msgdesc := &descriptorpb.DescriptorProto{
		Name: proto.String("ExampleMessage"),
	}
	meth := &descriptorpb.MethodDescriptorProto{
		Name:       proto.String("Example"),
		InputType:  proto.String("ExampleMessage"),
		OutputType: proto.String("ExampleMessage"),
	}
	svc := &descriptorpb.ServiceDescriptorProto{
		Name:   proto.String("ExampleService"),
		Method: []*descriptorpb.MethodDescriptorProto{meth},
	}
	msg := &descriptor2.Message{
		DescriptorProto: msgdesc,
	}
	file := descriptor2.File{
		FileDescriptorProto: &descriptorpb.FileDescriptorProto{
			Name:        proto.String("example.proto"),
			Package:     proto.String("example"),
			Dependency:  []string{"a.example/b/c.proto", "a.example/d/e.proto"},
			MessageType: []*descriptorpb.DescriptorProto{msgdesc},
			Service:     []*descriptorpb.ServiceDescriptorProto{svc},
		},
		GoPkg: descriptor2.GoPackage{
			Path: "example.com/path/to/example/example.pb",
			Name: "example_pb",
		},
		Messages: []*descriptor2.Message{msg},
		Services: []*descriptor2.Service{
			{
				ServiceDescriptorProto: svc,
				Methods: []*descriptor2.Method{
					{
						MethodDescriptorProto: meth,
						RequestType:           msg,
						ResponseType:          msg,
						Bindings: []*descriptor2.Binding{
							{
								HTTPMethod: "GET",
								Body:       &descriptor2.Body{FieldPath: nil},
							},
						},
					},
				},
			},
		},
	}
	got, err := applyTemplate(param{File: crossLinkFixture(&file), RegisterFuncSuffix: "Handler", AllowPatchFeature: true, Framework: FrameworkGin}, descriptor2.NewRegistry())
	if err != nil {
		t.Errorf("applyTemplate(%#v) failed with %v; want success", file, err)
		return
	}
	if want := "package example_pb\n"; !strings.Contains(got, want) {
		t.Errorf("applyTemplate(%#v) = %s; want to contain %s", file, got, want)
	}
}

func TestApplyTemplateRequestWithoutClientStreaming(t *testing.T) {
	msgdesc := &descriptorpb.DescriptorProto{
		Name: proto.String("ExampleMessage"),
		Field: []*descriptorpb.FieldDescriptorProto{
			{
				Name:     proto.String("nested"),
				Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
				Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
				TypeName: proto.String("NestedMessage"),
				Number:   proto.Int32(1),
			},
		},
	}
	nesteddesc := &descriptorpb.DescriptorProto{
		Name: proto.String("NestedMessage"),
		Field: []*descriptorpb.FieldDescriptorProto{
			{
				Name:   proto.String("int32"),
				Label:  descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
				Type:   descriptorpb.FieldDescriptorProto_TYPE_INT32.Enum(),
				Number: proto.Int32(1),
			},
			{
				Name:   proto.String("bool"),
				Label:  descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
				Type:   descriptorpb.FieldDescriptorProto_TYPE_BOOL.Enum(),
				Number: proto.Int32(2),
			},
		},
	}
	meth := &descriptorpb.MethodDescriptorProto{
		Name:            proto.String("Echo"),
		InputType:       proto.String("ExampleMessage"),
		OutputType:      proto.String("ExampleMessage"),
		ClientStreaming: proto.Bool(false),
	}
	svc := &descriptorpb.ServiceDescriptorProto{
		Name:   proto.String("ExampleService"),
		Method: []*descriptorpb.MethodDescriptorProto{meth},
	}
	for _, spec := range []struct {
		serverStreaming bool
		sigWant         string
	}{
		{
			serverStreaming: false,
			sigWant:         `gateway.UnaryCall(server.Echo)`,
		},
		{
			serverStreaming: true,
			sigWant:         `gateway.ServerSideStreamCall(server.Echo)`,
		},
	} {
		meth.ServerStreaming = proto.Bool(spec.serverStreaming)

		msg := &descriptor2.Message{
			DescriptorProto: msgdesc,
		}
		nested := &descriptor2.Message{
			DescriptorProto: nesteddesc,
		}

		nestedField := &descriptor2.Field{
			Message:              msg,
			FieldDescriptorProto: msg.GetField()[0],
		}
		intField := &descriptor2.Field{
			Message:              nested,
			FieldDescriptorProto: nested.GetField()[0],
		}
		boolField := &descriptor2.Field{
			Message:              nested,
			FieldDescriptorProto: nested.GetField()[1],
		}
		file := descriptor2.File{
			FileDescriptorProto: &descriptorpb.FileDescriptorProto{
				Name:        proto.String("example.proto"),
				Package:     proto.String("example"),
				MessageType: []*descriptorpb.DescriptorProto{msgdesc, nesteddesc},
				Service:     []*descriptorpb.ServiceDescriptorProto{svc},
			},
			GoPkg: descriptor2.GoPackage{
				Path: "example.com/path/to/example/example.pb",
				Name: "example_pb",
			},
			Messages: []*descriptor2.Message{msg, nested},
			Services: []*descriptor2.Service{
				{
					ServiceDescriptorProto: svc,
					Methods: []*descriptor2.Method{
						{
							MethodDescriptorProto: meth,
							RequestType:           msg,
							ResponseType:          msg,
							Bindings: []*descriptor2.Binding{
								{
									HTTPMethod: "POST",
									PathTmpl: httprule.Template{
										Version: 1,
										OpCodes: []int{0, 0},
									},
									PathParams: []descriptor2.Parameter{
										{
											FieldPath: descriptor2.FieldPath([]descriptor2.FieldPathComponent{
												{
													Name:   "nested",
													Target: nestedField,
												},
												{
													Name:   "int32",
													Target: intField,
												},
											}),
											Target: intField,
										},
									},
									Body: &descriptor2.Body{
										FieldPath: descriptor2.FieldPath([]descriptor2.FieldPathComponent{
											{
												Name:   "nested",
												Target: nestedField,
											},
											{
												Name:   "bool",
												Target: boolField,
											},
										}),
									},
								},
							},
						},
					},
				},
			},
		}
		got, err := applyTemplate(param{File: crossLinkFixture(&file), RegisterFuncSuffix: "Handler", AllowPatchFeature: true, Framework: FrameworkGin}, descriptor2.NewRegistry())
		if err != nil {
			t.Errorf("applyTemplate(%#v) failed with %v; want success", file, err)
			return
		}
		if want := spec.sigWant; !strings.Contains(got, want) {
			t.Errorf("applyTemplate(%#v) = %s; want to contain %s", file, got, want)
		}
		if want := `func RegisterExampleServiceHandlerServer(mux *gin.Engine, server ExampleServiceServer) {`; !strings.Contains(got, want) {
			t.Errorf("applyTemplate(%#v) = %s; want to contain %s", file, got, want)
		}
	}
}

func TestApplyTemplateRequestWithClientStreaming(t *testing.T) {
	msgdesc := &descriptorpb.DescriptorProto{
		Name: proto.String("ExampleMessage"),
		Field: []*descriptorpb.FieldDescriptorProto{
			{
				Name:     proto.String("nested"),
				Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
				Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
				TypeName: proto.String("NestedMessage"),
				Number:   proto.Int32(1),
			},
		},
	}
	nesteddesc := &descriptorpb.DescriptorProto{
		Name: proto.String("NestedMessage"),
		Field: []*descriptorpb.FieldDescriptorProto{
			{
				Name:   proto.String("int32"),
				Label:  descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
				Type:   descriptorpb.FieldDescriptorProto_TYPE_INT32.Enum(),
				Number: proto.Int32(1),
			},
			{
				Name:   proto.String("bool"),
				Label:  descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
				Type:   descriptorpb.FieldDescriptorProto_TYPE_BOOL.Enum(),
				Number: proto.Int32(2),
			},
		},
	}
	meth := &descriptorpb.MethodDescriptorProto{
		Name:            proto.String("Echo"),
		InputType:       proto.String("ExampleMessage"),
		OutputType:      proto.String("ExampleMessage"),
		ClientStreaming: proto.Bool(true),
	}
	svc := &descriptorpb.ServiceDescriptorProto{
		Name:   proto.String("ExampleService"),
		Method: []*descriptorpb.MethodDescriptorProto{meth},
	}
	for _, spec := range []struct {
		serverStreaming bool
		sigWant         string
	}{
		{
			serverStreaming: false,
			sigWant:         `gateway.ClientSideStreamCall(server.Echo)`,
		},
		{
			serverStreaming: true,
			sigWant:         `gateway.BidiStreamCall(server.Echo)`,
		},
	} {
		meth.ServerStreaming = proto.Bool(spec.serverStreaming)

		msg := &descriptor2.Message{
			DescriptorProto: msgdesc,
		}
		nested := &descriptor2.Message{
			DescriptorProto: nesteddesc,
		}

		nestedField := &descriptor2.Field{
			Message:              msg,
			FieldDescriptorProto: msg.GetField()[0],
		}
		intField := &descriptor2.Field{
			Message:              nested,
			FieldDescriptorProto: nested.GetField()[0],
		}
		boolField := &descriptor2.Field{
			Message:              nested,
			FieldDescriptorProto: nested.GetField()[1],
		}
		file := descriptor2.File{
			FileDescriptorProto: &descriptorpb.FileDescriptorProto{
				Name:        proto.String("example.proto"),
				Package:     proto.String("example"),
				MessageType: []*descriptorpb.DescriptorProto{msgdesc, nesteddesc},
				Service:     []*descriptorpb.ServiceDescriptorProto{svc},
			},
			GoPkg: descriptor2.GoPackage{
				Path: "example.com/path/to/example/example.pb",
				Name: "example_pb",
			},
			Messages: []*descriptor2.Message{msg, nested},
			Services: []*descriptor2.Service{
				{
					ServiceDescriptorProto: svc,
					Methods: []*descriptor2.Method{
						{
							MethodDescriptorProto: meth,
							RequestType:           msg,
							ResponseType:          msg,
							Bindings: []*descriptor2.Binding{
								{
									HTTPMethod: "POST",
									PathTmpl: httprule.Template{
										Version: 1,
										OpCodes: []int{0, 0},
									},
									PathParams: []descriptor2.Parameter{
										{
											FieldPath: descriptor2.FieldPath([]descriptor2.FieldPathComponent{
												{
													Name:   "nested",
													Target: nestedField,
												},
												{
													Name:   "int32",
													Target: intField,
												},
											}),
											Target: intField,
										},
									},
									Body: &descriptor2.Body{
										FieldPath: descriptor2.FieldPath([]descriptor2.FieldPathComponent{
											{
												Name:   "nested",
												Target: nestedField,
											},
											{
												Name:   "bool",
												Target: boolField,
											},
										}),
									},
								},
							},
						},
					},
				},
			},
		}
		got, err := applyTemplate(param{File: crossLinkFixture(&file), RegisterFuncSuffix: "Handler", AllowPatchFeature: true, Framework: FrameworkGin}, descriptor2.NewRegistry())
		if err != nil {
			t.Errorf("applyTemplate(%#v) failed with %v; want success", file, err)
			return
		}
		if want := spec.sigWant; !strings.Contains(got, want) {
			t.Errorf("applyTemplate(%#v) = %s; want to contain %s", file, got, want)
		}
		if want := `func RegisterExampleServiceHandlerServer(mux *gin.Engine, server ExampleServiceServer) {`; !strings.Contains(got, want) {
			t.Errorf("applyTemplate(%#v) = %s; want to contain %s", file, got, want)
		}
	}
}

func TestApplyTemplateInProcess(t *testing.T) {
	msgdesc := &descriptorpb.DescriptorProto{
		Name: proto.String("ExampleMessage"),
		Field: []*descriptorpb.FieldDescriptorProto{
			{
				Name:     proto.String("nested"),
				Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
				Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
				TypeName: proto.String("NestedMessage"),
				Number:   proto.Int32(1),
			},
		},
	}
	nesteddesc := &descriptorpb.DescriptorProto{
		Name: proto.String("NestedMessage"),
		Field: []*descriptorpb.FieldDescriptorProto{
			{
				Name:   proto.String("int32"),
				Label:  descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
				Type:   descriptorpb.FieldDescriptorProto_TYPE_INT32.Enum(),
				Number: proto.Int32(1),
			},
			{
				Name:   proto.String("bool"),
				Label:  descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
				Type:   descriptorpb.FieldDescriptorProto_TYPE_BOOL.Enum(),
				Number: proto.Int32(2),
			},
		},
	}
	meth := &descriptorpb.MethodDescriptorProto{
		Name:            proto.String("Echo"),
		InputType:       proto.String("ExampleMessage"),
		OutputType:      proto.String("ExampleMessage"),
		ClientStreaming: proto.Bool(true),
	}
	svc := &descriptorpb.ServiceDescriptorProto{
		Name:   proto.String("ExampleService"),
		Method: []*descriptorpb.MethodDescriptorProto{meth},
	}
	for _, spec := range []struct {
		clientStreaming bool
		serverStreaming bool
		sigWant         []string
	}{
		{
			clientStreaming: false,
			serverStreaming: false,
			sigWant: []string{
				`gateway.UnaryCall(server.Echo)`,
			},
		},
		{
			clientStreaming: true,
			serverStreaming: true,
			sigWant: []string{
				`gateway.BidiStreamCall(server.Echo)`,
			},
		},
		{
			clientStreaming: true,
			serverStreaming: false,
			sigWant: []string{
				`gateway.ClientSideStreamCall(server.Echo)`,
			},
		},
		{
			clientStreaming: false,
			serverStreaming: true,
			sigWant: []string{
				`gateway.ServerSideStreamCall(server.Echo)`,
			},
		},
	} {
		meth.ClientStreaming = proto.Bool(spec.clientStreaming)
		meth.ServerStreaming = proto.Bool(spec.serverStreaming)

		msg := &descriptor2.Message{
			DescriptorProto: msgdesc,
		}
		nested := &descriptor2.Message{
			DescriptorProto: nesteddesc,
		}

		nestedField := &descriptor2.Field{
			Message:              msg,
			FieldDescriptorProto: msg.GetField()[0],
		}
		intField := &descriptor2.Field{
			Message:              nested,
			FieldDescriptorProto: nested.GetField()[0],
		}
		boolField := &descriptor2.Field{
			Message:              nested,
			FieldDescriptorProto: nested.GetField()[1],
		}
		file := descriptor2.File{
			FileDescriptorProto: &descriptorpb.FileDescriptorProto{
				Name:        proto.String("example.proto"),
				Package:     proto.String("example"),
				MessageType: []*descriptorpb.DescriptorProto{msgdesc, nesteddesc},
				Service:     []*descriptorpb.ServiceDescriptorProto{svc},
			},
			GoPkg: descriptor2.GoPackage{
				Path: "example.com/path/to/example/example.pb",
				Name: "example_pb",
			},
			Messages: []*descriptor2.Message{msg, nested},
			Services: []*descriptor2.Service{
				{
					ServiceDescriptorProto: svc,
					Methods: []*descriptor2.Method{
						{
							MethodDescriptorProto: meth,
							RequestType:           msg,
							ResponseType:          msg,
							Bindings: []*descriptor2.Binding{
								{
									HTTPMethod: "POST",
									PathTmpl: httprule.Template{
										Version: 1,
										OpCodes: []int{0, 0},
									},
									PathParams: []descriptor2.Parameter{
										{
											FieldPath: descriptor2.FieldPath([]descriptor2.FieldPathComponent{
												{
													Name:   "nested",
													Target: nestedField,
												},
												{
													Name:   "int32",
													Target: intField,
												},
											}),
											Target: intField,
										},
									},
									Body: &descriptor2.Body{
										FieldPath: descriptor2.FieldPath([]descriptor2.FieldPathComponent{
											{
												Name:   "nested",
												Target: nestedField,
											},
											{
												Name:   "bool",
												Target: boolField,
											},
										}),
									},
								},
							},
						},
					},
				},
			},
		}
		got, err := applyTemplate(param{File: crossLinkFixture(&file), RegisterFuncSuffix: "Handler", AllowPatchFeature: true, Framework: FrameworkGin}, descriptor2.NewRegistry())
		if err != nil {
			t.Errorf("applyTemplate(%#v) failed with %v; want success", file, err)
			return
		}

		for _, want := range spec.sigWant {
			if !strings.Contains(got, want) {
				t.Errorf("applyTemplate(%#v) = %s; want to contain %s", file, got, want)
			}
		}

		if want := `func RegisterExampleServiceHandlerServer(mux *gin.Engine, server ExampleServiceServer)`; !strings.Contains(got, want) {
			t.Errorf("applyTemplate(%#v) = %s; want to contain %s", file, got, want)
		}
	}
}

func TestAllowPatchFeature(t *testing.T) {
	updateMaskDesc := &descriptorpb.FieldDescriptorProto{
		Name:     proto.String("UpdateMask"),
		Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
		Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
		TypeName: proto.String(".google.protobuf.FieldMask"),
		Number:   proto.Int32(1),
	}
	msgdesc := &descriptorpb.DescriptorProto{
		Name:  proto.String("ExampleMessage"),
		Field: []*descriptorpb.FieldDescriptorProto{updateMaskDesc},
	}
	meth := &descriptorpb.MethodDescriptorProto{
		Name:       proto.String("Example"),
		InputType:  proto.String("ExampleMessage"),
		OutputType: proto.String("ExampleMessage"),
	}
	svc := &descriptorpb.ServiceDescriptorProto{
		Name:   proto.String("ExampleService"),
		Method: []*descriptorpb.MethodDescriptorProto{meth},
	}
	msg := &descriptor2.Message{
		DescriptorProto: msgdesc,
	}
	updateMaskField := &descriptor2.Field{
		Message:              msg,
		FieldDescriptorProto: updateMaskDesc,
	}
	msg.Fields = append(msg.Fields, updateMaskField)
	file := descriptor2.File{
		FileDescriptorProto: &descriptorpb.FileDescriptorProto{
			Name:        proto.String("example.proto"),
			Package:     proto.String("example"),
			MessageType: []*descriptorpb.DescriptorProto{msgdesc},
			Service:     []*descriptorpb.ServiceDescriptorProto{svc},
		},
		GoPkg: descriptor2.GoPackage{
			Path: "example.com/path/to/example/example.pb",
			Name: "example_pb",
		},
		Messages: []*descriptor2.Message{msg},
		Services: []*descriptor2.Service{
			{
				ServiceDescriptorProto: svc,
				Methods: []*descriptor2.Method{
					{
						MethodDescriptorProto: meth,
						RequestType:           msg,
						ResponseType:          msg,
						Bindings: []*descriptor2.Binding{
							{
								HTTPMethod: "PATCH",
								Body: &descriptor2.Body{FieldPath: descriptor2.FieldPath{descriptor2.FieldPathComponent{
									Name:   "abe",
									Target: msg.Fields[0],
								}}},
							},
						},
					},
				},
			},
		},
	}
	want := `gateway.UnaryCall(server.Example)`
	for _, allowPatchFeature := range []bool{true, false} {
		got, err := applyTemplate(param{File: crossLinkFixture(&file), RegisterFuncSuffix: "Handler", AllowPatchFeature: allowPatchFeature, Framework: FrameworkGin}, descriptor2.NewRegistry())
		if err != nil {
			t.Errorf("applyTemplate(%#v) failed with %v; want success", file, err)
			return
		}
		if !strings.Contains(got, want) {
			t.Errorf("applyTemplate(%#v) = %s; want to contain %s", file, got, want)
		}
	}
}

func TestIdentifierCapitalization(t *testing.T) {
	msgdesc1 := &descriptorpb.DescriptorProto{
		Name: proto.String("Exam_pleRequest"),
	}
	msgdesc2 := &descriptorpb.DescriptorProto{
		Name: proto.String("example_response"),
	}
	meth1 := &descriptorpb.MethodDescriptorProto{
		Name:       proto.String("ExampleGe2t"),
		InputType:  proto.String("Exam_pleRequest"),
		OutputType: proto.String("example_response"),
	}
	meth2 := &descriptorpb.MethodDescriptorProto{
		Name:       proto.String("Exampl_eGet"),
		InputType:  proto.String("Exam_pleRequest"),
		OutputType: proto.String("example_response"),
	}
	svc := &descriptorpb.ServiceDescriptorProto{
		Name:   proto.String("Example"),
		Method: []*descriptorpb.MethodDescriptorProto{meth1, meth2},
	}
	msg1 := &descriptor2.Message{
		DescriptorProto: msgdesc1,
	}
	msg2 := &descriptor2.Message{
		DescriptorProto: msgdesc2,
	}
	file := descriptor2.File{
		FileDescriptorProto: &descriptorpb.FileDescriptorProto{
			Name:        proto.String("example.proto"),
			Package:     proto.String("example"),
			Dependency:  []string{"a.example/b/c.proto", "a.example/d/e.proto"},
			MessageType: []*descriptorpb.DescriptorProto{msgdesc1, msgdesc2},
			Service:     []*descriptorpb.ServiceDescriptorProto{svc},
		},
		GoPkg: descriptor2.GoPackage{
			Path: "example.com/path/to/example/example.pb",
			Name: "example_pb",
		},
		Messages: []*descriptor2.Message{msg1, msg2},
		Services: []*descriptor2.Service{
			{
				ServiceDescriptorProto: svc,
				Methods: []*descriptor2.Method{
					{
						MethodDescriptorProto: meth1,
						RequestType:           msg1,
						ResponseType:          msg1,
						Bindings: []*descriptor2.Binding{
							{
								HTTPMethod: "GET",
								Body:       &descriptor2.Body{FieldPath: nil},
							},
						},
					},
				},
			},
			{
				ServiceDescriptorProto: svc,
				Methods: []*descriptor2.Method{
					{
						MethodDescriptorProto: meth2,
						RequestType:           msg2,
						ResponseType:          msg2,
						Bindings: []*descriptor2.Binding{
							{
								HTTPMethod: "GET",
								Body:       &descriptor2.Body{FieldPath: nil},
							},
						},
					},
				},
			},
		},
	}

	got, err := applyTemplate(param{File: crossLinkFixture(&file), RegisterFuncSuffix: "Handler", AllowPatchFeature: true, Framework: FrameworkGin}, descriptor2.NewRegistry())
	if err != nil {
		t.Errorf("applyTemplate(%#v) failed with %v; want success", file, err)
		return
	}
	if want := `gateway.UnaryCall(server.ExampleGe2T)`; !strings.Contains(got, want) {
		t.Errorf("applyTemplate(%#v) = %s; want to contain %s", file, got, want)
	}
	if want := `gateway.UnaryCall(server.ExamplEGet)`; !strings.Contains(got, want) {
		t.Errorf("applyTemplate(%#v) = %s; want to contain %s", file, got, want)
	}
}

func TestApplyTemplateFrameworks(t *testing.T) {
	file := minimalServiceFile(t)
	for _, spec := range []struct {
		fw   Framework
		want []string
	}{
		{FrameworkGin, []string{`*gin.Engine`, `mux.Handle("GET"`}},
		{FrameworkFiber, []string{`*fiber.App`, `app.Add("GET"`}},
		{FrameworkNetHTTP, []string{`*http.ServeMux`, `mux.Handle("GET /v1/ping"`}},
	} {
		got, err := applyTemplate(param{File: crossLinkFixture(file), RegisterFuncSuffix: "Handler", AllowPatchFeature: true, Framework: spec.fw}, descriptor2.NewRegistry())
		if err != nil {
			t.Fatalf("framework %s: %v", spec.fw, err)
		}
		for _, w := range spec.want {
			if !strings.Contains(got, w) {
				t.Fatalf("framework %s: want %q in:\n%s", spec.fw, w, got)
			}
		}
	}
}

func minimalServiceFile(t *testing.T) *descriptor2.File {
	t.Helper()
	msgdesc := &descriptorpb.DescriptorProto{Name: proto.String("PingRequest")}
	meth := &descriptorpb.MethodDescriptorProto{
		Name: proto.String("Ping"), InputType: proto.String("PingRequest"), OutputType: proto.String("PingRequest"),
	}
	svc := &descriptorpb.ServiceDescriptorProto{Name: proto.String("DemoService"), Method: []*descriptorpb.MethodDescriptorProto{meth}}
	msg := &descriptor2.Message{DescriptorProto: msgdesc}
	return &descriptor2.File{
		FileDescriptorProto: &descriptorpb.FileDescriptorProto{
			Name: proto.String("demo.proto"), Package: proto.String("demo"),
			MessageType: []*descriptorpb.DescriptorProto{msgdesc},
			Service:     []*descriptorpb.ServiceDescriptorProto{svc},
		},
		GoPkg: descriptor2.GoPackage{Path: "example.com/demo", Name: "demo"},
		Messages: []*descriptor2.Message{msg},
		Services: []*descriptor2.Service{{
			ServiceDescriptorProto: svc,
			Methods: []*descriptor2.Method{{
				MethodDescriptorProto: meth,
				RequestType: msg, ResponseType: msg,
				Bindings: []*descriptor2.Binding{{HTTPMethod: "GET", PathTmpl: httprule.Template{Template: "/v1/ping"}}},
			}},
		}},
	}
}
