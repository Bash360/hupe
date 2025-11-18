package utils

import "reflect"

func ValueToInterface(values []reflect.Value) []any {
	result := make([]any, 0, len(values))

	for i, v := range values {
		result[i] = v.Interface()
	}
	return result
}
