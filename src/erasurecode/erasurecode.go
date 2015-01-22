package erasurecode

import (
	"fmt"
	"errors"
	"sync"
	"matrix"
	gf8 "galoisfield8"
)

type (
	EC struct {
		NData, NCode int
		encMatrix *matrix.Matrix8
		encMatrixT *matrix.Matrix8

		codeMatrix *matrix.Matrix8
		codeMatrixT *matrix.Matrix8
	}
)

var (
	ErrTooManyLost = errors.New( "too many data/code lost to recover" )
)

var (
	matrixInit sync.Once
)

func NewEC( nData int, nCode int ) *EC {

	matrixInit.Do( func(){ matrix.InitField() } )

	ec := EC{
		NData: nData,
		NCode: nCode,
	}

	mtx, err := matrix.NewMatrix( nCode, nData )
	// never happen
	if err != nil {
		return nil
	}

	cmtx, err := matrix.NewMatrix( nCode, nData + 1 )
	// never happen
	if err != nil {
		return nil
	}

	f := mtx.GetField()
	n, m := mtx.GetSize()

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			mtx.Set( i, j, f.Pow( gf8.Elt(i+1), gf8.Elt(j) ) )
			cmtx.Set( i, j, f.Pow( gf8.Elt(i+1), gf8.Elt(j) ) )
		}
		cmtx.Set( i, m, f.Sub( 0, f.E ) )
	}

	ec.encMatrix = mtx
	ec.encMatrixT = mtx.T()
	ec.codeMatrix = cmtx
	ec.codeMatrixT = cmtx.T()

	return &ec
}

func ( ec *EC ) EncTo( datas []byte, codes []byte ) {
	ec.encMatrixT.MulVecToBytes( datas, codes )
}

func ( ec *EC ) DecTo( datas []byte, codes []byte, datalost []int, codelost []int ) (err error) {

	// debug:
	for i := 0; i < len(datalost); i++ {
		if datas[datalost[i]] != 0 {
			panic( "data marked as broken is not 0" )
		}
	}
	for i := 0; i < len(codelost); i++ {
		if codes[codelost[i]] != 0 {
			panic( "code marked as broken is not 0" )
		}
	}

	// datavalid := inverseIndexes( len( datas ), datalost )
	codevalid := inverseIndexes( len( codes ), codelost )
	fmt.Println( "codevalid:", codevalid )
	fmt.Println( "codelost:", codelost )

	if len(datalost) > len( codevalid ) {
		err = ErrTooManyLost
		return
	}

	// need at most as many codes as datalost to recover
	codevalid = codevalid[ :len( datalost ) ]
	fmt.Println( "codevalid:", codevalid )

	// sub matrix that covers all valid codes and lost datas
	submtx :=  ec.encMatrix.SubMatrix( codevalid, datalost )
	fmt.Println( "submtx:" )
	submtx.Print()

	bb := make( []byte, 3 )
	submtx.MulByVecToBytes( []byte{ 1, 1, 1 }, bb )
	fmt.Println( "submtx * [1, 1, 1]=", bb )

	acodemtx :=  ec.encMatrixT.SubMatrix( inverseIndexes( len(datas), []int{} ), codevalid )
	fmt.Println( "acodemtx:" )
	acodemtx.Print()

	acodes := make( []byte, len(codevalid) )
	acodemtx.MulVecToBytes( datas, acodes )
	fmt.Println( "codes without data lost: ", acodes )

	f := ec.encMatrix.GetField()
	for i := 0; i < len( codevalid ); i++ {
		acodes[ i ] = byte(f.Sub( gf8.Elt(codes[ codevalid[i] ]), gf8.Elt(acodes[i]) ))
	}

	fmt.Println( "sub-ed codes:", acodes )

	imtx, err := submtx.Inverse()
	if err != nil {
		return
	}

	fmt.Println( "inversed:" )
	imtx.Print()

	adatas := make( []byte, len(datalost) )

	imtx.MulByVecToBytes( acodes, adatas )

	fmt.Println( "data recovered:", adatas )

	return
}

func tomap( n int, indexArr []int ) []bool {
	mp := make( []bool, n )

	for i := 0; i < len(indexArr); i++ {
		mp[indexArr[i]] = true
	}
	return mp
}

func toIndexes( mp []bool ) []int {
	ids := make( []int, 0, len( mp ) )
	for i := 0; i < len( mp ); i++ {
		if mp[i] == true {
			ids = append( ids, i )
		}
	}
	return ids
}

func inverseIndexes( n int, indexes []int ) []int {
	ids := make( []int, 0, n-len( indexes ) )
	mp := tomap( n, indexes )
	for i := 0; i < n; i++ {
		if mp[i] == false {
			ids = append( ids, i )
		}

	}
	return ids
}
