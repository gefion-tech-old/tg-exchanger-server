package utils

import (
	"errors"
	"fmt"
	"reflect"
)

func SetReflectIntValue(obj interface{}, field string, valueToSet int) (interface{}, error) {
	// Заполнение объекта querys
	// Может быть любым типом
	val := reflect.ValueOf(obj)

	// Если это указатель
	if val.Kind() == reflect.Ptr {
		val = reflect.Indirect(val)
	}

	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("failed to process the struct %s", reflect.TypeOf(obj).String())
	}

	f := val.FieldByName(field)
	if !f.IsValid() && f.Kind() != reflect.Int {
		return nil, fmt.Errorf("in struct %s, field %s is invalid", reflect.TypeOf(obj).String(), field)
	}

	f.SetInt(int64(valueToSet))
	return obj, nil
}

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
