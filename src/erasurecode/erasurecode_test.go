package erasurecode

import (
	"fmt"
	"testing"
	// gf8 "galoisfield8"
)

type (
	tcase struct {
	}
)

func TestNew( t *testing.T ) {
	ec := NewEC( 4, 4 )

	fmt.Println( "ec:", ec )
	fmt.Println( "encMatrixT:" )
	ec.encMatrixT.Print()

	datas := []byte{ 1, 2, 3, 4 }
	codes := make( []byte, 4 )

	ec.EncTo( datas, codes )

	fmt.Println( "codes: ", codes )

	datas[1] = 0
	datas[2] = 0
	datas[3] = 0
	ec.DecTo( datas, codes, []int{ 1, 2, 3 }, []int{} )

	if 0 == 1 {} else {
		t.Errorf( "bb" )
	}

}

