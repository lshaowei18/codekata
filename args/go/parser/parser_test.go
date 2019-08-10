package parser

import (
	"reflect"
	"strconv"
	"testing"
)

func TestParser(t *testing.T) {
	flags := "l,p#,d*"
	args := []string{"./parser.go", "-l", "-p", "8080", "-d", "/usr/logs"}
	parser, _ := Args(flags, args)

	tests := []struct {
		flag string
		want interface{}
	}{
		{"l", false},
		{"p", 8080},
		{"d", "/usr/logs"},
	}
	for _, tt := range tests {
		got, _ := parser.GetValue(tt.flag)
		if got != tt.want {
			t.Errorf("Wanted :%v, got: %v ", tt.want, got)
		}
	}
}

func TestArgs(t *testing.T) {
	flags := "l,p,d"
	args := []string{"go run", "-l", "-p"}
	got, _ := Args(flags, args)

	if reflect.ValueOf(got).Kind() != reflect.Struct {
		t.Errorf("Args func should return a struct type.")
	}

	flagsSplit := []string{"l", "p", "d"}

	if !reflect.DeepEqual(flagsSplit, got.flags) {
		t.Errorf("Wanted : %v, got : %v", flagsSplit, got.flags)
	}

	if !reflect.DeepEqual(got.args, args[1:]) {
		t.Errorf("Wanted: %s, got : %s", args, got.args)
	}

	if reflect.ValueOf(got.schema).Kind() != reflect.Map {
		t.Errorf("Args should return a struct where its schema is a map.")
	}
}

func TestCreateDefaultSchema(t *testing.T) {
	flags := "l,p*,d#"
	args := []string{"go run", "-l", "-p"}

	parser, _ := Args(flags, args)

	parser.createSchema()

	want := map[string]interface{}{
		"l": false,
		"p": "",
		"d": 0,
	}

	if !reflect.DeepEqual(parser.schema, want) {
		t.Errorf("Wanted: %v, got: %v", parser.schema, want)
	}
}

func TestInvalidSchema(t *testing.T) {
	flags := "l,p*,d#,g+"
	args := []string{"go run"}

	_, err := Args(flags, args)
	if err == nil {
		t.Errorf("Should be an error, invalid flags given %v", flags)
	}
}

func TestInvalidParseArgs(t *testing.T) {
	flags := "l,p*,d#"
	args := []string{"go run", "-l", "-p", "-a", "Hello"}

	_, err := Args(flags, args)
	if err == nil {
		t.Errorf("Should have an error as the flags are invalid.")
	}
}

func TestParseArgsWithCorrectValues(t *testing.T) {
	flags := "l,p*,d#"
	args := []string{"go run", "-l", "true", "-p", "hello", "-d", "123"}

	parser, _ := Args(flags, args)
	want := map[string]interface{}{
		"l": true,
		"p": "hello",
		"d": 123,
	}
	if !reflect.DeepEqual(parser.schema, want) {
		parser.printSchema()
		t.Errorf("Wanted: %v, got: %v", want, parser.schema)
	}
}

func TestInvalidArgumentTypes(t *testing.T) {
	tests := []struct {
		flags string
		args  []string
	}{
		{"l,p*,d#", []string{"go run", "-l", "10"}},
		{"l,p*,d#", []string{"go run", "-d", "hello"}},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			_, err := Args(tt.flags, tt.args)
			if err == nil {
				t.Errorf("Should be an error, invalid args given. flags: %v, args :%v \n",
					tt.flags, tt.args)
			}
		})
	}
}

func TestGetValueWhenThereIsNoArg(t *testing.T) {
	schema := "l,p#"
	args := []string{"go run"}

	parser, _ := Args(schema, args)
	got, _ := parser.GetValue("l")

	if got != false {
		t.Errorf("Default value for boolean arg should be false, got : %v", got)
	}

	got, _ = parser.GetValue("p")

	if got != 0 {
		t.Errorf("Default value for string arg should be empty. got : %v", got)
	}
}

func TestGetValue(t *testing.T) {
	tests := []struct {
		flag string
		want interface{}
	}{
		{"l", true},
		{"p", 8080},
		{"d", ""},
	}
	port := 8080
	schema := "l,p#,d*"
	strPort := strconv.Itoa(port)
	args := []string{"go run", "-p", strPort, "-l", "true"}

	parser, _ := Args(schema, args)
	for _, tt := range tests {
		got, _ := parser.GetValue(tt.flag)

		if got != tt.want {
			t.Errorf("wanted: %v of type %T, got: %v of type %T", tt.want, tt.want,
				got, got)
		}
	}
}

func TestGetKeyAndFlagType(t *testing.T) {
	tests := []struct {
		flag     string
		key      string
		flagType reflect.Kind
	}{
		{"p#", "p", reflect.Int},
		{"d*", "d", reflect.String},
		{"l", "l", reflect.Bool},
	}

	for _, tt := range tests {
		t.Run("Test parsing of flags", func(t *testing.T) {
			key, flagType := parseFlag(tt.flag)

			if key != tt.key {
				t.Errorf("wanted: %v, got: %v", tt.key, key)
			}

			if flagType != tt.flagType {
				t.Errorf("wanted: %v, got: %v", tt.flagType, flagType)
			}
		})
	}
}

func TestGetTypeFromSchema(t *testing.T) {
	tests := []struct {
		flag     string
		flagType reflect.Kind
	}{
		{"p", reflect.Int},
		{"d", reflect.String},
		{"l", reflect.Bool},
		{"g", reflect.Invalid},
	}

	schema := "l,p#,d*"
	args := []string{"go run"}

	parser, _ := Args(schema, args)

	for _, tt := range tests {
		t.Run("Test parsing of flags", func(t *testing.T) {
			flagType, _ := parser.getFlagType(tt.flag)

			if flagType != tt.flagType {
				t.Errorf("wanted: %v, got: %v", tt.flagType, flagType)
			}
		})
	}
}

func TestConvertStringToType(t *testing.T) {
	t.Run("Valid values", func(t *testing.T) {
		tests := []struct {
			value    string
			flagType reflect.Kind
			want     interface{}
		}{
			{"1", reflect.Int, 1},
			{"hi", reflect.String, "hi"},
			{"true", reflect.Bool, true},
		}

		for _, tt := range tests {
			got, _ := convertStringToType(tt.value, tt.flagType)

			if got != tt.want {
				t.Errorf("wanted: %v, got: %v", tt.want, got)
			}
		}
	})

	t.Run("Invalid values", func(t *testing.T) {
		tests := []struct {
			value    string
			flagType reflect.Kind
		}{
			{"hello", reflect.Int},
			{"321", reflect.Bool},
		}

		for _, tt := range tests {
			_, err := convertStringToType(tt.value, tt.flagType)

			if err == nil {
				t.Errorf("Value type : %v, and flag : %v don't match, should have err",
					tt.value, tt.flagType)
			}
		}
	})
}
