package udn

import (
	"fmt"
	"testing"
)

func testGet(t *testing.T, obj interface{}, key string, val interface{}, isErr bool) {
	actualVal, err := Get(obj, key)
	if !isErr {
		if err != nil {
			t.Errorf("Error fetching %v: %v", key, err)
		}
		if actualVal != val {
			t.Errorf("Wrong value when fetching %v: expected=%v actual=%v", key, val, actualVal)
		}
	} else {
		if err == nil {
			t.Errorf("Found missing key=%v: val=%v err=%v", key, val, err)
		}
	}
}

func TestSimpleMapGet(t *testing.T) {
	m := map[string]map[string]string{
		"foo": map[string]string{"bar": "baz"},
	}

	// Test the ones that are there
	testGet(t, m, "foo.bar", "baz", false)
	// TODO: figure out how to compare these...
	//testGet(t, m, "foo", obj["foo"], false)

	// test some missing ones
	testGet(t, m, "nothere", nil, true)
	testGet(t, m, "foo.qux", nil, true)
}

func TestSimpleStructGet(t *testing.T) {
	type Inner struct {
		Val string
	}
	type Outer struct {
		Child Inner
	}

	s := Outer{Inner{"value"}}

	// Test the ones that are there
	testGet(t, s, "Child.Val", "value", false)
	// TODO: figure out how to compare these...
	testGet(t, s, "Child", s.Child, false)

	// test some missing ones
	testGet(t, s, "nothere", nil, true)
}

func TestSimpleStructPtrGet(t *testing.T) {
	type Inner struct {
		Val string
	}
	type Outer struct {
		Child Inner
	}

	s := &Outer{Inner{"value"}}

	// Test the ones that are there
	testGet(t, s, "Child.Val", "value", false)
	// TODO: figure out how to compare these...
	testGet(t, s, "Child", s.Child, false)

	// test some missing ones
	testGet(t, s, "nothere", nil, true)
}

func TestSimpleArrayGet(t *testing.T) {
	var obj [2][3]string
	obj[0] = [3]string{"a", "b", "c"}
	obj[1] = [3]string{"1", "2", "3"}

	// Test the ones that are there
	for topI, stringSlice := range obj {
		for innerI, val := range stringSlice {
			testGet(t, obj, fmt.Sprintf("%d.%d", topI, innerI), val, false)
		}
	}
	// test some missing ones
	testGet(t, obj, "nothere", nil, true)
}

func TestSimpleSliceGet(t *testing.T) {
	obj := [][]string{}
	obj = append(obj, []string{"a", "b", "c"})
	obj = append(obj, []string{"1", "2", "3"})

	// Test the ones that are there
	for topI, stringSlice := range obj {
		for innerI, val := range stringSlice {
			testGet(t, obj, fmt.Sprintf("%d.%d", topI, innerI), val, false)
		}
	}

	// test some missing ones
	testGet(t, obj, "nothere", nil, true)
}

// Target for benchmarking-- note we don't verify the return here, as we assume
// that the tests cover that
var r interface{}

func benchmarkUDNGet(b *testing.B, base interface{}, key string) {
	var ret interface{}
	for n := 0; n < b.N; n++ {
		ret, _ = Get(base, key)
	}
	r = ret
}

func BenchmarkMapGetUDN(b *testing.B) {
	m := map[string]map[string]string{
		"foo": map[string]string{"bar": "baz"},
	}
	benchmarkUDNGet(b, m, "foo.bar")
}

func BenchmarkMapGet(b *testing.B) {
	m := map[string]map[string]string{
		"foo": map[string]string{"bar": "baz"},
	}
	var ret string
	for n := 0; n < b.N; n++ {
		ret, _ = m["foo"]["bar"]
	}
	r = ret
}

func BenchmarkStructGetUDN(b *testing.B) {
	type Inner struct {
		Val string
	}
	type Outer struct {
		Child Inner
	}

	s := Outer{Inner{"value"}}
	benchmarkUDNGet(b, s, "Child.Val")
}

func BenchmarkStructGet(b *testing.B) {
	type Inner struct {
		Val string
	}
	type Outer struct {
		Child Inner
	}

	s := Outer{Inner{"value"}}
	var ret string
	for n := 0; n < b.N; n++ {
		ret = s.Child.Val
	}
	r = ret
}

func BenchmarkArrayGetUDN(b *testing.B) {
	var obj [2][3]string
	obj[0] = [3]string{"a", "b", "c"}
	obj[1] = [3]string{"1", "2", "3"}

	benchmarkUDNGet(b, obj, "0.1")
}

func BenchmarkArrayGet(b *testing.B) {
	var obj [2][3]string
	obj[0] = [3]string{"a", "b", "c"}
	obj[1] = [3]string{"1", "2", "3"}

	var ret string
	for n := 0; n < b.N; n++ {
		ret = obj[0][1]
	}
	r = ret
}

func BenchmarkSliceGetUDN(b *testing.B) {
	obj := [][]string{}
	obj = append(obj, []string{"a", "b", "c"})
	obj = append(obj, []string{"1", "2", "3"})

	benchmarkUDNGet(b, obj, "0.1")
}

func BenchmarkSliceGet(b *testing.B) {
	obj := [][]string{}
	obj = append(obj, []string{"a", "b", "c"})
	obj = append(obj, []string{"1", "2", "3"})

	var ret string
	for n := 0; n < b.N; n++ {
		ret = obj[0][1]
	}
	r = ret
}
