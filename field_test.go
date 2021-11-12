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
	}

	for i, test := range table {
		if v := test.Field.Validate(test.Input); v != test.Expected {
			t.Errorf("%v: For input '%v', expected '%v' got '%v'", i, test.Input, test.Expected, v)
		}
	}
}
