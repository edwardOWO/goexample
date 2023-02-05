package learn

import (
	"fmt"
	"reflect"
)

/*
func GetReflect(v interface{}) {

	t := reflect.ValueOf(v)

	//fmt.Println("Data type: \n" + t.Type().String())

	for index := 0; index < t.NumField(); index++ {

		//fmt.Println(strconv.Itoa(index) + " :" + t.Type().Field(index).Name)

		fmt.Println(t.Field(index).Interface())

	}

}
*/
func GetReflect(v interface{}) {

	e := reflect.ValueOf(v)

	fmt.Printf("%v \n", e.Type())
	for i := 0; i < e.NumField(); i++ {
		varName := e.Type().Field(i).Name
		varType := e.Type().Field(i).Type
		varValue := e.Field(i)
		fmt.Printf("%v %v %v \n", varName, varType, varValue)
	}
}

func PrintValue(v interface{}) {

	fmt.Print(v.(int))

	switch v.(type) {

	case int:
		fmt.Println("b1.(type):", "int", v)

	default:
		break

	}
}
