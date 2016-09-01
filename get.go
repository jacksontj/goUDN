package udn

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func Get(base interface{}, key string) (interface{}, error) {
	keyParts := strings.Split(key, ".")
	val, err := GetParts(base, keyParts)
	if err != nil {
		return nil, err
	} else {
		return val.Interface(), nil
	}
}

// map of object ptr -> key (name) -> value
var reflectCache map[interface{}]map[string]reflect.Value

func GetParts(base interface{}, keyParts []string) (*reflect.Value, error) {
	var currVal reflect.Value
	currVal = reflect.ValueOf(base)

	// This is a very bad attempt at dealing with the caching.... this seems
	// to be no faster than just doing it at runtime
	var baseCacheMap map[string]reflect.Value
	var ok bool
	if currVal.Kind() == reflect.Ptr {
		// TODO: initialize elsewhere
		if reflectCache == nil {
			reflectCache = make(map[interface{}]map[string]reflect.Value)
		}
		baseCacheMap, ok = reflectCache[base]
		if !ok {
			baseCacheMap = make(map[string]reflect.Value)
			reflectCache[base] = baseCacheMap
		}

		// TODO: more than just the first hop
		if len(baseCacheMap) > 0 {
			// start most specific and move back
			for x := len(keyParts); x > 0; x-- {
				lookupKey := strings.Join(keyParts[:x], ".")
				cachedVal, ok := baseCacheMap[lookupKey]
				if ok {
					currVal = cachedVal
					keyParts = keyParts[x:]
					break
				}
			}

		}
	}

	// caching strategy:
	// if base is a pointer to something, we can cache.
	// What we'll do is store a map of base -> key -> Value (or type, we'll have to see)
	// then we can look for matches

	var subval reflect.Value
	// Since some layers might be pointers, we won't range over the keyParts
	// this way the pointer types can just dereference then continue
	for x := 0; x < len(keyParts); {
		keyPart := keyParts[x]
		incr := true
		switch currVal.Kind() {
		case reflect.Array:
			fallthrough
		case reflect.Slice:
			idx, err := strconv.Atoi(keyPart)
			if err != nil {
				return nil, fmt.Errorf("Accessing an array/slice with a non-int key: %v", keyPart)
			}
			subval = currVal.Index(idx)
		case reflect.Map:
			subval = currVal.MapIndex(reflect.ValueOf(keyPart))
		// if we are a pointer of some kind-- lets dereference the pointer and continue on
		case reflect.Ptr:
			subval = currVal.Elem()
			incr = false
		case reflect.Struct:
			subval = currVal.FieldByName(keyPart)
		default:
			return nil, fmt.Errorf("Unable to Get() past %v %v", currVal, currVal.Kind())
		}
		if subval.IsValid() {
			if baseCacheMap != nil {
				lookupKey := strings.Join(keyParts[:x+1], ".")
				baseCacheMap[lookupKey] = subval
			}
			currVal = subval
			if incr {
				x++
			}
		} else {
			return nil, fmt.Errorf("unable to find %v in %v %v: %v", keyPart, currVal, currVal.Kind(), subval)
		}
	}
	return &currVal, nil
}
