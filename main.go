package main

import (
	"fmt"
)

func main() {
	srv := NewServer("")
	err := srv.Run()
	if err != nil {
		panic(err)

	}

	fmt.Println("here we go")
}
