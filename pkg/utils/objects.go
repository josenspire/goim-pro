package utils

import "reflect"

func DeepEqual(if1 interface{}, if2 interface{}) bool {
	return reflect.DeepEqual(if1, if2)
}

func TransformStructToMap(obj interface{}) map[string]interface{}{
	obj1 := reflect.TypeOf(obj)
	obj2 := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < obj1.NumField(); i++ {
		data[obj1.Field(i).Name] = obj2.Field(i).Interface()
	}
	return data
}

