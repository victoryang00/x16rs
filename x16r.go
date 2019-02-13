package x16r

/*
#cgo LDFLAGS: libx16rs_hash.a
#include "x16rs.h"
*/
import "C"

import "fmt"

func Sum(data []byte) []byte {
	var res [32]C.char
	C.x16r_hash(C.CString(string(data)), &res[0])
	return []byte(C.GoStringN(&res[0], 32))
}

func HashX16RS(data []byte) []byte {
	var res [32]C.char
	C.x16rs_hash(C.CString(string(data)), &res[0])
	return []byte(C.GoStringN(&res[0], 32))
}

func main() {
	fmt.Println(Sum([]byte("test")))
	fmt.Println(HashX16RS([]byte("test")))
}

/*

//#include "x16r.h"

//#include "sha3/sph_blake.h"
//#include "sha3/sph_bmw.h"
//#include "sha3/sph_groestl.h"
//#include "sha3/sph_jh.h"
//#include "sha3/sph_keccak.h"
//#include "sha3/sph_skein.h"
//#include "sha3/sph_luffa.h"
//#include "sha3/sph_cubehash.h"
//#include "sha3/sph_shavite.h"
//#include "sha3/sph_simd.h"
//#include "sha3/sph_echo.h"
//#include "sha3/sph_hamsi.h"
//#include "sha3/sph_fugue.h"
//#include "sha3/sph_shabal.h"
//#include "sha3/sph_whirlpool.h"
//#include "sha3/sph_sha2.h"

*/
