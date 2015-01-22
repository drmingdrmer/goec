package galoisfield8

import (
	"fmt"
)

type (
	// ElementCaster interface {}

	Elt uint8
	BigElt uint16
	Table []Elt

	FieldOperator interface {

		Pow( Elt, Elt ) Elt
		Log( Elt ) Elt

		Add( Elt, Elt ) Elt
		Sub( Elt, Elt ) Elt
		MulBy( Elt, Elt ) Elt
		Div( Elt, Elt ) Elt

		GetE() Elt
		GetZero() Elt
	}

	F struct {
		E, Zero Elt
		PowerTable Table
		LogTable Table
	}
)

const (
	// galios field of size 256
	// x^8 + x^4 + x^3 + x^2 + 1
	// binary: 1-0001-1101
	Prime int = 0x11d
	Order int = 256
)

func New() *F {

	f := F{
		E: 1,
		Zero: 0,
		PowerTable: make( Table, Order ),
		LogTable: make( Table, Order ),
	}

	var (
		n int = 1
	)

	for i := 0; i < Order; i++ {

		f.PowerTable[ i ] = Elt(n)
		f.LogTable[ n ] = Elt(i)

		n *= 2
		if n >= Order  {
			n = n ^ Prime
		}
	}
	return &f
}

func PrintTable( t Table ) {
	for i := 0; i < Order; i++ {
		fmt.Printf( " %02x", t[ i ] )
		if i % 16 == 15 {
			fmt.Println( "" )
		}
	}
}

func ( f *F )Pow( n Elt, pow Elt ) Elt {
	if n == 0 {
		return 0
	}

	l := int(f.LogTable[ n ]) * int(pow)
	l %= Order-1
	return f.PowerTable[ l ]
}

func ( f *F )Log( n Elt ) Elt {
	if n == 0 {
		panic( "log of 0" )
	}
	return f.LogTable[ n ]
}

func ( f *F )MulBy( a Elt, b Elt ) Elt {
	if a == 0 || b == 0 {
		return 0
	}

	l := int(f.LogTable[ a ]) + int(f.LogTable[ b ])
	if l >= Order-1 {
		l -= Order-1
	}

	return f.PowerTable[ l ]
}

func ( f *F ) Div( a Elt, b Elt ) Elt {
	if b == 0 {
		panic( "divide by 0" )
	}
	if a == 0 {
		return 0
	}

	l := int(f.LogTable[ a ]) - int(f.LogTable[ b ])
	if l < 0 {
		l += Order-1
	}

	return f.PowerTable[ l ]
}

func ( f *F ) Add( a Elt, b Elt ) Elt {
	return a ^ b
}

func ( f *F ) Sub( a Elt, b Elt ) Elt {
	return f.Add( a, b )
}
