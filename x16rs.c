#include "x16rs.h"
#include <stdlib.h>
#include <stdint.h>
#include <string.h>
#include <stdio.h>

#pragma GCC diagnostic ignored "-Wpointer-sign"

void print_byte_list(char* name, void *data, int len, int wide);

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


void print_byte_list(char* name, void *data, int len, int wide)
{
    if (len > 0) {
        int i;
        printf("%s: %d", name, ((uint8_t*)data)[0]);
        for(i = 1; i < len; i++){
            if (wide > 0 && i % wide == 0)
				printf("\n");
            printf(",%d", ((uint8_t*)data)[i]);
        }
    }else{
        printf("%s", name);
    }

    printf("\n");
    fflush(stdout);
}

void sha3_256(const char *input, const int in_size, char *output)
{
    sha3_ctx ctx;
    rhash_sha3_256_init(&ctx);
    rhash_sha3_update(&ctx, input, in_size);
    rhash_sha3_final(&ctx, output);
}

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
    uint32_t hash[64 / 4];
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
            case SKEIN:
            sph_skein512_init(&ctx_skein);
            sph_skein512(&ctx_skein, in, size);
            sph_skein512_close(&ctx_skein, hash);
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
void x16rs_hash_old(const int loopnum, const char* input, char* output)
{
    int insize = 32;
    int k = 0;
	int n = 0;

    uint8_t in[32];
    uint8_t in_step[32];
    uint8_t out[32];
    memset(in,  0, 32);
    memset(in_step,  0, 32);
    memset(out, 0, 32);
    memcpy((void*)in,      input, 32); // in
    memcpy((void*)in_step, input, 32); // in_step
    // in[0] = 1; // Test for Correctness

    for(k = 0; k < 2 * loopnum; k++){
		/* Execute the hash algorithm, the first step */
        x16rs_hash_sz((char*)in_step, (char*)in_step, insize);
    }

    for(n = 0; n < loopnum; n++){
        memset(out, 0, 32);
		/* Execute the hash algorithm, the second step */
        x16rs_hash_sz((char*)in, (char*)out, insize);
        if(out[0]+out[1]+out[30]+out[31] != 0){
            memcpy((void*)in, out, 32); // output => in
        }
    }

    // return
    memcpy(output, out, 32);
}



// input length must more than 32
void x16rs_hash_sz(const char* input, char* output, int insize)
{
    uint32_t hash[64/4];

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

    uint8_t in[32];
    memset(in, 0, 32); // init
    memcpy((void*)in,   input, 32); // in
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

        memset(in, 0, 32);
        memcpy((void*)in, hash, 32);
    }

    uint8_t results[32];
    memset(results, 0, 32); // memset results to zero
    memcpy((void*)results,  hash, 32);
    memcpy(output, results, 32);
}


// input length must more than 32
static const size_t x16rs_hash_insize = 32;
void x16rs_hash(const int loopnum, const char* input_hash, char* output_hash)
{
    uint32_t inputoutput[64/4];

    uint32_t *input_hash_ptr32 = (uint32_t *) input_hash;
    inputoutput[0] = input_hash_ptr32[0];
    inputoutput[1] = input_hash_ptr32[1];
    inputoutput[2] = input_hash_ptr32[2];
    inputoutput[3] = input_hash_ptr32[3];
    inputoutput[4] = input_hash_ptr32[4];
    inputoutput[5] = input_hash_ptr32[5];
    inputoutput[6] = input_hash_ptr32[6];
    inputoutput[7] = input_hash_ptr32[7];

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

    int n;
    for(n = 0; n < loopnum; n++){

        uint8_t algo = inputoutput[7] % 16;
        switch (algo)
        {
        case BLAKE:
            sph_blake512_init(&ctx_blake);
            sph_blake512(&ctx_blake, inputoutput, x16rs_hash_insize);
            sph_blake512_close(&ctx_blake, inputoutput);
        break;
        case BMW:
            sph_bmw512_init(&ctx_bmw);
            sph_bmw512(&ctx_bmw, inputoutput, x16rs_hash_insize);
            sph_bmw512_close(&ctx_bmw, inputoutput);
        break;
        case GROESTL:
            sph_groestl512_init(&ctx_groestl);
            sph_groestl512(&ctx_groestl, inputoutput, x16rs_hash_insize);
            sph_groestl512_close(&ctx_groestl, inputoutput);
        break;
        case SKEIN:
            sph_skein512_init(&ctx_skein);
            sph_skein512(&ctx_skein, inputoutput, x16rs_hash_insize);
            sph_skein512_close(&ctx_skein, inputoutput);
        break;
        case JH:
            sph_jh512_init(&ctx_jh);
            sph_jh512(&ctx_jh, inputoutput, x16rs_hash_insize);
            sph_jh512_close(&ctx_jh, inputoutput);
        break;
        case KECCAK:
            sph_keccak512_init(&ctx_keccak);
            sph_keccak512(&ctx_keccak, inputoutput, x16rs_hash_insize);
            sph_keccak512_close(&ctx_keccak, inputoutput);
        break;
        case LUFFA:
            sph_luffa512_init(&ctx_luffa);
            sph_luffa512(&ctx_luffa, inputoutput, x16rs_hash_insize);
            sph_luffa512_close(&ctx_luffa, inputoutput);
        break;
        case CUBEHASH:
            sph_cubehash512_init(&ctx_cubehash);
            sph_cubehash512(&ctx_cubehash, inputoutput, x16rs_hash_insize);
            sph_cubehash512_close(&ctx_cubehash, inputoutput);
        break;
        case SHAVITE:
            sph_shavite512_init(&ctx_shavite);
            sph_shavite512(&ctx_shavite, inputoutput, x16rs_hash_insize);
            sph_shavite512_close(&ctx_shavite, inputoutput);
        break;
        case SIMD:
            sph_simd512_init(&ctx_simd);
            sph_simd512(&ctx_simd, inputoutput, x16rs_hash_insize);
            sph_simd512_close(&ctx_simd, inputoutput);
        break;
        case ECHO:
            sph_echo512_init(&ctx_echo);
            sph_echo512(&ctx_echo, inputoutput, x16rs_hash_insize);
            sph_echo512_close(&ctx_echo, inputoutput);
        break;
        case HAMSI:
            sph_hamsi512_init(&ctx_hamsi);
            sph_hamsi512(&ctx_hamsi, inputoutput, x16rs_hash_insize);
            sph_hamsi512_close(&ctx_hamsi, inputoutput);
        break;
        case FUGUE:
            sph_fugue512_init(&ctx_fugue);
            sph_fugue512(&ctx_fugue, inputoutput, x16rs_hash_insize);
            sph_fugue512_close(&ctx_fugue, inputoutput);
        break;
        case SHABAL:
            sph_shabal512_init(&ctx_shabal);
            sph_shabal512(&ctx_shabal, inputoutput, x16rs_hash_insize);
            sph_shabal512_close(&ctx_shabal, inputoutput);
        break;
        case WHIRLPOOL:
            sph_whirlpool_init(&ctx_whirlpool);
            sph_whirlpool(&ctx_whirlpool, inputoutput, x16rs_hash_insize);
            sph_whirlpool_close(&ctx_whirlpool, inputoutput);
        break;
        case SHA512:
            sph_sha512_init(&ctx_sha512);
            sph_sha512(&ctx_sha512, inputoutput, x16rs_hash_insize);
            sph_sha512_close(&ctx_sha512, inputoutput);
        break;
        }

    }

    uint32_t *output_hash_ptr32 = (uint32_t *) output_hash;
    output_hash_ptr32[0] = inputoutput[0];
    output_hash_ptr32[1] = inputoutput[1];
    output_hash_ptr32[2] = inputoutput[2];
    output_hash_ptr32[3] = inputoutput[3];
    output_hash_ptr32[4] = inputoutput[4];
    output_hash_ptr32[5] = inputoutput[5];
    output_hash_ptr32[6] = inputoutput[6];
    output_hash_ptr32[7] = inputoutput[7];
}

// input length must be 32
static const uint8_t diamond_hash_base_stuff[17] = "0WTYUIAHXVMEKBSZN";
void diamond_hash(const char* hash32, char* output16)
{
    uint8_t *stuff32 = (uint8_t*)hash32;
    int i, p = 13;
    for(i = 0; i < 16; i++){
        int num = p * (int)stuff32[i*2] * (int)stuff32[i*2+1];
        p = num % 17;
        output16[i] = diamond_hash_base_stuff[p];
        if(p == 0){
           p = 13; 
        }
    }
}

// input length must be 32
void miner_diamond_hash(const uint32_t hsstart, const uint32_t hsend, const int diamondnumber, const char* stop_mark1, const char* input32, const char* addr21, const char* extmsg32, char* nonce8, char* diamond16)
{
    int loopnum = diamondnumber / 8192 + 1; // Adjust the hashing times every 8192 diamonds (about 140 days and half a year)
    if( loopnum > 16 ){
        loopnum = 16; // atmost 16 rounds for x16rs algorithm
    }

    // mark for stop running or not
    uint8_t *is_stop = (uint8_t*)stop_mark1;

    uint32_t basestuff[8+2+6+8];
    int basestufftargetsize = 61;
    memcpy( (void*)basestuff, (void*)input32, 32);
    memcpy( (void*)basestuff+40, (void*)addr21, 21);

    // If 20000 diamonds have been produced, basestuff should add extend msg
    if(diamondnumber > 20000){
        memcpy( (void*)basestuff+40+21, (void*)extmsg32, 32);
        basestufftargetsize = 61 + 32;
    }

    // start calculate
    uint8_t sha3res[32];
    uint8_t hashnew[32];
    uint8_t diamond[16];

    uint32_t noncenum1;
    uint32_t noncenum2;
    for(noncenum1 = hsstart; noncenum1 < hsend; noncenum1++){
        basestuff[8] = noncenum1;
        for(noncenum2 = 1; noncenum2 < 4294967294; noncenum2++){
            basestuff[9] = noncenum2;

            // stop if get stop mark
            if(noncenum2 % 1000 == 0 && is_stop[0] != 0) {

                uint8_t noncenum_empty[8] = {0,0,0,0,0,0,0,0};
                memcpy(nonce8, noncenum_empty, 8);

                //////// TEST START ////////
                // uint32_t nonce[2] = {0, 0};
                // nonce[0] = noncenum1;
                // nonce[1] = noncenum2-1;
                // memcpy( (char*)nonce8, (char*)nonce, 8);
                // memcpy( (char*)diamond16, (char*)diamond, 16); 
                //////// TEST END //////// 

                return; // return null
            }
            // running hash calculate
            sha3_256((char*)basestuff, basestufftargetsize, (char*)sha3res);
            x16rs_hash(loopnum, (char*)sha3res, (char*)hashnew);
            diamond_hash((char*)hashnew, (char*)diamond);

            // to check it result is diamond
            uint8_t success = 1;
            int k;
			// Diamond is 16 bits. Here, directly traverse the string to determine whether it is a diamond
			// The rule of judgment is that the first 10 digits must be 0, and the last 6 digits cannot be 0
            for(k = 0; k < 16; k++) {
                if(k < 10){
					// The number 48 here represents the ASCII character '0'
                    if (diamond[k] != 48) {
                        success = 0;
                        break;
                    }
                } else {
                    if (diamond[k] == 48) {
                        success = 0;
                        break;
                    }
                }
            }

            // if mint diamond success
            if (success == 1) {
                // check difficulty
                // Referring to Moore's law, the excavation difficulty of every 42000 diamonds will double in about 2 years,
				// and the difficulty increment will tend to decrease to zero in 64 years
                uint8_t diadiffbits[32] = {
                    128,132,136,140,144,148,152,156, // Step + 4
                	160,164,168,172,176,180,184,188,
                	192,196,200,204,208,212,216,220,
                	224,228,232,236,240,244,248,252
                };
                int shnum = diamondnumber / 42000; // max 32step atmost 64 years
                int shmax = 255 - (diamondnumber / 65536);
                uint8_t diffyes = 1;
                int i;
                for (i = 0; i < 32; i++) {
                    if (i < shnum && sha3res[i] >= diadiffbits[i]) {
                        diffyes = 0; // Check failed, difficulty value does not meet requirements
                        break;
                    }
                    if (sha3res[i] > shmax) {
                        diffyes = 0; // Check failed, difficulty value does not meet requirements
                        break;
                    }
                }
                // difficulty value does not meet requirements
                if (diffyes != 1) {
                    continue; // to do next turn mining
                }
                // Adjust the difficulty for every 3277 diamonds 3277 = 16 ^ 6 / 256 / 20
                // When the difficulty is the highest, the first 20 bits of hash are 0, not all 32 bits are 0
                uint8_t diffok = 0;
                int diffnum = diamondnumber / 3277; // Adjust the difficulty for every 3277 diamonds
                int n;
                for (n = 0; n < 32; n++) {
                    int bt = hashnew[n];
                    if (diffnum < 255) {
                        if (bt + diffnum > 255) {
                            diffok = 0;
                            break; // Difficulty check failed
                        } else {
                            diffok = 1;
                            break; // success
                        }
                    } else if (diffnum >= 255) {
                        if (bt != 0) {
                            diffok = 0;
                            break; // Difficulty check failed
                        }
                        // next turn check
                        diffnum -= 255;
                    }
                }
                // difficulty value does not meet requirements
                if (diffok != 1) {
                    // next turn mining
                    continue;
                }

                // Difficulty check passed
                uint32_t nonce[2] = {0, 0};
                nonce[0] = noncenum1;
                nonce[1] = noncenum2;
                memcpy((char*)nonce8, (char*)nonce, 8);
                memcpy((char*)diamond16, (char*)diamond, 16);
                // printf("\n%s\n", diamond); fflush(stdout);
                return; // copy diamond to diamond16, copy nonce to nonce8 and return success
            }

        }

    }

    // Loop completed, return null, failed
    uint8_t noncenum_empty[8] = {0,0,0,0,0,0,0,0};
    memcpy( nonce8, noncenum_empty, 8);

}

// algorithm to mining hac
void miner_x16rs_hash(const int loopnum, const int retmaxhash, const char* stop_mark1, const uint32_t hsstart, const uint32_t hsend, const char* target_difficulty_hash32, const char* input_stuff89, char* stopkind1, char* success1, char* nonce4, char* reshash32)
{
//    printf("miner_x16rs_hash_v1()\n");
    // Signature information
    uint8_t stuffnew_base[90];
    uint8_t *stuffnew = stuffnew_base + 1;
    memcpy(stuffnew, input_stuff89, 89);
    uint32_t *stuffnew_uint32 = (uint32_t*)stuffnew_base;

    // value for sha3 result
    unsigned char sha3res[32];

    // value for x16rs hash
    uint8_t hashnew[32];
    uint8_t hashmaxpower[32];
    hashmaxpower[0] = 255; // init
    uint8_t noncemaxpower[4];

    // Judge whether the hash meets the requirements
    uint8_t iscalcok = 0;
    uint8_t ispowerok = 0;

    // Check the for loop value of the result
    uint8_t pk = 0;
    uint8_t pi = 0;

    // stop mining mark
    uint8_t *is_stop = (uint8_t*)stop_mark1;

    // nonce value
    uint32_t noncenum;
    for (noncenum = hsstart; noncenum < hsend; noncenum++) {
        // stop mining
        if (noncenum % 5000 == 0 && is_stop[0] != 0)
        {
            success1[0] = 0; // Record the mining failure
            stopkind1[0] = 1; // External signal forced stop mining hac
            if (retmaxhash == 1) {
                memcpy(nonce4, &noncemaxpower, 4);
                memcpy(reshash32, &hashmaxpower, 32);
            } else {
                memcpy(nonce4, &noncenum, 4);
                memcpy(reshash32, &hashnew, 32);
            }
            return;
        }

        // reset nonce value
        stuffnew_uint32[20] = noncenum;
        sha3_256(stuffnew, 89, sha3res);
        x16rs_hash(loopnum, ((void*)sha3res), ((void*)hashnew));

        if (retmaxhash == 1) {
            ispowerok = 0;
            for (pk = 0; pk < 32; pk++) {
                uint8_t o1 = hashmaxpower[pk];
                uint8_t o2 = hashnew[pk];
                if (o2 > o1) {
                    break;
                }
                if (o2 < o1) {
                    ispowerok = 1;
                    break;
                }
            }
            if (ispowerok == 1) {
                memcpy(&noncemaxpower, &noncenum, 4);
                memcpy(&hashmaxpower, &hashnew, 32);
            } else {
                continue;
            }
        }

        iscalcok = 0;
        for (pi = 0; pi < 32; pi++) {
            uint8_t o1 = target_difficulty_hash32[pi];
            uint8_t o2 = hashnew[pi];
            if (o2 > o1) {
                break;
            }
            if (o2 < o1) {
                iscalcok = 1;
                break;
            }
        }
        // check mining result
        if (iscalcok == 1) {
            // return copy to out
            success1[0] = 1; // success mark
            stopkind1[0] = 2; // Non forced stop, stop after successful excavation
            if (retmaxhash == 1) {
                memcpy(nonce4, &noncemaxpower, 4);
                memcpy(reshash32, &hashmaxpower, 32);
            } else {
                memcpy(nonce4, &noncenum, 4);
                memcpy(reshash32, &hashnew, 32);
            }
            return; // success
        }
        // continue doing next excavation
    }

    // Cycle complete, no successful mining
    success1[0] = 0;
    stopkind1[0] = 0; // Non forced stop, stop after natural circulation
    if (retmaxhash == 1) {
        memcpy(nonce4, &noncemaxpower, 4);
        memcpy(reshash32, &hashmaxpower, 32);
    } else {
        memcpy(nonce4, &noncenum, 4);
        memcpy(reshash32, &hashnew, 32);
    }
    return;
}

///////////////////// TEST /////////////////////
void test_print_x16rs(const char* input, char* output32x16)
{
    int insize = 32;
    int size = insize;
    int i = 0;
    uint32_t hash[64/4];

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
    uint8_t output[32];
    
    sph_blake512_init(&ctx_blake);
    sph_blake512(&ctx_blake, in, size);
    sph_blake512_close(&ctx_blake, hash);
    in = (void*) hash;
    memcpy(output32x16 + 32*0, (void*)hash, 32);
    
    sph_bmw512_init(&ctx_bmw);
    sph_bmw512(&ctx_bmw, in, size);
    sph_bmw512_close(&ctx_bmw, hash);
    in = (void*) hash;
    memcpy(output32x16 + 32*1, (void*)hash, 32);
    
    sph_groestl512_init(&ctx_groestl);
    sph_groestl512(&ctx_groestl, in, size);
    sph_groestl512_close(&ctx_groestl, hash);
    in = (void*) hash;
    memcpy(output32x16 + 32*2, (void*)hash, 32);
    
    sph_skein512_init(&ctx_skein);
    sph_skein512(&ctx_skein, in, size);
    sph_skein512_close(&ctx_skein, hash);
    in = (void*) hash;
    memcpy(output32x16 + 32*3, (void*)hash, 32);
    
    sph_jh512_init(&ctx_jh);
    sph_jh512(&ctx_jh, in, size);
    sph_jh512_close(&ctx_jh, hash);
    in = (void*) hash;
    memcpy(output32x16 + 32*4, (void*)hash, 32);

    sph_keccak512_init(&ctx_keccak);
    sph_keccak512(&ctx_keccak, in, size);
    sph_keccak512_close(&ctx_keccak, hash);
    in = (void*) hash;
    memcpy(output32x16 + 32*5, (void*)hash, 32);
    
    sph_luffa512_init(&ctx_luffa);
    sph_luffa512(&ctx_luffa, in, size);
    sph_luffa512_close(&ctx_luffa, hash);
    in = (void*) hash;
    memcpy(output32x16 + 32*6, (void*)hash, 32);
    
    sph_cubehash512_init(&ctx_cubehash);
    sph_cubehash512(&ctx_cubehash, in, size);
    sph_cubehash512_close(&ctx_cubehash, hash);
    in = (void*) hash;
    memcpy(output32x16 + 32*7, (void*)hash, 32);
    
    sph_shavite512_init(&ctx_shavite);
    sph_shavite512(&ctx_shavite, in, size);
    sph_shavite512_close(&ctx_shavite, hash);
    in = (void*) hash;
    memcpy(output32x16 + 32*8, (void*)hash, 32);
    
    sph_simd512_init(&ctx_simd);
    sph_simd512(&ctx_simd, in, size);
    sph_simd512_close(&ctx_simd, hash);
    in = (void*) hash;
    memcpy(output32x16 + 32*9, (void*)hash, 32); 
 
    sph_echo512_init(&ctx_echo);
    sph_echo512(&ctx_echo, in, size);
    sph_echo512_close(&ctx_echo, hash);
    in = (void*) hash;
    memcpy(output32x16 + 32*10, (void*)hash, 32);
    
    sph_hamsi512_init(&ctx_hamsi);
    sph_hamsi512(&ctx_hamsi, in, size);
    sph_hamsi512_close(&ctx_hamsi, hash);
    in = (void*) hash;
    memcpy(output32x16 + 32*11, (void*)hash, 32);  
    
    sph_fugue512_init(&ctx_fugue);
    sph_fugue512(&ctx_fugue, in, size);
    sph_fugue512_close(&ctx_fugue, hash);
    in = (void*) hash;
    memcpy(output32x16 + 32*12, (void*)hash, 32);
    
    sph_shabal512_init(&ctx_shabal);
    sph_shabal512(&ctx_shabal, in, size);
    sph_shabal512_close(&ctx_shabal, hash);
    in = (void*) hash;
    memcpy(output32x16 + 32*13, (void*)hash, 32);
    
    sph_whirlpool_init(&ctx_whirlpool);
    sph_whirlpool(&ctx_whirlpool, in, size);
    sph_whirlpool_close(&ctx_whirlpool, hash);
    in = (void*) hash;
    memcpy(output32x16 + 32*14, (void*)hash, 32);
    
    sph_sha512_init(&ctx_sha512);
    sph_sha512(&ctx_sha512,(const void*) in, size);
    sph_sha512_close(&ctx_sha512,(void*) hash);
    in = (void*) hash;
    memcpy(output32x16 + 32*15, (void*)hash, 32);
}
