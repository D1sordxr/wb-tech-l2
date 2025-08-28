package main

import (
	"fmt"
	"os"
)

func Foo() error {
	var err *os.PathError = nil // err - указатель на os.PathError, равен nil
	return err                  // возвращается как интерфейс error
}

func main() {
	err := Foo()            // переменная реализующая интерфейс error
	fmt.Println(err)        // nil - значение
	fmt.Println(err == nil) // false, так как хранит тип os.PathError

	// интерфейсы: имеют определённый набор методов
	// хранят в себе тип и значение

	// пустые интерфейсы:
	// не требуют реализации, но также хранят в себе тип и значение
	// могут содержать любой тип interface{} == any
}
