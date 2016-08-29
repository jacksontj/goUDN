package udn

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func Get(base interface{}, key string) (interface{}, error) {
	keyParts := strings.Split(key, ".")
	return GetParts(base, keyParts)
}

func GetParts(base interface{}, keyParts []string) (interface{}, error) {
	var currVal reflect.Value
	currVal = reflect.ValueOf(base)

	// Since some layers might be pointers, we won't range over the keyParts
	// this way the pointer types can just dereference then continue
	for x := 0; x < len(keyParts); {
		keyPart := keyParts[x]
		var subval reflect.Value
		switch currVal.Kind() {
		case reflect.Array:
		case reflect.Slice:
			idx, err := strconv.Atoi(keyPart)
			if err != nil {
				return nil, fmt.Errorf("Accessing an array/slice with a non-int key: %v", keyPart)
			}
			subval = currVal.Index(idx)
		case reflect.Map:
			subval = currVal.MapIndex(reflect.ValueOf(keyPart))
		// if we are a pointer of some kind-- lets dereference the pointer and continue on
		case reflect.Interface:
		case reflect.Ptr:
			currVal = currVal.Elem()
			continue
		case reflect.Struct:
			subval = currVal.FieldByName(keyPart)
		default:
			return nil, fmt.Errorf("Unable to Get() past %v %v", currVal, currVal.Kind())
		}
		if subval.IsValid() {
			currVal = subval
			x++
		} else {
			return nil, fmt.Errorf("unable to find %v in %v", keyPart, currVal)
		}
	}
	return currVal.Interface(), nil
}
