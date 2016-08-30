package udn

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// TODO: finish these
// The set stuff is not functional yet
func Set(base interface{}, key string, val interface{}) error {
	keyParts := strings.Split(key, ".")
	return SetParts(base, keyParts, val)
}

func SetParts(base interface{}, keyParts []string, v interface{}) error {
	keyPrefix := keyParts[:len(keyParts)-1]
	lastKey := keyParts[len(keyParts)-1]
	obj, err := GetParts(base, keyPrefix)
	if err != nil {
		return err
	}
	objR := reflect.ValueOf(obj)
	val := reflect.ValueOf(v)
	switch objR.Kind() {
	case reflect.Array:
		return fmt.Errorf("Arrays aren't supported")
		/*
			// Seem to only get errors like:
			//		reflect.Value.Set using unaddressable value [recovered]
			//	something to do with the array not having the values directly-- basically
			// the index isn't `CanSet()`
			idx, err := strconv.Atoi(lastKey)
			if err != nil {
				return fmt.Errorf("Accessing an array with a non-int key: %v", lastKey)
			}
			if idx >= objR.Len() {
				return fmt.Errorf("array isn't large enough")
			} else {
				logrus.Infof("val=%v idx=%v", val, objR.Index(idx))
				logrus.Infof("settability of objR=%v index=%v", objR.CanSet(), objR.Index(idx).CanSet())
				objR.Index(idx).Set(val)
			}
		*/
	case reflect.Slice:
		idx, err := strconv.Atoi(lastKey)
		if err != nil {
			return fmt.Errorf("Accessing an slice with a non-int key: %v", lastKey)
		}
		if idx >= objR.Len() {
			newSlice := reflect.Append(objR, val)
			SetParts(base, keyPrefix, newSlice.Interface())
		} else {
			objR.Index(idx).Set(val)
		}
	case reflect.Map:
		mapKey := reflect.ValueOf(lastKey)
		objR.SetMapIndex(mapKey, val)
	case reflect.Ptr:
		return fmt.Errorf("don't know what to do with pointers yet")
	case reflect.Struct:
		return fmt.Errorf("don't know what to do with struct yet")
	default:
		return fmt.Errorf("Unable to Get() past %v %v", objR, objR.Kind())
	}

	return nil
}
