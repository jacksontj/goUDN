package udn

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func Get(base interface{}, key string) (interface{}, error) {
	keyParts := strings.Split(key, ".")
	val, _, err := getSetParts(base, keyParts, nil, false)
	if err != nil {
		return nil, err
	} else {
		return val.Interface(), nil
	}
}

func getSetParts(base interface{}, keyParts []string, val interface{}, set bool) (*reflect.Value, bool, error) {
	var currVal reflect.Value
	currVal = reflect.ValueOf(base)

	var subval reflect.Value
	// Since some layers might be pointers, we won't range over the keyParts
	// this way the pointer types can just dereference then continue
	for x := 0; x < len(keyParts); {
		// check to see if this layer implements the setter interface
		if set && currVal.Type().Implements(setInterface) {
			tmp := currVal.Interface().(UDNSetter)
			err := tmp.SetUDN(keyParts[x:], val)
			if err == nil {
				return nil, true, nil
			}
		}

		// check to see if this layer implements the getter
		// if so, call that interface-- then continue on
		if currVal.Type().Implements(getInterface) {
			tmp := currVal.Interface().(UDNGetter)
			tmpVal, incrX, err := tmp.GetUDN(keyParts[x:])
			if err == nil {
				x = x + incrX
				currVal = reflect.ValueOf(tmpVal)
				continue
			}
		}

		keyPart := keyParts[x]
		switch currVal.Kind() {
		case reflect.Array:
			fallthrough
		case reflect.Slice:
			idx, err := strconv.Atoi(keyPart)
			if err != nil {
				return nil, false, fmt.Errorf("Accessing an array/slice with a non-int key: %v", keyPart)
			}
			subval = currVal.Index(idx)
		case reflect.Map:
			subval = currVal.MapIndex(reflect.ValueOf(keyPart))
		// if we are a pointer of some kind-- lets dereference the pointer and continue on
		case reflect.Ptr:
			currVal = currVal.Elem()
			continue
		case reflect.Struct:
			subval = currVal.FieldByName(keyPart)
		default:
			return nil, false, fmt.Errorf("Unable to Get() past %v %v", currVal, currVal.Kind())
		}
		if subval.IsValid() {
			currVal = subval
			x++
		} else {
			return nil, false, fmt.Errorf("unable to find %v in %v %v: %v", keyPart, currVal, currVal.Kind(), subval)
		}
	}
	return &currVal, false, nil
}
