package udn

import (
	"fmt"
	"testing"
)

func newFoo() *Foo {
	return &Foo{
		MapOStuff: map[string]string{"baz": "qux"},
	}
}

type Foo struct {
	MapOStuff map[string]string
}

func (f *Foo) GetUDN(keyParts []string) (interface{}, int, error) {
	if len(keyParts) > 0 {
		switch keyParts[0] {
		case "bar":
			return f.MapOStuff, 1, nil
		}
	}

	return nil, 0, fmt.Errorf("not found")
}

func (f *Foo) SetUDN(keyParts []string, val interface{}) error {
	if len(keyParts) == 2 {
		if keyParts[0] == "bar" && keyParts[1] == "foo" {
			f.MapOStuff["foo"] = val.(string)
			return nil
		}
	}
	return fmt.Errorf("nope")
}

func TestInterfaceGet(t *testing.T) {
	f := newFoo()

	testGet(t, f, "bar.baz", "qux", false)
}

func TestInterfaceSet(t *testing.T) {
	f := newFoo()

	testGet(t, f, "bar.foo", "fooqux", true)

	// make sure that pass-by-value errors out
	err := Set(&f, "bar.foo", "fooqux")
	if err != nil {
		t.Errorf("Unable to set: %v", err)
	}

	testGet(t, f, "bar.foo", "fooqux", false)
}
