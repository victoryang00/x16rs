#include "x16rs.h"
#include <stdlib.h>
#include <stdint.h>
#include <string.h>
#include <stdio.h>

#include "sha3/sph_blake.h"
#include "sha3/sph_bmw.h"
#include "sha3/sph_groestl.h"
#include "sha3/sph_jh.h"
#include "sha3/sph_keccak.h"
#include "sha3/sph_skein.h"
#include "sha3/sph_luffa.h"
#include "sha3/sph_cubehash.h"
#include "sha3/sph_shavite.h"
#include "sha3/sph_simd.h"
#include "sha3/sph_echo.h"
#include "sha3/sph_hamsi.h"
#include "sha3/sph_fugue.h"
#include "sha3/sph_shabal.h"
#include "sha3/sph_whirlpool.h"
#include "sha3/sph_sha2.h"
#include "sha3_256/sha3.h"

enum Algo {
    BLAKE = 0,
    BMW,
    GROESTL,
    JH,
    KECCAK,
    SKEIN,
    LUFFA,
    CUBEHASH,
    SHAVITE,
    SIMD,
    ECHO,
    HAMSI,
    FUGUE,
    SHABAL,
    WHIRLPOOL,
    SHA512,
    HASH_FUNC_COUNT
};

static void getAlgoString(const uint8_t* prevblock, char *output)
{
    char *sptr = output;
    int j;
    for (j = 0; j < HASH_FUNC_COUNT; j++) {
        char b = (15 - j) >> 1; // 16 ascii hex chars, reversed
        uint8_t algoDigit = (j & 1) ? prevblock[b] & 0xF : prevblock[b] >> 4;
        if (algoDigit >= 10)
            sprintf(sptr, "%c", 'A' + (algoDigit - 10));
        else
            sprintf(sptr, "%u", (uint32_t) algoDigit);
        sptr++;
    }
    *sptr = '\0';
}

void x16r_hash(const char* input, char* output)
{
    uint32_t hash[64/4];
    char hashOrder[HASH_FUNC_COUNT + 1] = { 0 };

    sph_blake512_context     ctx_blake;
    sph_bmw512_context       ctx_bmw;
    sph_groestl512_context   ctx_groestl;
    sph_skein512_context     ctx_skein;
    sph_jh512_context        ctx_jh;
    sph_keccak512_context    ctx_keccak;
    sph_luffa512_context     ctx_luffa;
    sph_cubehash512_context  ctx_cubehash;
    sph_shavite512_context   ctx_shavite;
    sph_simd512_context      ctx_simd;
    sph_echo512_context      ctx_echo;
    sph_hamsi512_context     ctx_hamsi;
    sph_fugue512_context     ctx_fugue;
    sph_shabal512_context    ctx_shabal;
    sph_whirlpool_context    ctx_whirlpool;
    sph_sha512_context       ctx_sha512;

    void *in = (void*) input;
    int size = 80;

    getAlgoString((uint8_t*)&input[4], hashOrder);
    
    int i;
    for (i = 0; i < 16; i++) {
        const char elem = hashOrder[i];
        const uint8_t algo = elem >= 'A' ? elem - 'A' + 10 : elem - '0';

        switch (algo) {
            case BLAKE:
            sph_blake512_init(&ctx_blake);
            sph_blake512(&ctx_blake, in, size);
            sph_blake512_close(&ctx_blake, hash);
            break;
            case BMW:
            sph_bmw512_init(&ctx_bmw);
            sph_bmw512(&ctx_bmw, in, size);
            sph_bmw512_close(&ctx_bmw, hash);
            break;
            case GROESTL:
            sph_groestl512_init(&ctx_groestl);
            sph_groestl512(&ctx_groestl, in, size);
            sph_groestl512_close(&ctx_groestl, hash);
            break;
            case SKEIN:
            sph_skein512_init(&ctx_skein);
            sph_skein512(&ctx_skein, in, size);
            sph_skein512_close(&ctx_skein, hash);
            break;
            case JH:
            sph_jh512_init(&ctx_jh);
            sph_jh512(&ctx_jh, in, size);
            sph_jh512_close(&ctx_jh, hash);
            break;
            case KECCAK:
            sph_keccak512_init(&ctx_keccak);
            sph_keccak512(&ctx_keccak, in, size);
            sph_keccak512_close(&ctx_keccak, hash);
            break;
            case LUFFA:
            sph_luffa512_init(&ctx_luffa);
            sph_luffa512(&ctx_luffa, in, size);
            sph_luffa512_close(&ctx_luffa, hash);
            break;
            case CUBEHASH:
            sph_cubehash512_init(&ctx_cubehash);
            sph_cubehash512(&ctx_cubehash, in, size);
            sph_cubehash512_close(&ctx_cubehash, hash);
            break;
            case SHAVITE:
            sph_shavite512_init(&ctx_shavite);
            sph_shavite512(&ctx_shavite, in, size);
            sph_shavite512_close(&ctx_shavite, hash);
            break;
            case SIMD:
            sph_simd512_init(&ctx_simd);
            sph_simd512(&ctx_simd, in, size);
            sph_simd512_close(&ctx_simd, hash);
            break;
            case ECHO:
            sph_echo512_init(&ctx_echo);
            sph_echo512(&ctx_echo, in, size);
            sph_echo512_close(&ctx_echo, hash);
            break;
            case HAMSI:
            sph_hamsi512_init(&ctx_hamsi);
            sph_hamsi512(&ctx_hamsi, in, size);
            sph_hamsi512_close(&ctx_hamsi, hash);
            break;
            case FUGUE:
            sph_fugue512_init(&ctx_fugue);
            sph_fugue512(&ctx_fugue, in, size);
            sph_fugue512_close(&ctx_fugue, hash);
            break;
            case SHABAL:
            sph_shabal512_init(&ctx_shabal);
            sph_shabal512(&ctx_shabal, in, size);
            sph_shabal512_close(&ctx_shabal, hash);
            break;
            case WHIRLPOOL:
            sph_whirlpool_init(&ctx_whirlpool);
            sph_whirlpool(&ctx_whirlpool, in, size);
            sph_whirlpool_close(&ctx_whirlpool, hash);
            break;
            case SHA512:
            sph_sha512_init(&ctx_sha512);
            sph_sha512(&ctx_sha512,(const void*) in, size);
            sph_sha512_close(&ctx_sha512,(void*) hash);
            break;
        }
        in = (void*) hash;
        size = 64;
    }
    memcpy(output, hash, 32);
}

// input length must more than 32
void x16rs_hash(const char* input, char* output)
{
    int insize = 32;
    x16rs_hash_sz(input, output, insize);
}
// input length must more than 32
void x16rs_hash_sz(const char* input, char* output, int insize)
{
//    printf("%s\n", input);

    uint32_t hash[64/4];

    // char hashOrder[HASH_FUNC_COUNT + 1] = { 0 };

    sph_blake512_context     ctx_blake;
    sph_bmw512_context       ctx_bmw;
    sph_groestl512_context   ctx_groestl;
    sph_skein512_context     ctx_skein;
    sph_jh512_context        ctx_jh;
    sph_keccak512_context    ctx_keccak;
    sph_luffa512_context     ctx_luffa;
    sph_cubehash512_context  ctx_cubehash;
    sph_shavite512_context   ctx_shavite;
    sph_simd512_context      ctx_simd;
    sph_echo512_context      ctx_echo;
    sph_hamsi512_context     ctx_hamsi;
    sph_fugue512_context     ctx_fugue;
    sph_shabal512_context    ctx_shabal;
    sph_whirlpool_context    ctx_whirlpool;
    sph_sha512_context       ctx_sha512;

    void *in = (void*) input;

    memcpy((void*)hash, input, 32); // first

    int size = insize;

    int i;
    for(i = 0; i < 1; i++) {

        uint8_t algo = hash[7] % 16;

        switch (algo) {
            case BLAKE:
            sph_blake512_init(&ctx_blake);
            sph_blake512(&ctx_blake, in, size);
            sph_blake512_close(&ctx_blake, hash);
            /*printf("switch (algo) = %d \n", algo);
            for(int i=0; i<8; i++){
                printf("hash %d = %d\n", i, hash[i]);
            }*/
            break;
            case BMW:
            sph_bmw512_init(&ctx_bmw);
            sph_bmw512(&ctx_bmw, in, size);
            sph_bmw512_close(&ctx_bmw, hash);
            /*printf("switch (algo) = %d \n", algo);
            for(int i=0; i<8; i++){
                printf("hash %d = %d\n", i, hash[i]);
            }*/
            break;
            case GROESTL:
            sph_groestl512_init(&ctx_groestl);
            sph_groestl512(&ctx_groestl, in, size);
            sph_groestl512_close(&ctx_groestl, hash);
            /*printf("switch (algo) = %d \n", algo);
            for(int i=0; i<8; i++){
                printf("hash %d = %d\n", i, hash[i]);
            }*/
            break;
            case SKEIN:
            sph_skein512_init(&ctx_skein);
            sph_skein512(&ctx_skein, in, size);
            sph_skein512_close(&ctx_skein, hash);
            /*printf("switch (algo) = %d \n", algo);
            for(int i=0; i<8; i++){
                printf("hash %d = %d\n", i, hash[i]);
            }*/
            break;
            case JH:
            sph_jh512_init(&ctx_jh);
            sph_jh512(&ctx_jh, in, size);
            sph_jh512_close(&ctx_jh, hash);
            /*printf("switch (algo) = %d \n", algo);
            for(int i=0; i<8; i++){
                printf("hash %d = %d\n", i, hash[i]);
            }*/
            break;
            case KECCAK:
            sph_keccak512_init(&ctx_keccak);
            sph_keccak512(&ctx_keccak, in, size);
            sph_keccak512_close(&ctx_keccak, hash);
            /*printf("switch (algo) = %d \n", algo);
            for(int i=0; i<8; i++){
                printf("hash %d = %d\n", i, hash[i]);
            }*/
            break;
            case LUFFA:
            sph_luffa512_init(&ctx_luffa);
            sph_luffa512(&ctx_luffa, in, size);
            sph_luffa512_close(&ctx_luffa, hash);
            /*printf("switch (algo) = %d \n", algo);
            for(int i=0; i<8; i++){
                printf("hash %d = %d\n", i, hash[i]);
            }*/
            break;
            case CUBEHASH:
            sph_cubehash512_init(&ctx_cubehash);
            sph_cubehash512(&ctx_cubehash, in, size);
            sph_cubehash512_close(&ctx_cubehash, hash);
            /*printf("switch (algo) = %d \n", algo);
            for(int i=0; i<8; i++){
                printf("hash %d = %d\n", i, hash[i]);
            }*/
            break;
            case SHAVITE:
            sph_shavite512_init(&ctx_shavite);
            sph_shavite512(&ctx_shavite, in, size);
            sph_shavite512_close(&ctx_shavite, hash);
            /*printf("switch (algo) = %d \n", algo);
            for(int i=0; i<8; i++){
                printf("hash %d = %d\n", i, hash[i]);
            }*/
            break;
            case SIMD:
            sph_simd512_init(&ctx_simd);
            sph_simd512(&ctx_simd, in, size);
            sph_simd512_close(&ctx_simd, hash);
            /*printf("switch (algo) = %d \n", algo);
            for(int i=0; i<8; i++){
                printf("hash %d = %d\n", i, hash[i]);
            }*/
            break;
            case ECHO:
            sph_echo512_init(&ctx_echo);
            sph_echo512(&ctx_echo, in, size);
            sph_echo512_close(&ctx_echo, hash);
            /*printf("switch (algo) = %d \n", algo);
            for(int i=0; i<8; i++){
                printf("hash %d = %d\n", i, hash[i]);
            }*/
            break;
            case HAMSI:
            sph_hamsi512_init(&ctx_hamsi);
            sph_hamsi512(&ctx_hamsi, in, size);
            sph_hamsi512_close(&ctx_hamsi, hash);
            /*printf("switch (algo) = %d \n", algo);
            for(int i=0; i<8; i++){
                printf("hash %d = %d\n", i, hash[i]);
            }*/
            break;
            case FUGUE:
            sph_fugue512_init(&ctx_fugue);
            sph_fugue512(&ctx_fugue, in, size);
            sph_fugue512_close(&ctx_fugue, hash);
            /*printf("switch (algo) = %d \n", algo);
            for(int i=0; i<8; i++){
                printf("hash %d = %d\n", i, hash[i]);
            }*/
            break;
            case SHABAL:
            sph_shabal512_init(&ctx_shabal);
            sph_shabal512(&ctx_shabal, in, size);
            sph_shabal512_close(&ctx_shabal, hash);
            /*printf("switch (algo) = %d \n", algo);
            for(int i=0; i<8; i++){
                printf("hash %d = %d\n", i, hash[i]);
            }*/
            break;
            case WHIRLPOOL:
            sph_whirlpool_init(&ctx_whirlpool);
            sph_whirlpool(&ctx_whirlpool, in, size);
            sph_whirlpool_close(&ctx_whirlpool, hash);
            /*printf("switch (algo) = %d \n", algo);
            for(int i=0; i<8; i++){
                printf("hash %d = %d\n", i, hash[i]);
            }*/
            break;
            case SHA512:
            sph_sha512_init(&ctx_sha512);
            sph_sha512(&ctx_sha512,(const void*) in, size);
            sph_sha512_close(&ctx_sha512,(void*) hash);
            /*printf("switch (algo) = %d \n", algo);
            for(int i=0; i<8; i++){
                printf("hash %d = %d\n", i, hash[i]);
            }*/
            break;
        }
        in = (void*) hash;
    }

    memcpy(output, hash, 32);

}



// input length must be 32
static const uint8_t diamond_hash_base_stuff[17] = "0WTYUIAHXVMEKBSZN";
void diamond_hash(const char* blkhash32, const char* addr21, const char* nonce8, char* output16)
{

    int insize = 32+21+8;
    uint8_t input[insize];
    memcpy(input, blkhash32, 32);
    memcpy(input+32, addr21, 21);
    memcpy(input+32+21, nonce8, 8);

//    printf("input: ");
//    for(int i=0; i<insize; i++){
//        printf("%d,", input[i]);
//    }
//    printf("\n");

    uint8_t output[32];
    x16rs_hash_sz(((void*)input), ((void*)output), insize);

//    printf("output: ");
//    for(int i=0; i<32; i++){
//        printf("%d,", output[i]);
//    }
//    printf("\n");

    //
    uint16_t diamond[16];
    memcpy(diamond, output, 32);
    int i;
    for(i=0; i<16; i++){
        output16[i] = diamond_hash_base_stuff[diamond[i] % 17];
    }

    // ok


}


// input length must be 32
void miner_diamond_hash(const char* input32, const char* addr21, char* nonce8, char* output16)
{
    int zerofront = 8;

    uint8_t diamond[16];

    uint8_t noncenum[8] = {0,0,0,0,0,0,0,0};
    uint8_t i0;
    for(i0=0; i0<255; i0++){
    noncenum[0] = i0;
    uint8_t i1;
    for(i1=0; i1<255; i1++){
    noncenum[1] = i1;
    uint8_t i2;
    for(i2=0; i2<255; i2++){
    noncenum[2] = i2;
    uint8_t i3;
    for(i3=0; i3<255; i3++){
    noncenum[3] = i3;
    uint8_t i4;
    for(i4=0; i4<255; i4++){
    noncenum[4] = i4;
    uint8_t i5;
    for(i5=0; i5<255; i5++){
    noncenum[5] = i5;
    uint8_t i6;
    for(i6=0; i6<255; i6++){
    noncenum[6] = i6;
    uint8_t i7;
    for(i7=0; i7<255; i7++){
    noncenum[7] = i7;




        diamond_hash(input32, addr21, noncenum, diamond);
        uint8_t isok = 1;
        uint8_t isnchar = 0;
        int k;
        for( k=0; k<16; k++ ) {
            if( k<zerofront && diamond[k] != 48 ){
                isok = 0;
                break;
            }
            if( k>=zerofront ) {
                if( diamond[k] == 48 ){
                    if( isnchar == 1 ){
                        isok = 0;
                        break;
                    }
                }else{
                    isnchar = 1;
                }
            }
        }


        if(isok == 1){


            memcpy(nonce8, noncenum, 8);
            memcpy(output16, diamond, 16);

            int i;

            printf("hash: ");
            uint8_t *input32p = (uint8_t*)input32;
            for(i=0; i<32; i++){
                printf("%u,", input32p[i]);
            }
            printf("  addr: ");
            uint8_t *addr21p = (uint8_t*)addr21;
            for(i=0; i<21; i++){
                printf("%u,", addr21p[i]);
            }
            printf("  nonce: ");
            for(i=0; i<8; i++){
                printf("%u,", noncenum[i]);
            }
            uint8_t noncenum_swap[8];
            for(i=0; i<8; i++){ noncenum_swap[i] = noncenum[7-i]; }
            uint64_t *nnum = (uint64_t*)noncenum;
            uint64_t *nnum_swap = (uint64_t*)noncenum_swap;
            printf("  diamond nonce = %ld/%ld, value = %16.16s \n", *nnum_swap, *nnum, diamond);


        }




    }}}}}}}}


}


// 挖矿算法
void miner_x16rs_hash_v1(const char* stop_mark1, const char* target_difficulty_hash32, const char* input_stuff89, char* nonce4)
{
//    printf("miner_x16rs_hash_v1()\n");

    uint8_t *is_stop = (uint8_t*)stop_mark1;

    uint8_t noncenum[4] = {0,0,0,0};
    uint8_t i0;
    for(i0=0; i0<255; i0++){
    noncenum[0] = i0;
    uint8_t i1;
    for(i1=0; i1<255; i1++){
    noncenum[1] = i1;
    uint8_t i2;
    for(i2=0; i2<255; i2++){
    noncenum[2] = i2;
    uint8_t i3;
    for(i3=0; i3<255; i3++){
    noncenum[3] = i3;
        // stop now
        if( is_stop[0] != 0 )
        {
            return;
        }
        // 重置nonce
        uint8_t stuffnew[89];
        memcpy(stuffnew, input_stuff89, 89);
        memcpy(&stuffnew[79], noncenum, 4);
        // 计算 sha3
	    unsigned char sha3res[32];
        sha3(sha3res, 256, stuffnew, 89*8);
        /*
            printf("  hash: ");
            uint8_t i;
            for(i=0; i<32; i++){
                printf("%u,", sha3res[i]);
            }
            printf("\n");
        */
        // x16rs hash
        uint8_t hashnew[32];
        x16rs_hash(((void*)sha3res), ((void*)hashnew));
        /*
            printf("  hash: ");
            uint8_t i;
            for(i=0; i<32; i++){
                printf("%u,", hashnew[i]);
            }
            printf("\n");
        */
        // 判断哈希是否满足要求
        uint8_t isok = 0;
        uint8_t pi;
        for(pi=0; pi<32; pi++){
            uint8_t o1 = target_difficulty_hash32[pi];
            uint8_t o2 = hashnew[pi];
            if(o2<o1){
                isok = 1;
                break;
            }
            if(o2>o1){
                isok = 0;
                break;
            }
        }
        // 结果判断
        if(isok == 1) {
        /*
            printf("finish hash: ");
            uint8_t i;
            for(i=0; i<32; i++){
                printf("%u,", hashnew[i]);
            }
            printf("\n");
        */
            // 返回
            memcpy(nonce4, noncenum, 4);
            return; // success
        }
        // 继续下一次计算
    }}}}

}
