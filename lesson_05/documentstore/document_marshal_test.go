package documentstore

import (
	"reflect"
	"testing"
)

type simpleStruct struct {
	X int
	Y string
	Z bool
}

type innerStruct struct {
	A int
	B string
}

type outerStruct struct {
	Name  string
	Inner innerStruct
	Nums  []int
}

func TestMarshalSimpleStruct(t *testing.T) {
	in := simpleStruct{
		X: 42,
		Y: "hello",
		Z: true,
	}

	doc, err := MarshalDocument(in)
	if err != nil {
		t.Fatalf("MarshalDocument(simpleStruct) error = %v", err)
	}

	if doc == nil {
		t.Fatalf("MarshalDocument(simpleStruct) returned nil doc")
	}

	if len(doc.Fields) != 3 {
		t.Fatalf("expected 3 fields, got %d", len(doc.Fields))
	}

	// X
	fx, ok := doc.Fields["X"]
	if !ok {
		t.Fatalf("field X not found in document")
	}
	if fx.Type != DocumentFieldTypeNumber {
		t.Fatalf("field X type = %s, want %s", fx.Type, DocumentFieldTypeNumber)
	}
	if v, ok := fx.Value.(int); !ok || v != 42 {
		t.Fatalf("field X value = %#v (%T), want 42 (int)", fx.Value, fx.Value)
	}

	// Y
	fy, ok := doc.Fields["Y"]
	if !ok {
		t.Fatalf("field Y not found in document")
	}
	if fy.Type != DocumentFieldTypeString {
		t.Fatalf("field Y type = %s, want %s", fy.Type, DocumentFieldTypeString)
	}
	if v, ok := fy.Value.(string); !ok || v != "hello" {
		t.Fatalf("field Y value = %#v (%T), want \"hello\" (string)", fy.Value, fy.Value)
	}

	// Z
	fz, ok := doc.Fields["Z"]
	if !ok {
		t.Fatalf("field Z not found in document")
	}
	if fz.Type != DocumentFieldTypeBool {
		t.Fatalf("field Z type = %s, want %s", fz.Type, DocumentFieldTypeBool)
	}
	if v, ok := fz.Value.(bool); !ok || v != true {
		t.Fatalf("field Z value = %#v (%T), want true (bool)", fz.Value, fz.Value)
	}
}

func TestMarshalNestedAndArray(t *testing.T) {
	in := outerStruct{
		Name: "outer",
		Inner: innerStruct{
			A: 5,
			B: "inner",
		},
		Nums: []int{1, 2, 3},
	}

	doc, err := MarshalDocument(&in) // pointer should also work
	if err != nil {
		t.Fatalf("MarshalDocument(outerStruct) error = %v", err)
	}

	// Name
	fName, ok := doc.Fields["Name"]
	if !ok {
		t.Fatalf("field Name not found")
	}
	if fName.Type != DocumentFieldTypeString {
		t.Fatalf("Name type = %s, want %s", fName.Type, DocumentFieldTypeString)
	}

	// Inner (object)
	fInner, ok := doc.Fields["Inner"]
	if !ok {
		t.Fatalf("field Inner not found")
	}
	if fInner.Type != DocumentFieldTypeObject {
		t.Fatalf("Inner type = %s, want %s", fInner.Type, DocumentFieldTypeObject)
	}
	nested, ok := fInner.Value.(*Document)
	if !ok || nested == nil {
		t.Fatalf("Inner value = %#v (%T), want *Document", fInner.Value, fInner.Value)
	}
	if fa, ok := nested.Fields["A"]; !ok || fa.Type != DocumentFieldTypeNumber {
		t.Fatalf("nested field A missing or wrong type")
	}

	// Nums (array)
	fNums, ok := doc.Fields["Nums"]
	if !ok {
		t.Fatalf("field Nums not found")
	}
	if fNums.Type != DocumentFieldTypeArray {
		t.Fatalf("Nums type = %s, want %s", fNums.Type, DocumentFieldTypeArray)
	}
	items, ok := fNums.Value.([]DocumentField)
	if !ok {
		t.Fatalf("Nums value = %#v (%T), want []DocumentField", fNums.Value, fNums.Value)
	}
	if len(items) != 3 {
		t.Fatalf("Nums length = %d, want 3", len(items))
	}
	for i, it := range items {
		if it.Type != DocumentFieldTypeNumber {
			t.Fatalf("Nums[%d].Type = %s, want %s", i, it.Type, DocumentFieldTypeNumber)
		}
		// we don't assert exact numeric type, just convert via reflect
		rv := reflect.ValueOf(it.Value)
		if !rv.IsValid() {
			t.Fatalf("Nums[%d].Value invalid", i)
		}
		got := rv.Convert(reflect.TypeOf(int(0))).Int()
		if int(got) != i+1 {
			t.Fatalf("Nums[%d] = %v, want %d", i, it.Value, i+1)
		}
	}
}

func TestMarshalErrors(t *testing.T) {
	// nil input
	if _, err := MarshalDocument(nil); err == nil {
		t.Fatalf("expected error for nil input, got nil")
	}

	// non-struct input
	if _, err := MarshalDocument(123); err == nil {
		t.Fatalf("expected error for non-struct input, got nil")
	}

	// nil pointer to struct
	var ps *simpleStruct
	if _, err := MarshalDocument(ps); err == nil {
		t.Fatalf("expected error for nil pointer input, got nil")
	}
}

// ---------- Unmarshal tests ----------

func TestUnmarshalSimpleDocument(t *testing.T) {
	doc := &Document{
		Fields: map[string]DocumentField{
			"X": {Type: DocumentFieldTypeNumber, Value: int64(10)},
			"Y": {Type: DocumentFieldTypeString, Value: "test"},
			"Z": {Type: DocumentFieldTypeBool, Value: true},
		},
	}

	var out simpleStruct
	if err := UnmarshalDocument(doc, &out); err != nil {
		t.Fatalf("UnmarshalDocument error = %v", err)
	}

	if out.X != 10 || out.Y != "test" || out.Z != true {
		t.Fatalf("unexpected result: %#v", out)
	}
}

func TestUnmarshalNestedAndArray(t *testing.T) {
	doc := &Document{
		Fields: map[string]DocumentField{
			"Name": {Type: DocumentFieldTypeString, Value: "outer"},
			"Inner": {
				Type: DocumentFieldTypeObject,
				Value: &Document{
					Fields: map[string]DocumentField{
						"A": {Type: DocumentFieldTypeNumber, Value: int64(7)},
						"B": {Type: DocumentFieldTypeString, Value: "inner"},
					},
				},
			},
			"Nums": {
				Type: DocumentFieldTypeArray,
				Value: []DocumentField{
					{Type: DocumentFieldTypeNumber, Value: int64(1)},
					{Type: DocumentFieldTypeNumber, Value: int64(2)},
					{Type: DocumentFieldTypeNumber, Value: int64(3)},
				},
			},
		},
	}

	var out outerStruct
	if err := UnmarshalDocument(doc, &out); err != nil {
		t.Fatalf("UnmarshalDocument error = %v", err)
	}

	if out.Name != "outer" {
		t.Fatalf("Name = %q, want %q", out.Name, "outer")
	}
	if out.Inner.A != 7 || out.Inner.B != "inner" {
		t.Fatalf("Inner = %#v, want {A:7 B:\"inner\"}", out.Inner)
	}
	if !reflect.DeepEqual(out.Nums, []int{1, 2, 3}) {
		t.Fatalf("Nums = %#v, want []int{1,2,3}", out.Nums)
	}
}

func TestUnmarshalErrors(t *testing.T) {
	doc := &Document{Fields: map[string]DocumentField{}}

	// nil doc
	var out simpleStruct
	if err := UnmarshalDocument(nil, &out); err == nil {
		t.Fatalf("expected error for nil doc, got nil")
	}

	// nil output
	if err := UnmarshalDocument(doc, nil); err == nil {
		t.Fatalf("expected error for nil output, got nil")
	}

	// non-pointer output
	if err := UnmarshalDocument(doc, out); err == nil {
		t.Fatalf("expected error for non-pointer output, got nil")
	}

	// pointer to non-struct
	var i int
	if err := UnmarshalDocument(doc, &i); err == nil {
		t.Fatalf("expected error for pointer to non-struct, got nil")
	}

	// type mismatch: doc has string, struct expects int
	doc2 := &Document{
		Fields: map[string]DocumentField{
			"X": {Type: DocumentFieldTypeString, Value: "not-a-number"},
		},
	}
	var out2 simpleStruct
	if err := UnmarshalDocument(doc2, &out2); err == nil {
		t.Fatalf("expected error for type mismatch (string into int), got nil")
	}
}

func TestRoundTrip(t *testing.T) {
	orig := outerStruct{
		Name: "roundtrip",
		Inner: innerStruct{
			A: 100,
			B: "nested",
		},
		Nums: []int{10, 20, 30},
	}

	doc, err := MarshalDocument(orig)
	if err != nil {
		t.Fatalf("MarshalDocument error = %v", err)
	}

	var decoded outerStruct
	if err := UnmarshalDocument(doc, &decoded); err != nil {
		t.Fatalf("UnmarshalDocument error = %v", err)
	}

	if !reflect.DeepEqual(orig, decoded) {
		t.Fatalf("round-trip mismatch:\n  orig   = %#v\n  decoded= %#v", orig, decoded)
	}
}