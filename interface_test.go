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

func TestInterfaceGet(t *testing.T) {
	f := newFoo()

	testGet(t, f, "bar.baz", "qux", false)
}
