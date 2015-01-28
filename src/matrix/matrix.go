package matrix

import (
	"errors"
	"fmt"
	gf8 "galoisfield8"
)

type (
	Matrix8 struct {
		n, m int

		elts [][]byte
		f *gf8.F
	}
)

var (
	f *gf8.F
)

var (
	ErrUninitializedField = errors.New( "supporting field is not initialized" )
	ErrDifferenSize = errors.New( "matrices with different n*m" )
	ErrDetZero = errors.New( "determinant is 0" )
)

func InitField() error {
	f = gf8.NewField()
	return nil
}

func NewMatrix( n int, m int ) ( mtx *Matrix8, err error ) {
	if f == nil {
		err = ErrUninitializedField
		return
	}
	mtx = &Matrix8{
		n:n,
		m:m,
		elts: make( [][]byte, n ),
		f: f,
	}
	for i := 0; i < n; i++ {
		mtx.elts[i] = make( []byte, m )
	}
	return
}
func NewUnit( n int ) ( mtx *Matrix8, err error ) {
	mtx, err = NewMatrix( n, n )
	if err != nil {
		return
	}

	for i := 0; i < n; i++ {
		mtx.elts[i][i] = 1
	}
	return
}
func NewSquare( n int ) ( mtx *Matrix8, err error ) {
	mtx, err = NewMatrix( n, n )
	return
}

func ( mtx *Matrix8 ) Copy() *Matrix8{
	if mtx == nil {
		return nil
	}

	m := Matrix8{
		n: mtx.n,
		m: mtx.m,
		f: mtx.f,
		elts: make( [][]byte, mtx.n ),
	}
	for i := 0; i < mtx.n; i++ {
		m.elts[i] = make( []byte, mtx.m )
		copy( m.elts[i], mtx.elts[i] )
	}
	return &m
}
func ( mtx *Matrix8 ) T() *Matrix8 {

	if mtx == nil {
		return nil
	}

	rst, err := NewMatrix( mtx.m, mtx.n )
	// should not happen
	if err != nil {
		return nil
	}
	for i := 0; i < mtx.n; i++ {
		for j := 0; j < mtx.m; j++ {
			rst.elts[j][i] = mtx.elts[i][j]
		}
	}

	return rst
}
func ( mtx *Matrix8 ) SubMatrix( rows []int, cols []int ) *Matrix8{
	m2, err := NewMatrix( len( rows ), len( cols ) )
	if err != nil {
		// should not happen
		return nil
	}
	for i := 0; i < len( rows ); i++ {
		for j := 0; j < len( cols ); j++ {
			m2.elts[ i ][ j ] = mtx.elts[ rows[ i ] ][ cols[ j ] ]
		}
	}

	return m2
}

func ( mtx *Matrix8 ) MulVec( v []byte ) []byte {
	rst := make( []byte, mtx.m )
	mtx.MulVecTo( v, rst )
	return rst
}

func ( mtx *Matrix8 ) MulVecTo( v []byte, rst []byte ) {
	var (
		sum byte = 0
	)

	f := mtx.f
	for j := 0; j < mtx.m; j++ {
		sum = 0
		for i := 0; i < mtx.n; i++ {
			sum = f.Add( sum, f.MulBy( v[i], mtx.elts[i][j] ) )
		}
		rst[j] = sum
	}
}
func ( mtx *Matrix8 ) MulVecToBytes( v []byte, rst []byte ) {
	var (
		sum byte = 0
	)

	f := mtx.f
	for j := 0; j < mtx.m; j++ {
		sum = 0
		for i := 0; i < mtx.n; i++ {
			sum = f.Add( sum, f.MulBy( byte(v[i]), byte(mtx.elts[i][j]) ) )
		}
		rst[j] = byte(sum)
	}
}
func ( mtx *Matrix8 ) MulByVecToBytes( v []byte, rst []byte ) {
	var (
		sum byte = 0
	)

	f := mtx.f
	for i := 0; i < mtx.n; i++ {
		sum = 0
		for j := 0; j < mtx.m; j++ {
			sum = f.Add( sum, f.MulBy( byte(mtx.elts[i][j]), byte(v[j]) ) )
		}
		rst[i] = byte(sum)
	}
}

func ( mtx *Matrix8 ) MulBy( m2 *Matrix8 ) ( rst *Matrix8, err error ) {
	if mtx.m != m2.n {
		err = ErrDifferenSize
		return
	}

	rst, err = NewMatrix( mtx.n, m2.m )
	if err != nil {
		return
	}

	var (
		sum byte
	)

	for i := 0; i < mtx.n; i++ {
		line := mtx.elts[i]

		for j := 0; j < m2.m; j++ {
			sum = 0
			for k := 0; k < mtx.m; k++ {
				sum = mtx.f.Add( sum, mtx.f.MulBy( line[k], m2.elts[k][j] ) )
			}
			rst.elts[i][j] = sum
		}
	}
	return
}

func ( mtx *Matrix8 ) MulInplace( m *Matrix8 ) ( err error ) {
	if mtx.n != mtx.m || mtx.n != m.n || m.n != m.m {
		return ErrDifferenSize
	}

	var (
		line = make( []byte, mtx.m )
		sum byte
	)

	for i := 0; i < mtx.n; i++ {
		for j := 0; j < mtx.m; j++ {
			line[j] = mtx.elts[i][j]
		}

		for j := 0; j < mtx.m; j++ {
			sum = 0
			for k := 0; k < m.n; k++ {
				sum = mtx.f.Add( sum, mtx.f.MulBy( line[k], m.elts[k][j] ) )
			}
			mtx.elts[i][j] = sum
		}
	}
	return
}
func ( mtx *Matrix8 ) Inverse() ( imtx *Matrix8, err error ) {
	cp := mtx.Copy()
	imtx, err = NewUnit( mtx.n )
	if err != nil {
		return
	}

	es := cp.elts
	f := cp.f

	// fmt.Printf( "invese of ------ \n" )
	// cp.Print()

	// x x x     y y y
	// x x x --> 0 y y
	// x x x     0 0 y
	for i := 0; i < cp.n; i++ {
		// make it:
		//	elts[i][i] != 0
		//	elts[i][j] == 0 for j > i
		if es[i][i] == cp.f.Zero {
			found := false
			j := 0
			for j = i+1; j < cp.m; j++ {
				if es[j][i] != cp.f.Zero {
					found = true
					break
				}
			}
			if found {
				cp.swapRow( i, j )
				imtx.swapRow( i, j )
			} else {
				err = ErrDetZero
				return
			}
		}
		// fmt.Printf( "m[%d][%d]=%02x should be non-zero\n", i, i, es[i][i] )
		// cp.Print()
		// imtx.Print()

		// es[i][i] != 0

		// make es[i][i] == 1
		if es[i][i] != f.E {
			r := f.Div( f.E, es[i][i] )
			cp.mulRow( i, r )
			imtx.mulRow( i, r )
		}
		// fmt.Printf( "made m[%d][%d] == 1\n", i, i )
		// cp.Print()
		// imtx.Print()

		// make the column all zero except es[i][i]
		for j := i+1; j < cp.m; j++ {
			d := f.Sub( 0, es[j][i] )
			cp.addRowTo( i, j, d )
			imtx.addRowTo( i, j, d )
		}
	}

	//  y y y     y 0 0
	//  0 y y --> 0 y 0
	//  0 0 y     0 0 y
	for i := cp.n-1; i >= 0; i-- {

		// make the row all zero except es[i][i]
		for j := 0; j < i; j++ {
			d := f.Sub( 0, es[j][i] )
			cp.addRowTo( i, j, d )
			imtx.addRowTo( i, j, d )
		}
	}
	return
}
func ( mtx *Matrix8 ) swapRow( src int, dst int ) {
	es := mtx.elts

	for i := 0; i < mtx.m; i++ {
		s0 := es[src][i]
		d0 := es[dst][i]

		es[dst][i] = s0
		es[src][i] = d0
	}
}
func ( mtx *Matrix8 ) mulRow( dst int, ratio byte ) {
	f := mtx.f
	es := mtx.elts

	for i := 0; i < mtx.m; i++ {
		es[dst][i] = f.MulBy( es[dst][i], ratio )
	}
}
func ( mtx *Matrix8 ) addRowTo( src int, dst int, ratio byte ) {
	f := mtx.f
	es := mtx.elts

	for i := 0; i < mtx.m; i++ {
		s0 := es[src][i]
		d0 := es[dst][i]

		es[dst][i] = f.Add( d0, f.MulBy( s0, ratio ) )
	}
}
func ( mtx *Matrix8 ) addColTo( src int, dst int, ratio byte ) {
	f := mtx.f
	es := mtx.elts

	for i := 0; i < mtx.n; i++ {
		s0 := es[i][src]
		d0 := es[i][dst]

		es[i][dst] = f.Add( d0, f.MulBy( s0, ratio ) )
	}
}

func ( mtx *Matrix8 ) Set( i int, j int, v byte ) {
	mtx.elts[i][j] = v
}
func ( mtx *Matrix8 ) GetField() *gf8.F {
	return mtx.f
}
func ( mtx *Matrix8 ) GetSize() ( int, int ) {
	return mtx.n, mtx.m
}

func ( mtx *Matrix8 ) Print() {
	for i := 0; i < mtx.n; i++ {
		for j := 0; j < mtx.m; j++ {
			fmt.Printf( " %02x",  mtx.elts[i][j] )
		}
		fmt.Println()
	}
}
func ( mtx *Matrix8 ) CopySlice( s [][]byte ) ( err error ) {

	if mtx.n != len(s) {
		err = ErrDifferenSize
		return
	}

	for i := 0; i < mtx.n; i++ {

		if mtx.m != len(s[i]) {
			err = ErrDifferenSize
			return
		}

		for j := 0; j < mtx.m; j++ {
			mtx.elts[i][j] = s[i][j]
		}
	}

	return
}
func NewFromSlice( s [][]byte ) ( mtx *Matrix8, err error ) {
	if s == nil {
		return
	}
	mtx, err = NewMatrix( len(s), len(s[0]) )
	if err != nil {
		return
	}

	err = mtx.CopySlice( s )
	return
}
func ( mtx *Matrix8 )Equal( b *Matrix8 ) bool {

	if mtx == nil || b == nil {
		return false
	}

	if mtx.n != b.n || mtx.m != b.m {
		return false
	}

	for i := 0; i < mtx.n; i++ {
		for j := 0; j < mtx.m; j++ {
			if mtx.elts[i][j] != b.elts[i][j] {
				return false
			}
		}
	}
	return true
}
