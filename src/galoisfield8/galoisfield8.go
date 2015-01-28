package galoisfield8

import (
	"fmt"
)

type (
	// ElementCaster interface {}

	OpTable []byte

	FieldOperator interface {

		Pow( byte, byte ) byte
		Log( byte ) byte

		Add( byte, byte ) byte
		Sub( byte, byte ) byte
		MulBy( byte, byte ) byte
		Div( byte, byte ) byte

		GetE() byte
		GetZero() byte
	}

	F struct {
		E, Zero byte
		PowerTable OpTable
		LogTable OpTable

		MulTable OpTable
	}
)

const (
	// galios field of size 256
	// x^8 + x^4 + x^3 + x^2 + 1
	// binary: 1-0001-1101
	Prime int = 0x11d
	Order int = 256
)

func NewField() *F {

	f := F{
		E: 1,
		Zero: 0,
		PowerTable: make( OpTable, Order ),
		LogTable: make( OpTable, Order ),
		MulTable: make( OpTable, Order * Order ),
	}

	var (
		n int = 1
	)

	for i := 0; i < Order; i++ {

		f.PowerTable[ i ] = byte(n)
		f.LogTable[ n ] = byte(i)

		n *= 2
		if n >= Order  {
			n = n ^ Prime
		}
	}

	for i := 0; i < Order; i++ {
		for j := 0; j < Order; j++ {
			f.MulTable[ i*Order + j ] = f.MulByAct( byte(i), byte(j) )
		}
	}

	return &f
}

func PrintTable( t OpTable ) {
	for i := 0; i < len(t); i++ {
		fmt.Printf( " %02x", t[ i ] )
		if i % 16 == 15 {
			fmt.Println( "" )
		}
	}
}

func ( f *F )Pow( n byte, pow byte ) byte {
	if n == 0 {
		return 0
	}

	l := int(f.LogTable[ n ]) * int(pow)
	l %= Order-1
	return f.PowerTable[ l ]
}

func ( f *F )Log( n byte ) byte {
	if n == 0 {
		panic( "log of 0" )
	}
	return f.LogTable[ n ]
}

func ( f *F )MulBy( a byte, b byte ) byte {
	return f.MulByAct( a, b )
	// return f.MulTable[ (int(a)<<8) + int(b) ]
}

func ( f *F )MulByAct( a byte, b byte ) byte {

	// if a == 0 || b == 0 {
	if a*b == 0 {
		return 0
	}

	l := int(f.LogTable[ a ]) + int(f.LogTable[ b ])
	if l >= Order-1 {
		l -= Order-1
	}

	return f.PowerTable[ l ]
}

func ( f *F ) Div( a byte, b byte ) byte {
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

func ( f *F ) Add( a byte, b byte ) byte {
	return a ^ b
}

func ( f *F ) Sub( a byte, b byte ) byte {
	return f.Add( a, b )
}
