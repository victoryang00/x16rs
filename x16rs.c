#include "x16rs.h"
#include <stdlib.h>
#include <stdint.h>
#include <string.h>
#include <stdio.h>

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
    if(len > 0){
        int i;
        printf("%s: %d", name, ((uint8_t*)data)[0]);
        for(i=1; i<len; i++){
            if(wide>0 && i%wide==0) printf("\n");
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
void x16rs_hash(const int loopnum, const char* input, char* output)
{
    int insize = 32;

    uint8_t in[32];
    uint8_t in_step[32];
    uint8_t out[32];
    memset(in,  0, 32);
    memset(in_step,  0, 32);
    memset(out, 0, 32);
    memcpy((void*)in,      input, 32); // in
    memcpy((void*)in_step, input, 32); // in_step
    // in[0] = 1; // 测试正确性

    int k;
    for(k=0; k<2*loopnum; k++){
        x16rs_hash_sz((char*)in_step, (char*)in_step, insize); // 执行哈希算法 第一步
    }
    int n;
    for(n=0; n<loopnum; n++){
        memset(out, 0, 32);
        x16rs_hash_sz((char*)in, (char*)out, insize); // 执行哈希算法 第二步
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
    //    printf("%s\n", input);
    // print_byte_list("x16rs_hash_sz input", (void*)input, 32, 0);

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

    uint8_t in[32];
    memset(in, 0, 32); // 初始化
    memcpy((void*)in,   input, 32); // in
    memcpy((void*)hash, input, 32); // first

    int size = insize;

    int i;
    for(i = 0; i < 1; i++) {

        uint8_t algo = hash[7] % 16;

        // print_byte_list("x16rs_hash_sz hash", (void*)hash, 4, 0);
        // print_byte_list("x16rs_hash_sz hash[7]", (void*)&hash[7], 4, 0);
        // print_byte_list("x16rs_hash_sz algo", (void*)&algo, 1, 0);

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

        // in = (void*) hash;
        memset(in, 0, 32);
        memcpy((void*)in, hash, 32);
    }

    // memcpy(output, hash, 32);

    // print_byte_list("x16rs_hash_sz output", (void*)output, 32, 0);

    uint8_t results[32];
    memset(results, 0, 32); // 初始化
    memcpy((void*)results,  hash, 32);
    memcpy(output, results, 32);
}


// input length must more than 32
static const size_t x16rs_hash_insize = 32;
void x16rs_hash__development(const int loopnum, const char* input_hash, char* output_hash)
{
    // uint32_t input[64/4];
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

    // memcpy((void*)inputoutput, input_hash, 32); // first

    // uint32_t hash[64/4];
    // uint32_t *hash = (uint32_t*) input;
    // uint8_t algo = hash[7] % 16;
    // uint32_t algo = ( // 大端模式
    //     ucharn[28]*256*256*256+
    //     ucharn[29]*256*256 +
    //     ucharn[30]*256 +
    //     ucharn[31] ) % 16;


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
    for(n=0; n<loopnum; n++){

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


    // memcpy((void*)output_hash, output, 32); // first

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
    for(i=0; i<16; i++){
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
    int loopnum = diamondnumber / 8192 + 1; // 每 8192 颗钻石（约140天小半年）调整一下哈希次数
    if( loopnum > 16 ){
        loopnum = 16; // 最多 16 次
    }

    // 停止标记
    uint8_t *is_stop = (uint8_t*)stop_mark1;

    uint32_t basestuff[8+2+6+8];
    int basestufftargetsize = 61;
    memcpy( (void*)basestuff, (void*)input32, 32);
    memcpy( (void*)basestuff+40, (void*)addr21, 21);

    // 2万以上钻石需要附加 extend msg
    if(diamondnumber > 20000){
        memcpy( (void*)basestuff+40+21, (void*)extmsg32, 32);
        basestufftargetsize = 61 + 32;
    }

    // 开始计算
    uint8_t sha3res[32];
    uint8_t hashnew[32];
    uint8_t diamond[16];

    uint32_t noncenum1;
    uint32_t noncenum2;
    for(noncenum1 = hsstart; noncenum1 < hsend; noncenum1++){
        basestuff[8] = noncenum1;
        for(noncenum2=1; noncenum2<4294967294; noncenum2++){
            basestuff[9] = noncenum2;

            // 停止
            if( noncenum2%1000==0 && is_stop[0] != 0 ) {

                uint8_t noncenum_empty[8] = {0,0,0,0,0,0,0,0};
                memcpy( nonce8, noncenum_empty, 8);

                //////// TEST START ////////
                // uint32_t nonce[2] = {0, 0};
                // nonce[0] = noncenum1;
                // nonce[1] = noncenum2-1;
                // memcpy( (char*)nonce8, (char*)nonce, 8);
                // memcpy( (char*)diamond16, (char*)diamond, 16); 
                //////// TEST END //////// 

                return; // 返回空
            }
            // print_byte_list("1: ", (void*)basestuff, 61, 0);
            // 哈希计算
            sha3_256((char*)basestuff, basestufftargetsize, (char*)sha3res);
            // print_byte_list("2: ", (void*)sha3res, 32, 0);
            // x16rs_hash(loopnum, (char*)sha3res, (char*)hashnew);
            x16rs_hash__development(loopnum, (char*)sha3res, (char*)hashnew);
            // print_byte_list("3: ", (void*)hashnew, 32, 0);
            diamond_hash((char*)hashnew, (char*)diamond);
            // print_byte_list("4: ", (void*)diamond, 16, 0);
            /*
            printf("hash: ");
            uint8_t *input32p = (uint8_t*)hashnew;
            int j;
            for(j=0; j<32; j++){
                printf("%u,", input32p[j]);
            }
            printf("\n");
            */
            // 判断结果是否为钻石
            uint8_t success = 1;
            int k;
            for( k=0; k<16; k++ ) {
                if(k<10){
                    if( diamond[k] != 48 ){
                        success = 0;
                        break;
                    }
                }else{
                    if( diamond[k] == 48 ){
                        success = 0;
                        break;
                    }
                }
            }

            // 挖出钻石
            if( success == 1 ) {
                // 检查难度
                // 每 3277 颗钻石调整一下难度 3277 = 16^6 / 256 / 20
                // 难度最高时hash前20位为0，而不是32位都为0。
                uint8_t diffok = 0;
                int diffnum = diamondnumber / 3277; // 每 3277 颗钻石调整一下难度
                int n;
                for( n=0; n<32; n++ ) {
                    int bt = hashnew[n];
                    if(diffnum < 255){
                        if( bt + diffnum > 255) {
                            diffok = 0;
                            break; // 难度检查失败
                        } else {
                            diffok = 1;
                            break; // success
                        }
                    } else if(diffnum >= 255) {
                        if(bt != 0) {
                            diffok = 0;
                            break; // 难度检查失败
                        }
                        // 下一轮检查
                        diffnum -= 255;
                    }
                }
                // 难度满足要求
                if( diffok == 1 ) {
                    uint32_t nonce[2] = {0, 0};
                    nonce[0] = noncenum1;
                    nonce[1] = noncenum2;
                    memcpy( (char*)nonce8, (char*)nonce, 8);
                    memcpy( (char*)diamond16, (char*)diamond, 16);
                    // printf("\n%s\n", diamond); fflush(stdout);
                    return; // 拷贝值，返回成功
                }

                // 下一轮挖掘
            }


        }

    }

    // 循环完成，返回空，失败
    uint8_t noncenum_empty[8] = {0,0,0,0,0,0,0,0};
    memcpy( nonce8, noncenum_empty, 8);

}




// 挖矿算法
void miner_x16rs_hash(const int loopnum, const int retmaxhash, const char* stop_mark1, const uint32_t hsstart, const uint32_t hsend, const char* target_difficulty_hash32, const char* input_stuff89, char* stopkind1, char* success1, char* nonce4, char* reshash32)
{
//    printf("miner_x16rs_hash_v1()\n");
    // 签名信息
    uint8_t stuffnew_base[90];
    uint8_t *stuffnew = stuffnew_base+1;
    memcpy( stuffnew, input_stuff89, 89);
    uint32_t *stuffnew_uint32 = (uint32_t*)stuffnew_base;

    // 计算 sha3的结果
    unsigned char sha3res[32];

    // x16rs hash的结果
    uint8_t hashnew[32];
    uint8_t hashmaxpower[32];
    hashmaxpower[0] = 255; // 初始化
    uint8_t noncemaxpower[4];

    // 判断哈希是否满足要求
    uint8_t iscalcok = 0;
    uint8_t ispowerok = 0;

    // 检查结果的for循环值
    uint8_t pk = 0;
    uint8_t pi = 0;

    // 停止标记
    uint8_t *is_stop = (uint8_t*)stop_mark1;

    // nonce值
    uint32_t noncenum;
    for(noncenum = hsstart; noncenum < hsend; noncenum++){
        // 停止标记检测
        if( noncenum%5000==0 && is_stop[0] != 0 )
        {
            success1[0] = 0; // 失败
            stopkind1[0] = 1; // 外部信号强制停止
            if(retmaxhash == 1){
                memcpy(nonce4, &noncemaxpower, 4);
                memcpy(reshash32, &hashmaxpower, 32);
            }else{
                memcpy(nonce4, &noncenum, 4);
                memcpy(reshash32, &hashnew, 32);
            }
            return;
        }
        // 重置nonce
        stuffnew_uint32[20] = noncenum;
        // memcpy(&stuffnew[79], &noncenum, 4);
        // 计算 sha3
        sha3_256(stuffnew, 89, sha3res);
        /*
            printf("  hash: ");
            uint8_t i;
            for(i=0; i<32; i++){
                printf("%u,", sha3res[i]);
            }
            printf("\n");
        */
        // x16rs_hash(loopnum, ((void*)sha3res), ((void*)hashnew));
        x16rs_hash__development(loopnum, ((void*)sha3res), ((void*)hashnew));
        /*
            printf("  hash: ");
            uint8_t i;
            for(i=0; i<32; i++){
                printf("%u,", hashnew[i]);
            }
            printf("\n");
        */

        if(retmaxhash == 1){
            ispowerok = 0;
            for(pk=0; pk<32; pk++){
                uint8_t o1 = hashmaxpower[pk];
                uint8_t o2 = hashnew[pk];
                if(o2>o1){
                    break;
                }
                if(o2<o1){
                    ispowerok = 1;
                    break;
                }
            }
            if(ispowerok==1){
                memcpy(&noncemaxpower, &noncenum, 4);
                memcpy(&hashmaxpower, &hashnew, 32);
            }else{
                continue;
            }
        }

        iscalcok = 0;
        for(pi=0; pi<32; pi++){
            uint8_t o1 = target_difficulty_hash32[pi];
            uint8_t o2 = hashnew[pi];
            if(o2>o1){
                break;
            }
            if(o2<o1){
                iscalcok = 1;
                break;
            }
        }
        // 结果判断
        if(iscalcok == 1) {
        /*
            printf("finish hash: ");
            uint8_t i;
            for(i=0; i<32; i++){
                printf("%u,", hashnew[i]);
            }
            printf("\n");
        */
            // 返回 copy to out
            success1[0] = 1; // 成功的标记
            stopkind1[0] = 2; // 非强制停止，挖出成功后停止
            if(retmaxhash == 1){
                memcpy(nonce4, &noncemaxpower, 4);
                memcpy(reshash32, &hashmaxpower, 32);
            }else{
                memcpy(nonce4, &noncenum, 4);
                memcpy(reshash32, &hashnew, 32);
            }
            return; // success
        }
        // 继续下一次计算
    }

    // 循环完毕，失败
    success1[0] = 0;
    stopkind1[0] = 0; // 非强制停止，自然循环完毕停止
    if(retmaxhash == 1){
        memcpy(nonce4, &noncemaxpower, 4);
        memcpy(reshash32, &hashmaxpower, 32);
    }else{
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
    uint8_t output[32];
    // memcpy((void*)hash, input, 32); // first

    
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