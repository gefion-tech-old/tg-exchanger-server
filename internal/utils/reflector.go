package utils

import (
	"errors"
	"reflect"
)

/*
	Метод для получение функции рефлектора, если она существует
*/
func GetReflectMethod(i interface{}, methodName string) (*reflect.Value, error) {
	var ptr reflect.Value
	var value reflect.Value
	var finalMethod reflect.Value

	value = reflect.ValueOf(i)

	// Если мы начнем с указателя, нам нужно получить значение, указанное на
	// Если мы начнем со значения, нам нужно получить указатель на это значение
	if value.Type().Kind() == reflect.Ptr {
		ptr = value
		value = ptr.Elem()
	} else {
		ptr = reflect.New(reflect.TypeOf(i))
		temp := ptr.Elem()
		temp.Set(value)
	}

	// Проверка метода по значению
	method := value.MethodByName(methodName)
	if method.IsValid() {
		finalMethod = method
	}
	// Проверка метода по указателю
	method = ptr.MethodByName(methodName)
	if method.IsValid() {
		finalMethod = method
	}

	if finalMethod.IsValid() {
		return &finalMethod, nil
	}

	return nil, errors.New("reflect method `" + methodName + "` is not found")
}
