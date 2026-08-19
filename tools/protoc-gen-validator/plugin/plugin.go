package plugin

import (
	"fmt"

	"github.com/hopeio/protobuf/utils/validator"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"

	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"google.golang.org/protobuf/proto"
	descriptor "google.golang.org/protobuf/types/descriptorpb"
)

const uuidPattern = "^([a-fA-F0-9]{8}-" +
	"[a-fA-F0-9]{4}-" +
	"[%s][a-fA-F0-9]{3}-" +
	"[8|9|aA|bB][a-fA-F0-9]{3}-" +
	"[a-fA-F0-9]{12})?$"

type plugin struct {
	*protogen.Plugin
	regexPkg     protogen.GoImportPath
	validatorPkg protogen.GoImportPath
	mixPkg       protogen.GoImportPath
}

// New creates a new validator plugin wrapping the given protogen.Plugin.
func New(p *protogen.Plugin) *plugin {
	return &plugin{
		Plugin:       p,
		regexPkg:     "regexp",
		validatorPkg: "github.com/hopeio/gox/validator",
		mixPkg:       "github.com/hopeio/mix",
	}
}

// Name returns the plugin identifier name.
func (p *plugin) Name() string {
	return "validator"
}

// Generate produces a ".validator.pb.go" file for each proto file that contains messages.
func (p *plugin) Generate() error {
	genFileMap := make(map[string]*protogen.GeneratedFile)

	for _, protoFile := range p.Files {
		if !protoFile.Generate {
			continue
		}

		if len(protoFile.Messages) == 0 {
			continue
		}

		fileName := protoFile.GeneratedFilenamePrefix
		g := p.NewGeneratedFile(fileName+".validator.pb.go", protoFile.GoImportPath)

		genFileMap[fileName] = g
		g.P("package ", protoFile.GoPackageName)

		for _, msg := range protoFile.Messages {
			for _, nmsg := range msg.Messages {
				p.gen(protoFile, nmsg, g)
			}
			p.gen(protoFile, msg, g)
		}
	}
	return nil
}

// gen generates validation code for a single message, dispatching to the appropriate proto syntax handler.
func (p *plugin) gen(file *protogen.File, message *protogen.Message, g *protogen.GeneratedFile) {
	if message.Desc.Options().(*descriptor.MessageOptions).GetMapEntry() {
		return
	}
	p.generateRegexVars(message, g)
	if file.Proto.GetSyntax() == "proto3" {
		p.generateProto3Message(message, g)
	} else {
		p.generateProto2Message(message, g)
	}
}

// getFieldValidatorIfAny returns the FieldValidator option for the field, or nil if absent.
func getFieldValidatorIfAny(field *protogen.Field) *validator.FieldValidator {
	if field.Desc.Options() != nil {
		v := proto.GetExtension(field.Desc.Options(), validator.E_Field)
		if v.(*validator.FieldValidator) != nil {
			return v.(*validator.FieldValidator)
		}
	}
	return nil
}

// getOneofValidatorIfAny returns the OneofValidator option for the oneof field, or nil if absent.
func getOneofValidatorIfAny(oneof *protogen.Oneof) *validator.OneofValidator {
	if oneof.Desc.Options() != nil {
		v := proto.GetExtension(oneof.Desc.Options(), validator.E_Oneof)
		if v.(*validator.OneofValidator) != nil {
			return (v.(*validator.OneofValidator))
		}
	}
	return nil
}

// isSupportedInt reports whether the field kind is one of the integer types with range validation support.
func (p *plugin) isSupportedInt(field *protogen.Field) bool {
	switch field.Desc.Kind() {
	case protoreflect.Int32Kind, protoreflect.Int64Kind:
		return true
	case protoreflect.Uint32Kind, protoreflect.Uint64Kind:
		return true
	case protoreflect.Sint32Kind, protoreflect.Sint64Kind:
		return true
	}
	return false
}

// isSupportedFloat reports whether the field kind is one of the floating-point types with range validation support.
func (p *plugin) isSupportedFloat(field *protogen.Field) bool {
	switch field.Desc.Kind() {
	case protoreflect.FloatKind, protoreflect.DoubleKind:
		return true
	case protoreflect.Fixed32Kind, protoreflect.Fixed64Kind:
		return true
	case protoreflect.Sfixed32Kind, protoreflect.Sfixed64Kind:
		return true
	}
	return false
}

// generateRegexVars emits package-level compiled regex variables for fields with regex or UUID validators.
func (p *plugin) generateRegexVars(message *protogen.Message, g *protogen.GeneratedFile) {
	ccTypeName := message.GoIdent
	for _, field := range message.Fields {
		validator := getFieldValidatorIfAny(field)
		if validator != nil {
			fieldName := field.GoName
			if validator.Regex != nil && validator.UuidVer != nil {
				fmt.Fprintf(os.Stderr, "WARNING: regex and uuid validator is set for field %v.%v, only one of them can be set. Regex and UUID validator is ignored for this field.", ccTypeName, fieldName)
			} else if validator.UuidVer != nil {
				uuid, err := getUUIDRegex(validator.UuidVer)
				if err != nil {
					fmt.Fprintf(os.Stderr, "WARNING: field %v.%v error %s.\n", ccTypeName, fieldName, err)
				} else {
					validator.Regex = &uuid
					g.P(`var `, p.regexName(ccTypeName.GoName, fieldName), ` = `, p.regexPkg.Ident("MustCompile"), `(`, "`", *validator.Regex, "`", `)`)
				}
			} else if validator.Regex != nil {
				g.P(`var `, p.regexName(ccTypeName.GoName, fieldName), ` = `, p.regexPkg.Ident("MustCompile"), `(`, "`", *validator.Regex, "`", `)`)
			}
		}
	}
}

// generateProto2Message emits a Validate() method for a proto2 message, handling nullable pointer semantics.
func (p *plugin) generateProto2Message(message *protogen.Message, g *protogen.GeneratedFile) {
	ccTypeName := message.GoIdent

	g.P(`func (x *`, ccTypeName, `) Validate() error {`)
	for _, field := range message.Fields {
		fieldName := field.GoName
		fieldValidator := getFieldValidatorIfAny(field)
		if fieldValidator == nil && field.Desc.Kind() != protoreflect.MessageKind {
			continue
		}
		if p.validatorWithMessageExists(fieldValidator) {
			fmt.Fprintf(os.Stderr, "WARNING: field %v.%v is a proto2 message, validator.msg_exists has no effect\n", ccTypeName, fieldName)
		}
		variableName := "x." + fieldName
		repeated := field.Desc.IsList()

		nullable := field.Parent != nil
		// For proto2 syntax, only Gogo generates non-pointer fields

		if repeated {
			p.generateRepeatedCountValidator(variableName, fieldName, fieldValidator, g)
			if field.Desc.Kind() == protoreflect.MessageKind || p.validatorWithNonRepeatedConstraint(fieldValidator) {
				g.P(`for _, item := range `, variableName, `{`)
				variableName = "item"
			}
		} else if nullable {
			g.P(`if `, variableName, ` != nil {`)
			if field.Desc.Kind() != protoreflect.BytesKind {
				variableName = "*(" + variableName + ")"
			}
		} else if field.Desc.Kind() != protoreflect.MessageKind {
			variableName = `x.Get` + fieldName + `()`
		}

		if !repeated && fieldValidator != nil {
			if fieldValidator.RepeatedCountMin != nil {
				fmt.Fprintf(os.Stderr, "WARNING: field %v.%v is not repeated, validator.min_elts has no effects\n", ccTypeName, fieldName)
			}
			if fieldValidator.RepeatedCountMax != nil {
				fmt.Fprintf(os.Stderr, "WARNING: field %v.%v is not repeated, validator.max_elts has no effects\n", ccTypeName, fieldName)
			}
		}
		if field.Desc.Kind() == protoreflect.StringKind {
			p.generateStringValidator(variableName, ccTypeName.GoName, fieldName, fieldValidator, g)
		} else if p.isSupportedInt(field) {
			p.generateIntValidator(variableName, ccTypeName.GoName, fieldName, fieldValidator, g)
		} else if field.Desc.Kind() == protoreflect.EnumKind {
			p.generateEnumValidator(field, variableName, ccTypeName.GoName, fieldName, fieldValidator, g)
		} else if p.isSupportedFloat(field) {
			p.generateFloatValidator(variableName, ccTypeName.GoName, fieldName, fieldValidator, g)
		} else if field.Desc.Kind() == protoreflect.BytesKind {
			p.generateLengthValidator(variableName, fieldName, fieldValidator, g)
		} else if field.Desc.Kind() == protoreflect.MessageKind {
			if repeated && nullable {
				variableName = "*(item)"
			}
			g.P(`if err := `, p.validatorPkg.Ident("ValidateStruct"), `(&(`, variableName, `)); err != nil {`)
			g.P(`return err`)
			g.P(`}`)
		}
		if repeated {
			// end the repeated loop
			if field.Desc.Kind() == protoreflect.MessageKind || p.validatorWithNonRepeatedConstraint(fieldValidator) {
				// This internal 'if' cannot be refactored as it would change semantics with respect to the corresponding prelude 'if's
				g.P(`}`)
			}
		} else if nullable {
			g.P(`}`)
		}
	}
	g.P(`return nil`)
	g.P(`}`)
}

// generateProto3Message emits a Validate() method for a proto3 message, handling oneof and zero-value semantics.
func (p *plugin) generateProto3Message(message *protogen.Message, g *protogen.GeneratedFile) {
	ccTypeName := message.GoIdent
	g.P(`func (x *`, ccTypeName, `) Validate() error {`)

	for _, oneof := range message.Oneofs {
		oneofValidator := getOneofValidatorIfAny(oneof)
		if oneofValidator == nil {
			continue
		}
		if oneofValidator.GetRequired() {
			oneOfName := oneof.GoName
			g.P(`if x.Get` + oneOfName + `() == nil {`)
			p.generateErrorString("validator.oneof_required", nil, map[string]string{"field": oneOfName}, g)
			g.P(`}`)
		}
	}
	for _, field := range message.Fields {
		fieldValidator := getFieldValidatorIfAny(field)
		if fieldValidator == nil && field.Desc.Kind() != protoreflect.MessageKind {
			continue
		}
		isOneOf := field.Oneof != nil
		fieldName := field.GoName
		variableName := "x." + fieldName
		repeated := field.Desc.IsList()
		// Golang's proto3 has no concept of unset primitive fields
		nullable := field.Desc.Kind() == protoreflect.MessageKind && field.Parent != nil
		if field.Desc.IsMap() {
			g.P(`// Validation of proto3 map<> fields is unsupported.`)
			continue
		}
		if isOneOf {
			// if x, ok := m.GetType().(*OneOfMessage3_OneInt); ok {
			g.P("if x, ok := x.Get", field.Oneof.GoName, "().(*", field.GoIdent, "); ok {")
			variableName = "x." + field.GoName
		}
		if repeated {
			p.generateRepeatedCountValidator(variableName, fieldName, fieldValidator, g)
			if field.Desc.Kind() == protoreflect.MessageKind || p.validatorWithNonRepeatedConstraint(fieldValidator) {
				g.P(`for _, item := range `, variableName, `{`)
				variableName = "item"
			}
		} else if fieldValidator != nil {
			if fieldValidator.RepeatedCountMin != nil {
				fmt.Fprintf(os.Stderr, "WARNING: field %v.%v is not repeated, validator.min_elts has no effects\n", ccTypeName, fieldName)
			}
			if fieldValidator.RepeatedCountMax != nil {
				fmt.Fprintf(os.Stderr, "WARNING: field %v.%v is not repeated, validator.max_elts has no effects\n", ccTypeName, fieldName)
			}
		}
		if field.Desc.Kind() == protoreflect.StringKind {
			p.generateStringValidator(variableName, ccTypeName.GoName, fieldName, fieldValidator, g)
		} else if p.isSupportedInt(field) {
			p.generateIntValidator(variableName, ccTypeName.GoName, fieldName, fieldValidator, g)
		} else if field.Desc.Kind() == protoreflect.EnumKind {
			p.generateEnumValidator(field, variableName, ccTypeName.GoName, fieldName, fieldValidator, g)
		} else if p.isSupportedFloat(field) {
			p.generateFloatValidator(variableName, ccTypeName.GoName, fieldName, fieldValidator, g)
		} else if field.Desc.Kind() == protoreflect.BytesKind {
			p.generateLengthValidator(variableName, fieldName, fieldValidator, g)
		} else if field.Desc.Kind() == protoreflect.MessageKind {
			if p.validatorWithMessageExists(fieldValidator) {
				if nullable && !repeated {
					g.P(`if nil == `, variableName, `{`)
					p.generateErrorString("validator.msg_exists", fieldValidator, map[string]string{"field": fieldName}, g)
					g.P(`}`)
				} else if repeated {
					fmt.Fprintf(os.Stderr, "WARNING: field %v.%v is repeated, validator.msg_exists has no effect\n", ccTypeName, fieldName)
				} else if !nullable {
					fmt.Fprintf(os.Stderr, "WARNING: field %v.%v is a nullable=false, validator.msg_exists has no effect\n", ccTypeName, fieldName)
				}
			}

			if nullable {
				g.P(`if `, variableName, ` != nil {`)
			} else {
				variableName = "&(" + variableName + ")"
			}
			g.P(`if err := `, p.validatorPkg.Ident("ValidateStruct"), `(`, variableName, `); err != nil {`)
			g.P(`return err`)
			g.P(`}`)
			if nullable {
				g.P(`}`)
			}
		}
		if repeated && (field.Desc.Kind() == protoreflect.MessageKind || p.validatorWithNonRepeatedConstraint(fieldValidator)) {
			// end the repeated loop
			g.P(`}`)
		}
		if isOneOf {
			// end the oneof if statement
			g.P(`}`)
		}
	}
	g.P(`return nil`)
	g.P(`}`)
}

// generateIntValidator emits range check code for integer fields based on IntGt/IntLt constraints.
func (p *plugin) generateIntValidator(variableName string, ccTypeName string, fieldName string, fv *validator.FieldValidator, g *protogen.GeneratedFile) {
	if fv.IntGt != nil {
		g.P(`if !(`, variableName, ` > `, fv.GetIntGt(), `) {`)
		p.generateErrorString("validator.int_gt", fv, map[string]string{"int_gt": strconv.FormatInt(fv.GetIntGt(), 10)}, g)
		g.P(`}`)
	}
	if fv.IntLt != nil {
		g.P(`if !(`, variableName, ` < `, fv.GetIntLt(), `) {`)
		p.generateErrorString("validator.int_lt", fv, map[string]string{"int_lt": strconv.FormatInt(fv.GetIntLt(), 10)}, g)
		g.P(`}`)
	}
}

// generateEnumValidator emits a membership check ensuring the value is a valid enum constant.
func (p *plugin) generateEnumValidator(
	field *protogen.Field,
	variableName, ccTypeName, fieldName string,
	fv *validator.FieldValidator, g *protogen.GeneratedFile) {
	if fv.GetIsInEnum() {
		enum := field.Enum
		g.P(`if _, ok := `, enum.GoIdent.GoName, "_name[int32(", variableName, ")]; !ok {")
		p.generateErrorString("validator.is_in_enum", fv, map[string]string{"enum": enum.GoIdent.GoName}, g)
		g.P(`}`)
	}
}

// generateLengthValidator emits length range checks for bytes/string fields.
func (p *plugin) generateLengthValidator(variableName string, fieldName string, fv *validator.FieldValidator, g *protogen.GeneratedFile) {
	if fv.LengthGt != nil {
		g.P(`if !( len(`, variableName, `) > `, fv.GetLengthGt(), `) {`)
		p.generateErrorString("validator.length_gt", fv, map[string]string{"length_gt": strconv.FormatInt(fv.GetLengthGt(), 10)}, g)
		g.P(`}`)
	}

	if fv.LengthLt != nil {
		g.P(`if !( len(`, variableName, `) < `, fv.GetLengthLt(), `) {`)
		p.generateErrorString("validator.length_lt", fv, map[string]string{"length_lt": strconv.FormatInt(fv.GetLengthLt(), 10)}, g)
		g.P(`}`)
	}

	if fv.LengthEq != nil {
		g.P(`if !( len(`, variableName, `) == `, fv.GetLengthEq(), `) {`)
		p.generateErrorString("validator.length_eq", fv, map[string]string{"length_eq": strconv.FormatInt(fv.GetLengthEq(), 10)}, g)
		g.P(`}`)
	}
}

// generateFloatValidator emits float range checks, resolving conflicts between strict and non-strict bounds.
func (p *plugin) generateFloatValidator(variableName string, ccTypeName string, fieldName string, fv *validator.FieldValidator, g *protogen.GeneratedFile) {
	upperIsStrict := true
	lowerIsStrict := true

	// First check for incompatible constraints (i.e flt_lt & flt_lte both defined, etc) and determine the real limits.
	if fv.FloatEpsilon != nil && fv.FloatLt == nil && fv.FloatGt == nil {
		fmt.Fprintf(os.Stderr, "WARNING: field %v.%v has no 'float_lt' or 'float_gt' field so setting 'float_epsilon' has no effect.", ccTypeName, fieldName)
	}
	if fv.FloatLt != nil && fv.FloatLte != nil {
		fmt.Fprintf(os.Stderr, "WARNING: field %v.%v has both 'float_lt' and 'float_lte' constraints, only the strictest will be used.", ccTypeName, fieldName)
		strictLimit := fv.GetFloatLt()
		if fv.FloatEpsilon != nil {
			strictLimit += fv.GetFloatEpsilon()
		}

		if fv.GetFloatLte() < strictLimit {
			upperIsStrict = false
		}
	} else if fv.FloatLte != nil {
		upperIsStrict = false
	}

	if fv.FloatGt != nil && fv.FloatGte != nil {
		fmt.Fprintf(os.Stderr, "WARNING: field %v.%v has both 'float_gt' and 'float_gte' constraints, only the strictest will be used.", ccTypeName, fieldName)
		strictLimit := fv.GetFloatGt()
		if fv.FloatEpsilon != nil {
			strictLimit -= fv.GetFloatEpsilon()
		}

		if fv.GetFloatGte() > strictLimit {
			lowerIsStrict = false
		}
	} else if fv.FloatGte != nil {
		lowerIsStrict = false
	}

	compareStr := ""
	if fv.FloatGt != nil || fv.FloatGte != nil {
		compareStr = fmt.Sprint(`if !(`, variableName)
		if lowerIsStrict {
			if fv.FloatEpsilon != nil {
				compareStr += fmt.Sprint(` + `, fv.GetFloatEpsilon())
			}
			compareStr += fmt.Sprint(` > `, fv.GetFloatGt(), `) {`)
		} else {
			compareStr += fmt.Sprint(` >= `, fv.GetFloatGte(), `) {`)
		}
		g.P(compareStr)
		if lowerIsStrict {
			params := map[string]string{"float_gt": strconv.FormatFloat(fv.GetFloatGt(), 'g', -1, 64)}
			if fv.FloatEpsilon != nil {
				params["float_epsilon"] = strconv.FormatFloat(fv.GetFloatEpsilon(), 'g', -1, 64)
			}
			p.generateErrorString("validator.float_gt", fv, params, g)
		} else {
			p.generateErrorString("validator.float_gte", fv, map[string]string{"float_gte": strconv.FormatFloat(fv.GetFloatGte(), 'g', -1, 64)}, g)
		}
		g.P(`}`)
	}

	if fv.FloatLt != nil || fv.FloatLte != nil {
		compareStr = fmt.Sprint(`if !(`, variableName)
		if upperIsStrict {
			if fv.FloatEpsilon != nil {
				compareStr += fmt.Sprint(` - `, fv.GetFloatEpsilon())
			}
			compareStr += fmt.Sprint(` < `, fv.GetFloatLt(), `) {`)
		} else {
			compareStr += fmt.Sprint(` <= `, fv.GetFloatLte(), `) {`)
		}
		g.P(compareStr)
		if upperIsStrict {
			params := map[string]string{"float_lt": strconv.FormatFloat(fv.GetFloatLt(), 'g', -1, 64)}
			if fv.FloatEpsilon != nil {
				params["float_epsilon"] = strconv.FormatFloat(fv.GetFloatEpsilon(), 'g', -1, 64)
			}
			p.generateErrorString("validator.float_lt", fv, params, g)
		} else {
			p.generateErrorString("validator.float_lte", fv, map[string]string{"float_lte": strconv.FormatFloat(fv.GetFloatLte(), 'g', -1, 64)}, g)
		}
		g.P(`}`)
	}
}

// getUUIDRegex returns a regex to validate that a string is in UUID
// format. The version parameter specified the UUID version. If version is 0,
// the returned regex is valid for any UUID version
func getUUIDRegex(version *int32) (string, error) {
	if version == nil {
		return "", nil
	} else if *version < 0 || *version > 5 {
		return "", fmt.Errorf("UUID version should be between 0-5, Got %d", *version)
	} else if *version == 0 {
		return fmt.Sprintf(uuidPattern, "1-5"), nil
	} else {
		return fmt.Sprintf(uuidPattern, strconv.Itoa(int(*version))), nil
	}
}

// generateStringValidator emits regex, non-empty, and length checks for string fields.
func (p *plugin) generateStringValidator(variableName string, ccTypeName string, fieldName string, fv *validator.FieldValidator, g *protogen.GeneratedFile) {
	if fv.Regex != nil || fv.UuidVer != nil {
		if fv.UuidVer != nil {
			uuid, err := getUUIDRegex(fv.UuidVer)
			if err != nil {
				fmt.Fprintf(os.Stderr, "WARNING: field %v.%v error %s.\n", ccTypeName, fieldName, err)
			} else {
				fv.Regex = &uuid
			}
		}

		g.P(`if !`, p.regexName(ccTypeName, fieldName), `.MatchString(`, variableName, `) {`)
		g.P()
		p.generateErrorString("validator.regex", fv, map[string]string{"regex": fv.GetRegex()}, g)
		g.P()
		g.P(`}`)
	}
	if fv.StringNotEmpty != nil && fv.GetStringNotEmpty() {
		g.P(`if `, variableName, ` == "" {`)
		g.P()
		p.generateErrorString("validator.string_not_empty", fv, nil, g)
		g.P()
		g.P(`}`)
	}
	p.generateLengthValidator(variableName, fieldName, fv, g)
}

// generateRepeatedCountValidator emits min/max element count checks for repeated fields.
func (p *plugin) generateRepeatedCountValidator(variableName string, fieldName string, fv *validator.FieldValidator, g *protogen.GeneratedFile) {
	if fv == nil {
		return
	}
	if fv.RepeatedCountMin != nil {
		compareStr := fmt.Sprint(`if len(`, variableName, `) < `, fv.GetRepeatedCountMin(), ` {`)
		g.P(compareStr)
		g.P()
		p.generateErrorString("validator.repeated_count_min", fv, map[string]string{"repeated_count_min": strconv.FormatInt(fv.GetRepeatedCountMin(), 10)}, g)
		g.P()
		g.P(`}`)
	}
	if fv.RepeatedCountMax != nil {
		compareStr := fmt.Sprint(`if len(`, variableName, `) > `, fv.GetRepeatedCountMax(), ` {`)
		g.P(compareStr)
		g.P()
		p.generateErrorString("validator.repeated_count_max", fv, map[string]string{"repeated_count_max": strconv.FormatInt(fv.GetRepeatedCountMax(), 10)}, g)
		g.P()
		g.P(`}`)
	}
}

// generateErrorString emits *mix.ErrResp: msg 为 i18n key（custom_error 优先），data 为约束占位符。
func (p *plugin) generateErrorString(fallbackKey string, fv *validator.FieldValidator, params map[string]string, g *protogen.GeneratedFile) {
	key := fallbackKey
	if fv != nil && fv.GetCustomError() != "" {
		key = fv.GetCustomError()
	}
	g.P(`return `, p.mixPkg.Ident("NewErrResp"), `(`, p.mixPkg.Ident("InvalidArgument"), `, `, strconv.Quote(key), `, `, p.emitStringMap(params), `)`)
}

func (p *plugin) emitStringMap(params map[string]string) string {
	if len(params) == 0 {
		return "nil"
	}
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var b strings.Builder
	b.WriteString(`map[string]string{`)
	for i, k := range keys {
		if i > 0 {
			b.WriteString(`, `)
		}
		b.WriteString(strconv.Quote(k))
		b.WriteString(`: `)
		b.WriteString(strconv.Quote(params[k]))
	}
	b.WriteString(`}`)
	return b.String()
}

// validatorWithMessageExists reports whether the field validator has the msg_exists constraint set.
func (p *plugin) validatorWithMessageExists(fv *validator.FieldValidator) bool {
	return fv != nil && fv.MsgExists != nil && *(fv.MsgExists)
}

// validatorWithNonRepeatedConstraint reports whether the validator has any constraint other than repeated-count bounds.
func (p *plugin) validatorWithNonRepeatedConstraint(fv *validator.FieldValidator) bool {
	if fv == nil {
		return false
	}

	// Need to use reflection in order to be future-proof for new types of constraints.
	// Elem 而非解引用拷贝：proto 消息内嵌 MessageState 含锁，按值拷贝会触发 vet copylocks
	v := reflect.ValueOf(fv).Elem()
	for i := 0; i < v.NumField(); i++ {
		fieldName := v.Type().Field(i).Name

		// All known validators will have a pointer type and we should skip any fields
		// that are not pointers (i.e unknown fields, etc) as well as 'nil' pointers that
		// don't lead to anything.
		if v.Type().Field(i).Type.Kind() != reflect.Ptr || v.Field(i).IsNil() {
			continue
		}

		// Identify non-repeated constraints based on their name.
		if fieldName != "RepeatedCountMin" && fieldName != "RepeatedCountMax" {
			return true
		}
	}
	return false
}

// regexName constructs the package-level variable name for a compiled regex for the given type and field.
func (p *plugin) regexName(ccTypeName string, fieldName string) string {
	return "regex" + ccTypeName + fieldName
}
