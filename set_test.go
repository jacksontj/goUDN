package udn

import "testing"

func TestSimpleArraySet(t *testing.T) {
	var obj [2][3]string
	obj[0][0] = "a"
	obj[0][1] = "b"
	obj[0][2] = "c"
	obj[1] = [3]string{"1", "2", "3"}

	// make sure that pass-by-value errors out
	err := Set(obj, "0.2", "bar")
	if err == nil {
		t.Errorf("No error when passing-by-value!")
	}

	// ensure that pass-by-value works
	err = Set(&obj, "0.2", "bar")
	if err != nil {
		t.Errorf("Unable to set: %v", err)
	}
	testGet(t, obj, "0.2", "bar", false)
}

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
	if err == nil {
		t.Errorf("No error when pass-by-value")
	}

	err = Set(&s, "Child", inner)
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

func TestFoo(t *testing.T) {
	type Inner struct {
		Val string
	}
	type Outer struct {
		Child *Inner
	}

	s := Outer{&Inner{"value"}}

	err := Set(s, "Child.Val", "bar")
	if err != nil {
		t.Errorf("Unable to set: %v", err)
	}
	testGet(t, s, "Child.Val", "bar", false)
}

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

// Target for benchmarking-- note we don't verify the return here, as we assume
// that the tests cover that
func benchmarkUDNSet(b *testing.B, base interface{}, key string, val interface{}) {
	for n := 0; n < b.N; n++ {
		Set(base, key, val)
	}
}

func BenchmarkMapSetUDN(b *testing.B) {
	m := map[string]map[string]string{
		"foo": map[string]string{"bar": "baz"},
	}
	benchmarkUDNSet(b, m, "foo.bar", "somethingelse")
}

func BenchmarkMapSet(b *testing.B) {
	m := map[string]map[string]string{
		"foo": map[string]string{"bar": "baz"},
	}
	for n := 0; n < b.N; n++ {
		m["foo"]["bar"] = "somethingelse"
	}
}

func BenchmarkStructSetUDN(b *testing.B) {
	type Inner struct {
		Val string
	}
	type Outer struct {
		Child Inner
	}

	s := Outer{Inner{"value"}}
	benchmarkUDNSet(b, s, "Child.Val", "somethingelse")
}

func BenchmarkStructSet(b *testing.B) {
	type Inner struct {
		Val string
	}
	type Outer struct {
		Child Inner
	}

	s := Outer{Inner{"value"}}
	for n := 0; n < b.N; n++ {
		s.Child.Val = "somethingelse"
	}
}

func BenchmarkArraySetUDN(b *testing.B) {
	var obj [2][3]string
	obj[0] = [3]string{"a", "b", "c"}
	obj[1] = [3]string{"1", "2", "3"}

	benchmarkUDNSet(b, obj, "0.1", "somethingelse")
}

func BenchmarkArraySet(b *testing.B) {
	var obj [2][3]string
	obj[0] = [3]string{"a", "b", "c"}
	obj[1] = [3]string{"1", "2", "3"}

	for n := 0; n < b.N; n++ {
		obj[0][1] = "somethingelse"
	}
}

func BenchmarkSliceSetUDN(b *testing.B) {
	obj := [][]string{}
	obj = append(obj, []string{"a", "b", "c"})
	obj = append(obj, []string{"1", "2", "3"})

	benchmarkUDNSet(b, obj, "0.1", "somethingelse")
}

func BenchmarkSliceSet(b *testing.B) {
	obj := [][]string{}
	obj = append(obj, []string{"a", "b", "c"})
	obj = append(obj, []string{"1", "2", "3"})

	for n := 0; n < b.N; n++ {
		obj[0][1] = "somethingelse"
	}
}
