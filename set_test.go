package udn

import "testing"

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
