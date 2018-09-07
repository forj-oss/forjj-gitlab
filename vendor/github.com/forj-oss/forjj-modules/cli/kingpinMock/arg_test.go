package kingpinMock

import (
	"reflect"
	"testing"
)

func TestNewArg(t *testing.T) {
	t.Log("NewArg creates a new Arg with name and help")
	a := NewArg("test", "help")
	if a.name != "test" {
		t.Errorf("name expected to be 'test'. Got %s", a.name)
	}
	if a.help != "help" {
		t.Errorf("help expected to be 'help'. Got %s", a.help)
	}
}

func TestArgClause_Bool(t *testing.T) {
	t.Log("Setting bool type")
	a := NewArg("test", "help")

	if a.GetType() != "any" {
		t.Errorf("Expected arg to be initialized as any. Got %s", a.GetType())
	}

	b := a.Bool()

	bt := reflect.TypeOf(b).String()
	if bt != "*bool" {
		t.Errorf("Expected returned value type to be *bool. Got: %s", bt)
	}

	if a.GetType() != "bool" {
		t.Errorf("Expected arg to be set as bool. Got %s", a.GetType())
	}
}

func TestArgClause_String(t *testing.T) {
	t.Log("Setting string type")
	a := NewArg("test", "help")

	if a.GetType() != "any" {
		t.Errorf("Expected arg to be initialized as any. Got %s", a.GetType())
	}

	b := a.String()

	bt := reflect.TypeOf(b).String()
	if bt != "*string" {
		t.Errorf("Expected returned value type to be *string. Got: %s", bt)
	}

	if a.GetType() != "string" {
		t.Errorf("Expected arg to be set as string. Got : %s", a.GetType())
	}

}

func TestArgClause_Default(t *testing.T) {
	value := "default"
	function := "Default"
	t.Logf("Running %s(\"%s\")", function, value)
	a := NewArg("test", "help")

	if a.vdefault != nil {
		t.Errorf("Expected %s() to not be set. Got '%s'", function, a.vdefault)
	}

	b := a.Default(value)

	if a != b {
		t.Fail()
	}

	if *a.vdefault != value {
		t.Errorf("Expected %s() to be set to '%s'. Got '%s'", function, value, a.vdefault)
	}
}

func TestArgClause_Envar(t *testing.T) {
	value := "ARG"
	function := "Envar"
	t.Logf("Running %s(\"%s\")", function, value)
	a := NewArg("test", "help")
	b := a.Envar("ARG")

	if a != b {
		t.Fail()
	}

	if a.envar != value {
		t.Errorf("Expected %s() to be set to '%s'. Got '%s'", function, value, a.vdefault)
	}

}

func TestArgClause_Required(t *testing.T) {
	value := true
	function := "Required"
	t.Logf("Running %s(\"%s\")", function, value)

	a := NewArg("test", "help")
	b := a.Required()

	if a != b {
		t.Fail()
	}

	if a.required != value {
		t.Errorf("Expected %s() to be true. Got '%s'", function, value, a.required)
	}

}
