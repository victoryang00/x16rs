#include <stdio.h>
#include <string.h>

#include "x16rs.h"

/*

gcc x16rs_test.c  x16rs.c sha3/blake.c sha3/bmw.c sha3/groestl.c sha3/jh.c sha3/keccak.c sha3/skein.c sha3/cubehash.c sha3/echo.c sha3/luffa.c sha3/simd.c sha3/hamsi.c sha3/hamsi_helper.c sha3/fugue.c sha3/shavite.c sha3/shabal.c sha3/whirlpool.c sha3/sha2big.c

*/

#if 0
int main () {
    const uint8_t input[33] = "http://www.tutorialspointxxx.com";

    for(int i=0; i<10; i++){
        printf("input %d = %d\n", i, input[i]);
    }
    uint8_t output[32];
    x16rs_hash(&input, &output);

    for(int i=0; i<10; i++){
        printf("output %d = %d\n", i, output[i]);
    }
    printf("output = %s\n", output);
}
#endif

#if 0
int main () {
    const uint8_t blkhash[32] = { 87, 206, 240, 151, 249, 167, 204, 12, 69, 188, 172, 99, 37, 181, 182, 229, 129, 153, 200, 25, 119, 99, 115, 76, 172, 102, 100, 232, 210, 184, 230, 62 };
    const uint8_t addr[21] = {0,22,82,215,51,169,135,248,255,189,82,120,227,122,205,117,64,162,253,1,38};

    uint8_t nonce8[8];
    uint8_t output16[16];

    miner_diamond_hash(blkhash, addr, nonce8, output16);
}
#endif
