package main

import (
	"fmt"
	gf8 "galoisfield8"
)

type (
	A []byte
)

func main() {

	t := []byte{ 1 }
	// u := []byte{ 1 }

	fmt.Println( A(t) )

	f := gf8.New()

	fmt.Printf( "%b\n", gf8.Prime )
	gf8.PrintTable( f.PowerTable )

	fmt.Println( "log:" )
	gf8.PrintTable( f.LogTable )
}
