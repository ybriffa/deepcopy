package deepcopy

import (
	"reflect"
)

func copyArray(v reflect.Value) interface{} {
	if v.IsNil() {
		return nil
	}

	array := reflect.MakeSlice(v.Type(), v.Len(), v.Len())
	for i := 0; i < v.Len(); i++ {
		rawSubValue := Copy(v.Index(i).Interface())
		if rawSubValue == nil {
			continue
		}
		subValue := reflect.ValueOf(rawSubValue)

		arrayValue := array.Index(i)
		if !arrayValue.CanSet() {
			continue
		}
		arrayValue.Set(subValue.Convert(arrayValue.Type()))
	}
	return array.Interface()
}

func copyMap(v reflect.Value) interface{} {
	if v.IsNil() {
		return nil
	}

	resultMap := reflect.MakeMap(v.Type())
	for _, key := range v.MapKeys() {
		value := v.MapIndex(key)
		duplicatedKey := reflect.ValueOf(Copy(key.Interface()))
		duplicatedValue := reflect.ValueOf(Copy(value.Interface()))

		resultMap.SetMapIndex(duplicatedKey.Convert(key.Type()), duplicatedValue.Convert(value.Type()))
	}
	return resultMap.Interface()
}

func copyPtr(v reflect.Value) interface{} {
	if v.IsNil() {
		return nil
	}
	retValue := Copy(v.Elem().Interface())
	ptr := reflect.New(reflect.ValueOf(retValue).Type())
	ptr.Elem().Set(reflect.ValueOf(retValue))
	return ptr.Interface()
}

func copyStruct(v reflect.Value) interface{} {
	newStruct := reflect.New(v.Type()).Elem()

	for i := 0; i < v.NumField(); i++ {
		newStructValue := newStruct.Field(i)
		if !newStructValue.CanSet() {
			continue
		}

		rawCopied := Copy(v.Field(i).Interface())
		if rawCopied == nil {
			continue
		}

		copied := reflect.ValueOf(rawCopied)
		if copied.Kind() == reflect.Ptr && copied.IsNil() {
			continue
		}

		newStructValue.Set(copied.Convert(newStructValue.Type()))
	}
	return newStruct.Interface()
}

// Copy takes any objects and create a deep copy of it, and all its fields or elements.
// In case of reflect.Invalid, reflect.Chan and reflect.UnsafePointer types, it returns nil.
func Copy(i interface{}) interface{} {
	v := reflect.ValueOf(i)

	switch v.Kind() {
	case reflect.Invalid:
		return nil

		// Bool
	case reflect.Bool:
		return v.Bool()

		// Int
	case reflect.Int:
		return int(v.Int())
	case reflect.Int8:
		return int8(v.Int())
	case reflect.Int16:
		return int16(v.Int())
	case reflect.Int32:
		return int32(v.Int())
	case reflect.Int64:
		return v.Int()

		// Unit
	case reflect.Uint:
		return uint(v.Uint())
	case reflect.Uint8:
		return uint8(v.Uint())
	case reflect.Uint16:
		return uint16(v.Uint())
	case reflect.Uint32:
		return uint32(v.Uint())
	case reflect.Uint64:
		return v.Uint()

		// Float
	case reflect.Float32:
		return float32(v.Float())
	case reflect.Float64:
		return v.Float()

		// Complex
	case reflect.Complex64:
		return complex64(v.Complex())
	case reflect.Complex128:
		return v.Complex()

		// String
	case reflect.String:
		return v.String()

		// Array
	case reflect.Array, reflect.Slice:
		return copyArray(v)

		// Map
	case reflect.Map:
		return copyMap(v)

		// Ptr
	case reflect.Ptr:
		return copyPtr(v)

		// Struct
	case reflect.Struct:
		return copyStruct(v)

		// Func
	case reflect.Func:
		return i
	}

	return nil
}
