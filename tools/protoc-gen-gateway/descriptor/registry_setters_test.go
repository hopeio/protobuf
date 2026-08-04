package descriptor

import "testing"

func TestRegistrySettersAndGetters(t *testing.T) {
	reg := NewRegistry()

	reg.SetAllowMerge(true)
	if !reg.IsAllowMerge() {
		t.Fatal("IsAllowMerge want true")
	}
	reg.SetAllowMerge(false)
	if reg.IsAllowMerge() {
		t.Fatal("IsAllowMerge want false")
	}

	reg.SetAllowRepeatedFieldsInBody(true)
	if !reg.IsAllowRepeatedFieldsInBody() {
		t.Fatal("IsAllowRepeatedFieldsInBody want true")
	}
	reg.SetAllowRepeatedFieldsInBody(false)
	if reg.IsAllowRepeatedFieldsInBody() {
		t.Fatal("IsAllowRepeatedFieldsInBody want false")
	}

	reg.SetIncludePackageInTags(true)
	if !reg.IsIncludePackageInTags() {
		t.Fatal("IsIncludePackageInTags want true")
	}
	reg.SetIncludePackageInTags(false)
	if reg.IsIncludePackageInTags() {
		t.Fatal("IsIncludePackageInTags want false")
	}

	reg.SetMergeFileName("merged.swagger.json")
	if got := reg.GetMergeFileName(); got != "merged.swagger.json" {
		t.Fatalf("GetMergeFileName=%q", got)
	}

	reg.SetAllowDeleteBody(true)

	reg.SetUseJSONNamesForFields(true)
	if !reg.GetUseJSONNamesForFields() {
		t.Fatal("GetUseJSONNamesForFields want true")
	}

	reg.SetUseFQNForOpenAPIName(true)
	if !reg.GetUseFQNForOpenAPIName() {
		t.Fatal("GetUseFQNForOpenAPIName want true")
	}

	reg.SetUseGoTemplate(true)
	if !reg.GetUseGoTemplate() {
		t.Fatal("GetUseGoTemplate want true")
	}

	reg.SetEnumsAsInts(true)
	if !reg.GetEnumsAsInts() {
		t.Fatal("GetEnumsAsInts want true")
	}

	reg.SetDisableDefaultErrors(true)
	if !reg.GetDisableDefaultErrors() {
		t.Fatal("GetDisableDefaultErrors want true")
	}

	reg.SetSimpleOperationIDs(true)
	if !reg.GetSimpleOperationIDs() {
		t.Fatal("GetSimpleOperationIDs want true")
	}

	reg.SetOmitPackageDoc(true)
	if !reg.GetOmitPackageDoc() {
		t.Fatal("GetOmitPackageDoc want true")
	}

	if err := reg.SetRepeatedPathParamSeparator("csv"); err != nil {
		t.Fatal(err)
	}
	if got := reg.GetRepeatedPathParamSeparator(); got != ',' {
		t.Fatalf("GetRepeatedPathParamSeparator=%q want ','", got)
	}
	if got := reg.GetRepeatedPathParamSeparatorName(); got != "csv" {
		t.Fatalf("GetRepeatedPathParamSeparatorName=%q", got)
	}
	if err := reg.SetRepeatedPathParamSeparator("unknown"); err == nil {
		t.Fatal("want error for unknown separator")
	}
}
