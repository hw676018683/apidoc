// Set default values for struct
package defaults

import (
	"log"
	"reflect"
)

func Set(ptr interface{}) interface{} {
	if ptr == nil {
		return nil
	}

	v := reflect.ValueOf(ptr)
	if reflect.TypeOf(ptr).Kind() == reflect.Ptr {
		v = reflect.ValueOf(ptr).Elem()
	}
	types := make([]reflect.Type, 0)
	setField(types, v)
	return v.Interface()
}

func setField(types []reflect.Type, field reflect.Value) {
	if !field.CanSet() {
		// fmt.Println(`can not set :`, field.Type().Name())
		return
	}
	switch field.Kind() {
	case reflect.Array:
		for j := 0; j < field.Len(); j++ {
			setField(types, field.Index(j))
		}
	case reflect.Map:
		mType := field.Type()
		field.Set(reflect.MakeMap(mType))

		key := reflect.New(mType.Key()).Elem()
		val := reflect.New(mType.Elem()).Elem()
		setField(types, val)
		field.SetMapIndex(key, val)
	case reflect.Slice:
		field.Set(reflect.MakeSlice(field.Type(), 1, 1))
		setField(types, field.Index(0))
	case reflect.Ptr:
		field.Set(reflect.New(field.Type().Elem()))
		setField(types, field.Elem())
	case reflect.Interface:
		if !field.Elem().IsValid() {
			return // 没有给接口赋值，无法初始化
		}
		newField := reflect.New(field.Elem().Type()).Elem()
		setField(types, newField)
		field.Set(newField)
	case reflect.Struct:
		t := field.Type()
		for _, v := range types {
			if v == t {
				log.Println(`-- Loop struct: ` + t.String())
				log.Println(types)
				return
			}
		}
		curTypes := make([]reflect.Type, len(types))
		copy(curTypes, types)
		curTypes = append(curTypes, t)
		for i := 0; i < field.NumField(); i++ {
			setField(curTypes, field.Field(i))
		}
	}
}
