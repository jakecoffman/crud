package crud

import "testing"

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
	}

	for i, test := range table {
		if v := test.Field.Validate(test.Input); v != test.Expected {
			t.Errorf("%v: For input '%v', expected '%v' got '%v'", i, test.Input, test.Expected, v)
		}
	}
}
