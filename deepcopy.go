package deepcopy

import (
	"reflect"
)

type copier struct {
	ptrs map[reflect.Type]map[uintptr]reflect.Value
}

func (c *copier) copyArray(v reflect.Value) reflect.Value {
	if v.IsNil() {
		return reflect.Zero(v.Type())
	}

	array := reflect.MakeSlice(v.Type(), v.Len(), v.Len())
	for i := 0; i < v.Len(); i++ {
		subValue := c.copy(v.Index(i))
		if !subValue.IsValid() {
			continue
		}

		arrayValue := array.Index(i)
		if !arrayValue.CanSet() {
			continue
		}
		arrayValue.Set(subValue.Convert(arrayValue.Type()))
	}
	return array
}

func (c *copier) copyMap(v reflect.Value) reflect.Value {
	if v.IsNil() {
		return reflect.Zero(v.Type())
	}

	resultMap := reflect.MakeMap(v.Type())
	for _, key := range v.MapKeys() {
		value := v.MapIndex(key)
		duplicatedKey := c.copy(key)
		duplicatedValue := c.copy(value)

		if !isNillable(duplicatedKey) || !duplicatedKey.IsNil() {
			duplicatedKey = duplicatedKey.Convert(key.Type())
		}

		if (!isNillable(duplicatedValue) || !duplicatedValue.IsNil()) && duplicatedValue.IsValid() {
			duplicatedValue = duplicatedValue.Convert(value.Type())
		}

		if !duplicatedValue.IsValid() {
			duplicatedValue = reflect.Zero(resultMap.Type().Elem())
		}

		resultMap.SetMapIndex(duplicatedKey, duplicatedValue)
	}
	return resultMap
}

func isNillable(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Chan, reflect.Interface, reflect.Ptr, reflect.Func:
		return true
	}
	return false
}

func (c *copier) copyPtr(v reflect.Value) reflect.Value {
	if v.IsNil() {
		return reflect.Zero(v.Type())
	}

	var newValue reflect.Value

	if v.Elem().CanAddr() {
		ptrs, exists := c.ptrs[v.Type()]
		if exists {
			if newValue, exists := ptrs[v.Elem().UnsafeAddr()]; exists {
				return newValue
			}
		}
	}

	newValue = c.copy(v.Elem())

	if v.Elem().CanAddr() {
		ptrs, exists := c.ptrs[v.Type()]
		if exists {
			if newValue, exists := ptrs[v.Elem().UnsafeAddr()]; exists {
				return newValue
			}
		}
	}

	ptr := reflect.New(newValue.Type())
	ptr.Elem().Set(newValue)
	return ptr
}

func (c *copier) copyStruct(v reflect.Value) reflect.Value {
	newStructPtr := reflect.New(v.Type())
	newStruct := newStructPtr.Elem()

	if v.CanAddr() {
		ptrs := c.ptrs[newStructPtr.Type()]
		if ptrs == nil {
			ptrs = make(map[uintptr]reflect.Value)
			c.ptrs[newStructPtr.Type()] = ptrs
		}
		ptrs[v.UnsafeAddr()] = newStructPtr
	}

	for i := 0; i < v.NumField(); i++ {
		newStructValue := newStruct.Field(i)
		if !newStructValue.CanSet() {
			continue
		}

		copied := c.copy(v.Field(i))
		if !copied.IsValid() {
			continue
		}

		newStructValue.Set(copied.Convert(newStructValue.Type()))
	}
	return newStruct
}

// Copy takes any objects and create a deep copy of it, and all its fields or elements.
// In case of reflect.Invalid, reflect.Chan and reflect.UnsafePointer types, it returns nil.
func Copy(i interface{}) interface{} {
	c := copier{
		ptrs: map[reflect.Type]map[uintptr]reflect.Value{},
	}
	ret := c.copy(reflect.ValueOf(i))
	if ret.Kind() == reflect.Invalid {
		return nil
	}
	return ret.Interface()
}

func (c *copier) copy(v reflect.Value) reflect.Value {

	switch v.Kind() {
	case reflect.Invalid:
		return reflect.ValueOf(nil)

		// Bool
	case reflect.Bool:
		return reflect.ValueOf(v.Bool())

		// Int
	case reflect.Int:
		return reflect.ValueOf(int(v.Int()))
	case reflect.Int8:
		return reflect.ValueOf(int8(v.Int()))
	case reflect.Int16:
		return reflect.ValueOf(int16(v.Int()))
	case reflect.Int32:
		return reflect.ValueOf(int32(v.Int()))
	case reflect.Int64:
		return reflect.ValueOf(v.Int())

		// Unit
	case reflect.Uint:
		return reflect.ValueOf(uint(v.Uint()))
	case reflect.Uint8:
		return reflect.ValueOf(uint8(v.Uint()))
	case reflect.Uint16:
		return reflect.ValueOf(uint16(v.Uint()))
	case reflect.Uint32:
		return reflect.ValueOf(uint32(v.Uint()))
	case reflect.Uint64:
		return reflect.ValueOf(v.Uint())

		// Float
	case reflect.Float32:
		return reflect.ValueOf(float32(v.Float()))
	case reflect.Float64:
		return reflect.ValueOf(v.Float())

		// Complex
	case reflect.Complex64:
		return reflect.ValueOf(complex64(v.Complex()))
	case reflect.Complex128:
		return reflect.ValueOf(v.Complex())

		// String
	case reflect.String:
		return reflect.ValueOf(v.String())

		// Array
	case reflect.Array, reflect.Slice:
		return c.copyArray(v)

		// Map
	case reflect.Map:
		return c.copyMap(v)

		// Ptr
	case reflect.Ptr:
		return c.copyPtr(v)

		// Struct
	case reflect.Struct:
		return c.copyStruct(v)

		// Func
	case reflect.Func:
		return v

		// Interface
	case reflect.Interface:
		return c.copy(v.Elem())
	}

	return reflect.Zero(v.Type())
}
