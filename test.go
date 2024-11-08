package main

import (
	"fmt"
	"reflect"
)

func main() {
	var x int = 42
	var y *int = &x

	// TypeOf
	t := reflect.TypeOf(y)
	zzzz := reflect.ValueOf(x)
	fmt.Println("Type:", t) // *int

	// Kind
	k := t.Kind()
	fmt.Println("Kind:", k) // ptr
	fmt.Println("111:", t.Elem())
	fmt.Println("111:", zzzz.Elem()) // int

	// Elem
	e := t.Elem()
	fmt.Println("Elem:", e) // int

	// 通过反射获取指针指向的值
	v := reflect.ValueOf(y).Elem()
	fmt.Println("Value:", v) // 42

	// 修改指针指向的值
	v.SetInt(100)
	fmt.Println("New Value:", x) // 100
}
