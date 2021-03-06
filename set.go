package udn

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func Set(base interface{}, key string, val interface{}) error {
	keyParts := strings.Split(key, ".")
	return SetParts(base, keyParts, val)
}

func SetParts(base interface{}, keyParts []string, v interface{}) error {
	keyPrefix := keyParts[:len(keyParts)-1]
	lastKey := keyParts[len(keyParts)-1]
	obj, done, err := getSetParts(base, keyPrefix, v, true)
	if err != nil {
		return err
	}
	// if the setter stuff already did it-- we are all clear
	if done {
		return err
	}
	objR := *obj
	// This for loop is only for pointers-- so we don't have to duplicate the code all over
	for {
		val := reflect.ValueOf(v)
		switch objR.Kind() {
		case reflect.Array:
			idx, err := strconv.Atoi(lastKey)
			if err != nil {
				return fmt.Errorf("Accessing an array with a non-int key: %v", lastKey)
			}
			if idx >= objR.Len() {
				return fmt.Errorf("array isn't large enough")
			} else {
				entry := objR.Index(idx)
				if entry.CanSet() {
					objR.Index(idx).Set(val)
				} else {
					// TODO: this
					return fmt.Errorf("Unable to set %v to %v in array %v, usually means the `base` was passed-by-value", lastKey, val, entry)
				}
			}
		case reflect.Slice:
			idx, err := strconv.Atoi(lastKey)
			if err != nil {
				return fmt.Errorf("Accessing an slice with a non-int key: %v", lastKey)
			}
			if idx >= objR.Len() {
				newSlice := reflect.Append(objR, val)
				SetParts(base, keyPrefix, newSlice.Interface())
			} else {
				entry := objR.Index(idx)
				if entry.CanSet() {
					objR.Index(idx).Set(val)
				} else {
					// TODO: this
					return fmt.Errorf("Unable to set %v to %v in slice %v, usually means the `base` was passed-by-value", lastKey, val, entry)
				}
			}
		case reflect.Map:
			mapKey := reflect.ValueOf(lastKey)
			objR.SetMapIndex(mapKey, val)
		case reflect.Ptr:
			objR = objR.Elem()
			continue
		case reflect.Struct:
			field := objR.FieldByName(lastKey)
			if field.CanSet() {
				field.Set(val)
			} else {
				// TODO: this
				return fmt.Errorf("Unable to set %v to %v in struct %v, usually means the `base` was passed-by-value", lastKey, val, objR)
			}
		default:
			return fmt.Errorf("Unable to Get() past %v %v", objR, objR.Kind())
		}
		break
	}

	return nil
}
