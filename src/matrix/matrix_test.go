package matrix
import (
	"fmt"
	"testing"
	"math/rand"
)

type (
	tmatrix [][]byte
	tcase struct {
		a, b, ab, ba tmatrix
		mes string
		err error
	}
)

func TestNewUnit( t *testing.T ) {

	InitField()

	mtx, err := NewUnit( 3 )
	if err == nil {} else {
		t.Errorf( "err should be nil", err )
	}

	for i := 0; i < 3; i++ {
		if mtx.elts[i][i] == 1 {} else {
			t.Errorf( "es[i][i] should be 1 but ", mtx.elts[i][i] )
		}
	}

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if i == j {
				continue
			}
			if mtx.elts[i][j] == 0 {} else {
				t.Errorf( "es[i][i] should be 0 but ", mtx.elts[i][i] )
			}
		}
	}
}

func TestCopy( t *testing.T ) {
	InitField()

	mtx, err := NewUnit( 3 )
	if err == nil {} else {
		t.Errorf( "err should be nil but ", err )
	}
	m2 := mtx.Copy()

	if mtx.n == m2.n && mtx.m == m2.m {} else {
		t.Errorf( "m and n are copied" )
	}
	if mtx.f == m2.f {} else {
		t.Errorf( "field is copied" )
	}

	mtx.elts[1][1] = 3
	if m2.elts[1][1] == 1 {} else {
		t.Errorf( "not changed, but", m2.elts )
	}
}

func TestInverse( t *testing.T ) {
	InitField()

	unit, err := NewUnit( 3 )
	if err == nil {} else {
		t.Errorf( "err should be nil", err )
	}

	cases := []tcase{
		{
			a: tmatrix{
				{1, 0, 0},
				{0, 1, 0},
				{0, 0, 0},
			},
			err: ErrDetZero,
		},
		{
			a: tmatrix{
				{1, 0, 0},
				{0, 1, 0},
				{0, 0, 1},
			},
			b: tmatrix{
				{1, 0, 0},
				{0, 1, 0},
				{0, 0, 1},
			},
		},
		{
			a: tmatrix{
				{1, 0, 0},
				{0, 1, 0},
				{0, 0, 1},
			},
			b: tmatrix{
				{1, 0, 0},
				{0, 1, 0},
				{0, 0, 1},
			},
		},
		{
			a: tmatrix{
				{2, 0, 0},
				{0, 1, 0},
				{0, 0, 1},
			},
			b: tmatrix{
				{0x8e, 0, 0},
				{0, 1, 0},
				{0, 0, 1},
			},
		},
		{
			a: tmatrix{
				{1, 2, 0},
				{0, 1, 0},
				{0, 0, 1},
			},
			b: tmatrix{
				{1, 2, 0},
				{0, 1, 0},
				{0, 0, 1},
			},
		},
		{
			a: tmatrix{
				{1, 0, 2},
				{0, 1, 0},
				{0, 0, 1},
			},
			b: tmatrix{
				{1, 0, 2},
				{0, 1, 0},
				{0, 0, 1},
			},
		},
		{
			a: tmatrix{
				{2, 3, 4},
				{0, 1, 0},
				{0, 0, 1},
			},
			b: tmatrix{
				{0x8e, 0x8f, 0x02},
				{0, 1, 0},
				{0, 0, 1},
			},
		},
	}

	for _, c := range cases {
		ma, err := NewFromSlice( c.a )
		if err == nil {} else {
			t.Errorf( "new matrix expect no error but: ", err )
		}
		mb, err := NewFromSlice( c.b )
		if err == nil {} else {
			t.Errorf( "new matrix expect no error but: ", err )
		}

		imtx, err := ma.Inverse()
		if err == c.err {} else {
			t.Errorf( "expect err=%s but %s", c.err, err )
		}

		if c.err != nil {
			continue
		}

		fmt.Println( "inversed:" )
		imtx.Print()
		if imtx.Equal( mb ) {} else {
			t.Errorf( "expected to be: ", c.b, imtx )
		}

		err = ma.MulInplace( imtx )
		if err == nil {} else {
			t.Errorf( "err should be nil, but", err )
		}

		fmt.Println( "ma * imtx:" )
		ma.Print()

		eq := ma.Equal( unit )
		if eq == true {} else {
			t.Errorf( "ma * imtx should be unit" )
		}
	}
}

func TestInverseRandom( t *testing.T ) {
	InitField()
	fmt.Println( "---- inverse random" )

	mtx, err := NewSquare( 3 )
	if err == nil {} else {
		t.Errorf( "err should be nil but ", err )
	}

	unit, err := NewUnit( 3 )

	for q := 0; q < 100; q++ {
		for i := 0; i < mtx.n; i++ {
			for j := 0; j < mtx.m; j++ {
				mtx.elts[i][j] = byte(rand.Int())
			}
		}
		imtx, err := mtx.Inverse()
		if err != nil {
			continue
		}

		fmt.Println( "mtx:" )
		mtx.Print()

		err = mtx.MulInplace( imtx )
		if err == nil {} else {
			t.Errorf( "err should be nil when mtx.MulInplace, but ", err )
		}

		fmt.Println( "inverse: err=", err )
		imtx.Print()
		fmt.Println( "mtx * inverse:" )
		mtx.Print()

		r := mtx.Equal( unit )
		if r == true {} else {
			t.Errorf( "mtx * inversed = unit!" )
		}

	}
}

func TestEqual( t *testing.T ) {
	cases := []tcase{
		{ a:tmatrix{{ 1 }}, b:tmatrix{{ 1 }} },
		{
			a: tmatrix{
				{1, 2, 3},
				{2, 3, 4},
				{3, 4, 5},
			},
			b: tmatrix{
				{1, 2, 3},
				{2, 3, 4},
				{3, 4, 5},
			},
		},
	}

	var (
		m *Matrix8 = nil
		m1 *Matrix8 = nil
	)

	InitField()

	m1, err := NewMatrix( 1, 1 )
	if err == nil {} else {
		t.Errorf( "err should be nil ", err )
	}

	if m.Equal( m ) == false {} else {
		t.Errorf( "expected to be false but true, nil eq nil" )
	}

	if m.Equal( m1 ) == false {} else {
		t.Errorf( "m != m1" )
	}

	if m1.Equal( m ) == false {} else {
		t.Errorf( "m1 != m" )
	}

	if m1.Equal( m1 ) == true {} else {
		t.Errorf( "m1 != m1" )
	}

	for _, c := range cases {
		ma, err := NewFromSlice( c.a )
		if err == nil {} else {
			t.Errorf( "new matrix expect no error but: ", err )
		}
		mb, err := NewFromSlice( c.b )
		if err == nil {} else {
			t.Errorf( "new matrix expect no error but: ", err )
		}

		r := ma.Equal( mb ) && mb.Equal( ma )
		if r == true {} else {
			t.Errorf( "expect a == b ", err )
		}

		ma.elts[0][0] += 1
		r = !ma.Equal( mb ) && !mb.Equal( ma )
		if r == true {} else {
			t.Errorf( "expect a != b ", err )
		}
	}

}

func TestMul( t *testing.T ) {
	cases := []tcase{
		{
			a: tmatrix{
				{1, 0, 0},
				{0, 2, 0},
				{0, 4, 3},
			},
			b: tmatrix{
				{1, 0, 0},
				{1, 1, 0},
				{0, 0, 1},
			},
			ab: tmatrix{
				{1, 0, 0},
				{2, 2, 0},
				{4, 4, 3},
			},
		},
		{
			a: tmatrix{
				{1, 0, 0},
				{0, 2, 0},
				{0, 4, 3},
			},
			b: tmatrix{
				{1, 1, 0},
				{0, 1, 0},
				{0, 0, 1},
			},
			ab: tmatrix{
				{1, 1, 0},
				{0, 2, 0},
				{0, 4, 3},
			},
		},
		{
			a: tmatrix{
				{1, 0, 0},
				{0, 2, 0},
				{0, 4, 3},
			},
			b: tmatrix{
				{1, 0, 0},
				{0, 0, 1},
				{0, 1, 0},
			},
			ab: tmatrix{
				{1, 0, 0},
				{0, 0, 2},
				{0, 3, 4},
			},
		},
		{
			a: tmatrix{
				{1, 0, 0},
				{0, 2, 0},
				{0, 0, 3},
			},
			b: tmatrix{
				{3, 0, 0},
				{0, 9, 0},
				{0, 0, 0x0c},
			},
			ab: tmatrix{
				{3, 0, 0},
				{0, 0x12, 0},
				{0, 0, 0x14},
			},
		},
		{
			a: tmatrix{
				{1, 2, 3},
				{4, 5, 6},
				{7, 8, 9},
			},
			b: tmatrix{
				{1, 3, 5},
				{7, 9, 0x0a},
				{0x0b, 0x0c, 0x0d},
			},
			ab: tmatrix{
				{0x12, 5, 6},
				{0x25, 0x09, 0x18},
				{0x6c, 0x2d, 0x2e},
			},
		},
		{
			a: tmatrix{
				{1, 2},
				{0, 1},
				{1, 2},
			},
			b: tmatrix{
				{1, 3, 5},
				{7, 9, 1},
			},
			ab: tmatrix{
				{0x0f, 0x11, 0x07},
				{0x07, 0x09, 0x01},
				{0x0f, 0x11, 0x07},
			},
			mes: "3*2 * 2*3",
		},
	}

	InitField()

	for _, c := range cases {
		ma, err := NewFromSlice( c.a )
		if err == nil {} else {
			t.Errorf( "new matrix expect no error but: ", err )
		}
		mb, err := NewFromSlice( c.b )
		if err == nil {} else {
			t.Errorf( "new matrix expect no error but: ", err )
		}

		r, err := ma.MulBy( mb )
		if err == nil {} else {
			t.Errorf( "expect no err, but ", err )
		}

		for i := 0; i < r.n; i++ {
			for j := 0; j < r.m; j++ {
				if r.elts[i][j] == c.ab[i][j] {} else {
					t.Errorf( "expect r[%d][%d] to be %02x but %02x %s", i, j, c.ab[i][j], r.elts[i][j], c.mes )
				}
			}
		}
	}
}

func TestMulInplace( t *testing.T ) {
	cases := []tcase{
		{
			a: tmatrix{
				{1, 0, 0},
				{0, 2, 0},
				{0, 4, 3},
			},
			b: tmatrix{
				{1, 0, 0},
				{1, 1, 0},
				{0, 0, 1},
			},
			ab: tmatrix{
				{1, 0, 0},
				{2, 2, 0},
				{4, 4, 3},
			},
			ba: tmatrix{
				{1, 0, 0},
				{1, 2, 0},
				{0, 4, 3},
			},
		},
		{
			a: tmatrix{
				{1, 0, 0},
				{0, 2, 0},
				{0, 4, 3},
			},
			b: tmatrix{
				{1, 1, 0},
				{0, 1, 0},
				{0, 0, 1},
			},
			ab: tmatrix{
				{1, 1, 0},
				{0, 2, 0},
				{0, 4, 3},
			},
			ba: tmatrix{
				{1, 2, 0},
				{0, 2, 0},
				{0, 4, 3},
			},
		},
		{
			a: tmatrix{
				{1, 0, 0},
				{0, 2, 0},
				{0, 4, 3},
			},
			b: tmatrix{
				{1, 0, 0},
				{0, 0, 1},
				{0, 1, 0},
			},
			ab: tmatrix{
				{1, 0, 0},
				{0, 0, 2},
				{0, 3, 4},
			},
			ba: tmatrix{
				{1, 0, 0},
				{0, 4, 3},
				{0, 2, 0},
			},
		},
	}

	InitField()

	for _, c := range cases {
		ma, err := NewFromSlice( c.a )
		if err == nil {} else {
			t.Errorf( "new matrix expect no error but: ", err )
		}
		mb, err := NewFromSlice( c.b )
		if err == nil {} else {
			t.Errorf( "new matrix expect no error but: ", err )
		}

		err = ma.MulInplace( mb )
		if err == nil {} else {
			t.Errorf( "expect no err, but ", err )
		}

		for i := 0; i < ma.n; i++ {
			for j := 0; j < ma.m; j++ {
				if ma.elts[i][j] == c.ab[i][j] {} else {
					t.Errorf( "expect ma[%d][%d] to be %d but %d", i, j, c.ab[i][j], ma.elts[i][j] )
				}
			}
		}
	}
}

func BenchmarkMulVec( b *testing.B ) {
	InitField()

	// a := 3
	// c := 4

	ab := tmatrix{
		{0x0f, 0x11, 0x07},
		{0x07, 0x09, 0x01},
		{0x0f, 0x11, 0x07},
	}
	ma, _ := NewFromSlice( ab )

	v := []byte{ 3, 4, 5 }
	rst := []byte{ 0, 0, 0 }

	for i := 0; i < b.N; i++ {
		// ma.MulVecToBytes( v, rst )
		ma.MulByVecToBytes( v, rst )
	}
}
