package crud

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestField_Max_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	String().Min(2).Max(1)
}

func TestField_Min_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	String().Max(1).Min(2)
}

func TestField_String(t *testing.T) {
	table := []struct {
		Field    Field
		Input    interface{}
		Expected error
	}{
		{
			Field:    String(),
			Input:    nil,
			Expected: nil,
		},
		{
			Field:    String().Required(),
			Input:    nil,
			Expected: errRequired,
		},
		{
			Field:    String().Required(),
			Input:    7,
			Expected: errWrongType,
		},
		{
			Field:    String().Required(),
			Input:    "",
			Expected: errRequired,
		},
		{
			Field:    String().Required().Allow(""),
			Input:    "",
			Expected: nil,
		},
		{
			Field:    String().Required(),
			Input:    "anything",
			Expected: nil,
		},
		{
			Field:    String().Min(1),
			Input:    "",
			Expected: errMinimum,
		},
		{
			Field:    String().Min(1),
			Input:    "1",
			Expected: nil,
		},
		{
			Field:    String().Max(1),
			Input:    "1",
			Expected: nil,
		},
		{
			Field:    String().Max(1),
			Input:    "12",
			Expected: errMaximum,
		},
		{
			Field:    String().Enum("hi"),
			Input:    "",
			Expected: errEnumNotFound,
		},
		{
			Field:    String().Enum("hi"),
			Input:    "hi",
			Expected: nil,
		},
	}

	for i, test := range table {
		if v := test.Field.Validate(test.Input); v != test.Expected {
			t.Errorf("%v: For input '%v', expected '%v' got '%v'", i, test.Input, test.Expected, v)
		}
	}
}

func TestField_Integer(t *testing.T) {
	table := []struct {
		Field    Field
		Input    interface{}
		Expected error
	}{
		{
			Field:    Integer(),
			Input:    7.2,
			Expected: errWrongType,
		},
		{
			Field:    Integer(),
			Input:    "7",
			Expected: errWrongType,
		},
		{
			Field:    Integer(),
			Input:    7., // Allowed since body will contain float64 with JSON
			Expected: nil,
		},
		{
			Field:    Integer(),
			Input:    nil,
			Expected: nil,
		},
		{
			Field:    Integer().Required(),
			Input:    nil,
			Expected: errRequired,
		},
		{
			Field:    Integer().Required(),
			Input:    0,
			Expected: nil,
		},
		{
			Field:    Integer().Min(1),
			Input:    0,
			Expected: errMinimum,
		},
		{
			Field:    Integer().Min(1),
			Input:    1,
			Expected: nil,
		},
		{
			Field:    Integer().Max(1),
			Input:    1,
			Expected: nil,
		},
		{
			Field:    Integer().Max(1),
			Input:    12,
			Expected: errMaximum,
		},
		{
			Field:    Integer().Enum(1, 2, 3),
			Input:    4,
			Expected: errEnumNotFound,
		},
		{
			Field:    Integer().Enum(1, 2, 3),
			Input:    2,
			Expected: nil,
		},
	}

	for i, test := range table {
		if v := test.Field.Validate(test.Input); v != test.Expected {
			t.Errorf("%v: For input '%v', expected '%v' got '%v'", i, test.Input, test.Expected, v)
		}
	}
}

func TestField_Float64(t *testing.T) {
	table := []struct {
		Field    Field
		Input    interface{}
		Expected error
	}{
		{
			Field:    Number(),
			Input:    "7",
			Expected: errWrongType,
		},
		{
			Field:    Number(),
			Input:    nil,
			Expected: nil,
		},
		{
			Field:    Number().Required(),
			Input:    nil,
			Expected: errRequired,
		},
		{
			Field:    Number().Required(),
			Input:    0.,
			Expected: nil,
		},
		{
			Field:    Number().Min(1.1),
			Input:    0.,
			Expected: errMinimum,
		},
		{
			Field:    Number().Min(1.1),
			Input:    1.1,
			Expected: nil,
		},
		{
			Field:    Number().Max(1.2),
			Input:    1.2,
			Expected: nil,
		},
		{
			Field:    Number().Max(1),
			Input:    12.,
			Expected: errMaximum,
		},
		{
			Field:    Number().Enum(1., 2.),
			Input:    3.,
			Expected: errEnumNotFound,
		},
		{
			Field:    Number().Enum(1., 2.),
			Input:    1.,
			Expected: nil,
		},
	}

	for i, test := range table {
		if v := test.Field.Validate(test.Input); v != test.Expected {
			t.Errorf("%v: For input '%v', expected '%v' got '%v'", i, test.Input, test.Expected, v)
		}
	}
}

func TestField_Boolean(t *testing.T) {
	table := []struct {
		Field    Field
		Input    interface{}
		Expected error
	}{
		{
			Field:    Boolean(),
			Input:    "true",
			Expected: errWrongType,
		},
		{
			Field:    Boolean(),
			Input:    nil,
			Expected: nil,
		},
		{
			Field:    Boolean().Required(),
			Input:    nil,
			Expected: errRequired,
		},
		{
			Field:    Boolean().Required(),
			Input:    true,
			Expected: nil,
		},
		{
			Field:    Boolean(),
			Input:    false,
			Expected: nil,
		},
		{
			Field:    Boolean().Enum(true),
			Input:    false,
			Expected: errEnumNotFound,
		},
		{
			Field:    Boolean().Enum(true),
			Input:    true,
			Expected: nil,
		},
	}

	for i, test := range table {
		if v := test.Field.Validate(test.Input); v != test.Expected {
			t.Errorf("%v: For input '%v', expected '%v' got '%v'", i, test.Input, test.Expected, v)
		}
	}
}

func TestField_Array(t *testing.T) {
	table := []struct {
		Field    Field
		Input    interface{}
		Expected error
	}{
		{
			Field:    Array(),
			Input:    true,
			Expected: errWrongType,
		},
		{
			Field:    Array(),
			Input:    nil,
			Expected: nil,
		},
		{
			Field:    Array().Required(),
			Input:    nil,
			Expected: errRequired,
		},
		{
			Field:    Array().Required(),
			Input:    []interface{}{},
			Expected: nil,
		},
		{
			Field:    Array(),
			Input:    []interface{}{1},
			Expected: nil,
		},
		{
			Field:    Array().Min(1),
			Input:    []interface{}{},
			Expected: errMinimum,
		},
		{
			Field:    Array().Min(1),
			Input:    []interface{}{1},
			Expected: nil,
		},
		{
			Field:    Array().Max(1),
			Input:    []interface{}{1, 2},
			Expected: errMaximum,
		},
		{
			Field:    Array().Max(1),
			Input:    []interface{}{1},
			Expected: nil,
		},
		{
			Field:    Array().Items(String()),
			Input:    []interface{}{},
			Expected: nil,
		},
		{
			Field:    Array().Items(String().Required()),
			Input:    []interface{}{""},
			Expected: errRequired,
		},
		{
			Field:    Array().Items(String().Required()),
			Input:    []interface{}{"hi"},
			Expected: nil,
		},
	}

	for i, test := range table {
		if v := test.Field.Validate(test.Input); v != test.Expected {
			t.Errorf("%v: For input '%v', expected '%v' got '%v'", i, test.Input, test.Expected, v)
		}
	}
}

func TestField_Object(t *testing.T) {
	table := []struct {
		Field    Field
		Input    interface{}
		Expected error
	}{
		{
			Field:    Object(map[string]Field{}),
			Input:    "",
			Expected: errWrongType,
		},
		{
			Field:    Object(map[string]Field{}),
			Input:    nil,
			Expected: nil,
		},
		{
			Field:    Object(map[string]Field{}).Required(),
			Input:    nil,
			Expected: errRequired,
		},
		{
			Field:    Object(map[string]Field{}).Required(),
			Input:    map[string]interface{}{},
			Expected: nil,
		},
		{
			Field: Object(map[string]Field{
				"nested": Integer().Required().Min(1),
			}),
			Input:    map[string]interface{}{},
			Expected: errRequired,
		},
		{
			Field: Object(map[string]Field{
				"nested": Integer().Required().Min(1),
			}),
			Input: map[string]interface{}{
				"nested": 1.1,
			},
			Expected: errWrongType,
		},
		{
			Field: Object(map[string]Field{
				"nested": Integer().Required().Min(1),
			}),
			Input: map[string]interface{}{
				"nested": 1.,
			},
			Expected: nil,
		},
	}

	for i, test := range table {
		if v := test.Field.Validate(test.Input); !errors.Is(v, test.Expected) {
			t.Errorf("%v: For input '%v', expected '%v' got '%v'", i, test.Input, test.Expected, v)
		}
	}
}

// Tests children inheriting their parent's settings
func TestField_Object_Setting_Inheritance(t *testing.T) {
	obj := Object(map[string]Field{
		"child": Object(map[string]Field{}),
	}).Strip(false)

	input := map[string]interface{}{
		"child": map[string]interface{}{
			"grandchild": true,
		},
		"another": "hi",
	}

	err := obj.Validate(input)
	if err != nil {
		t.Errorf(err.Error())
	}

	if v, ok := input["another"]; !ok {
		t.Errorf("another is missing")
	} else {
		if v != "hi" {
			t.Error("expected hi got", v)
		}
	}

	if v, ok := input["child"]; !ok {
		t.Errorf("child missing")
	} else {
		child := v.(map[string]interface{})
		if v, ok := child["grandchild"]; !ok {
			t.Errorf("grandchild missing")
		} else if v != true {
			t.Error("grandchild expected true, got", v)
		}
	}

	obj = Object(map[string]Field{
		"child": Object(map[string]Field{}).Strip(true),
	}).Strip(false)

	err = obj.Validate(input)
	if err != nil {
		t.Errorf(err.Error())
	}

	if v, ok := input["another"]; !ok {
		t.Errorf("another is missing")
	} else {
		if v != "hi" {
			t.Error("expected hi got", v)
		}
	}

	if v, ok := input["child"]; !ok {
		t.Errorf("child missing")
	} else {
		child := v.(map[string]interface{})
		if _, ok := child["grandchild"]; ok {
			t.Errorf("expected grandchild to be stripped, but it still exists")
		}
	}
}

// Tests children inheriting their parent's settings
func TestField_Array_Setting_Inheritance(t *testing.T) {
	obj := Array().Items(Object(map[string]Field{})).Strip(false)

	input := []interface{}{
		map[string]interface{}{
			"hello": "world",
		},
	}

	err := obj.Validate(input)
	if err != nil {
		t.Errorf(err.Error())
	}

	if fmt.Sprint(input) != "[map[hello:world]]" {
		t.Errorf(fmt.Sprint(input))
	}

	obj = Array().Items(Object(map[string]Field{}).Strip(true)).Strip(false)

	err = obj.Validate(input)
	if err != nil {
		t.Errorf(err.Error())
	}

	if fmt.Sprint(input) != "[map[]]" {
		t.Errorf(fmt.Sprint(input))
	}
}

func TestField_Validate_DateTime(t *testing.T) {
	dt := DateTime()
	err := dt.Validate("q")
	switch v := err.(type) {
	case *time.ParseError:
	default:
		t.Errorf("Expected a ParseError, got %s", v)
	}

	if err = dt.Validate(time.Now().Format(time.RFC3339)); err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
}

func TestField_Validate_Date(t *testing.T) {
	dt := Date()
	err := dt.Validate("q")
	switch v := err.(type) {
	case *time.ParseError:
	default:
		t.Errorf("Expected a ParseError, got %s", v)
	}

	if err = dt.Validate(time.Now().Format(fullDate)); err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
}
