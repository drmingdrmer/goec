package galoisfield8
import (
	"testing"
	// "fmt"
)

type (
	case2 struct {
		a, b, r byte
		mes string
	}
)

func TestNew( t *testing.T ) {
	f := NewField()
	if f.E == 1 {} else {
		t.Errorf( "f.E should be 1 but ", f.E )
	}
	if f.Zero == 0 {} else {
		t.Errorf( "f.Zero should be 0 but ", f.Zero )
	}
}
func TestPowAndLog( t *testing.T ) {

	f := NewField()

	cases := []case2{
		{ 0, 0, 0, "" },
		{ 0, 1, 0, "" },
		{ 0, 2, 0, "" },
		{ 1, 1, 1, "" },
		{ 1, 2, 1, "" },
		{ 2, 1, 2, "" },
		{ 2, 0x0f, 0x26, "" },
		{ 0x8e, 2, 0x47, "" },
		{ 0x8e, 3, 0xad, "" },
	}

	var (
		r byte
	)

	for _, c := range cases {
		r = f.Pow( c.a, c.b )
		if r == c.r {} else {
			t.Errorf( "expect %d + %d = %d, but %d %s", c.a, c.b, c.r, r, c.mes )
		}

		// // TODO test log
		// if c.r != 0 {
		//     r = f.Log( c.r, c.a )
		//     if r == c.b {} else {
		//         t.Errorf( "expect %d - %d = %d, but %d %s", c.r, c.a, c.b, r, c.mes )
		//     }
		// }

	}
}
func TestAddAndSub( t *testing.T ) {

	f := NewField()

	// 2^254 = 0x8e
	cases := []case2{
		{ 0, 0, 0, "" },
		{ 0, 1, 1, "" },
		{ 1, 1, 0, "" },
		{ 1, 2, 3, "" },
		{ 0xde, 0xde, 0, "" },
		{ 0x80, 0xff, 0x7f, "" },
		// { 2, 0x8e, 1, "" },
		// { 0x8e, 0x8e, 0x47, "" },
	}

	PrintTable( f.PowerTable )
	PrintTable( f.LogTable )

	var (
		r byte
	)

	for _, c := range cases {
		r = f.Add( c.a, c.b )
		if r == c.r {} else {
			t.Errorf( "expect %d + %d = %d, but %d %s", c.a, c.b, c.r, r, c.mes )
		}

		r = f.Add( c.b, c.a )
		if r == c.r {} else {
			t.Errorf( "expect %d + %d = %d, but %d %s", c.b, c.a, c.r, r, c.mes )
		}

		r = f.Sub( c.r, c.a )
		if r == c.b {} else {
			t.Errorf( "expect %d - %d = %d, but %d %s", c.r, c.a, c.b, r, c.mes )
		}

		r := f.Sub( c.r, c.b )
		if r == c.a {} else {
			t.Errorf( "expect %d - %d = %d, but %d %s", c.r, c.b, c.a, r, c.mes )
		}
	}
}
func TestMulAndDiv( t *testing.T ) {

	f := NewField()

	// 2^254 = 0x8e
	cases := []case2{
		{ 0, 0, 0, "" },
		{ 0, 1, 0, "" },
		{ 1, 1, 1, "" },
		{ 1, 2, 2, "" },
		{ 1, 0x8e, 0x8e, "" },
		{ 2, 0x8e, 1, "" },
		{ 0x8e, 0x8e, 0x47, "" },
	}

	var (
		r byte
	)

	PrintTable( f.PowerTable )
	PrintTable( f.LogTable )

	for _, c := range cases {
		r = f.MulBy( c.a, c.b )
		if r == c.r {} else {
			t.Errorf( "expect %d * %d = %d, but %d %s", c.a, c.b, c.r, r, c.mes )
		}

		r = f.MulBy( c.b, c.a )
		if r == c.r {} else {
			t.Errorf( "expect %d * %d = %d, but %d %s", c.b, c.a, c.r, r, c.mes )
		}

		if c.a != 0 {
			r := f.Div( c.r, c.a )
			if r == c.b {} else {
				t.Errorf( "expect %d / %d = %d, but %d %s", c.r, c.a, c.b, r, c.mes )
			}
		}

		if c.b != 0 {
			r := f.Div( c.r, c.b )
			if r == c.a {} else {
				t.Errorf( "expect %d / %d = %d, but %d %s", c.r, c.b, c.a, r, c.mes )
			}
		}
	}
}


func BenchmarkMul( b *testing.B ) {
	f := NewField()

	// a := 3
	// c := 4

	var (
		c byte = 0
	)
	for i := 0; i < b.N; i++ {
		// a = a ^ c
		// c = f.MulTable[ ((i&0xff)<<8) & 34  ]
		// c = f.MulBy( byte(i&0xff), 34 )
		c = f.MulByAct( 24, 34 )
	}
	_ = c
}
