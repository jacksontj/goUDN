package udn

import "testing"

/*
// Seem to only get errors like:
//		reflect.Value.Set using unaddressable value [recovered]
//	something to do with the array not having the values directly-- basically
// the index isn't `CanSet()`

func TestSimpleArraySet(t *testing.T) {
	var obj [2][3]string
	obj[0][0] = "a"
	obj[0][1] = "b"
	obj[0][2] = "c"
	obj[1] = [3]string{"1", "2", "3"}

	err := Set(obj, "0.2", "bar")
	if err != nil {
		t.Errorf("Unable to set: %v", err)
	}
	testGet(t, obj, "0.2", "bar", false)
}
*/

func TestSimpleSliceSet(t *testing.T) {
	obj := [][]string{}
	obj = append(obj, []string{"a", "b", "c"})
	obj = append(obj, []string{"1", "2", "3"})

	err := Set(obj, "0.2", "bar")
	if err != nil {
		t.Errorf("Unable to set: %v", err)
	}
	testGet(t, obj, "0.2", "bar", false)

	err = Set(obj, "0.3", "bar")
	if err != nil {
		t.Errorf("Unable to set: %v", err)
	}
	testGet(t, obj, "0.3", "bar", false)

	err = Set(obj, "0.4", "qux")
	if err != nil {
		t.Errorf("Unable to set: %v", err)
	}
	testGet(t, obj, "0.4", "qux", false)
}

func TestSimpleMapSet(t *testing.T) {
	m := map[string]map[string]string{
		"foo": map[string]string{"bar": "baz"},
	}

	err := Set(m, "foo.qux", "bar")
	if err != nil {
		t.Errorf("Unable to set: %v", err)
	}

	testGet(t, m, "foo.qux", "bar", false)
}

/*

TODO: finish, errors like:

	panic: reflect: reflect.Value.Set using unaddressable value [recovered]

This has to do with copying structs where necessary and re-setting


func TestSimpleStructSet(t *testing.T) {
	type Inner struct {
		Val string
	}
	type Outer struct {
		Child Inner
	}

	s := Outer{}
	inner := Inner{"value"}
	err := Set(s, "Child", inner)
	if err != nil {
		t.Errorf("Unable to set: %v", err)
	}

	// Test the ones that are there
	testGet(t, s, "Child.Val", "value", false)
	// TODO: figure out how to compare these...
	testGet(t, s, "Child", s.Child, false)

	// test some missing ones
	testGet(t, s, "nothere", nil, true)
}
*/

func TestSimpleStructPtrSet(t *testing.T) {
	type Inner struct {
		Val string
	}
	type Outer struct {
		Child Inner
	}

	s := &Outer{}
	inner := Inner{"value"}

	err := Set(s, "Child", inner)
	if err != nil {
		t.Errorf("Unable to set: %v", err)
	}
	// Test the ones that are there
	testGet(t, s, "Child.Val", "value", false)
	// TODO: figure out how to compare these...
	testGet(t, s, "Child", s.Child, false)

	// test some missing ones
	testGet(t, s, "nothere", nil, true)
}
