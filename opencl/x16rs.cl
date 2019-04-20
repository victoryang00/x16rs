#ifndef X16RX_CL
#define X16RX_CL

#define DEBUG(x)

#if __ENDIAN_LITTLE__
  #define SPH_LITTLE_ENDIAN 1
#else
  #define SPH_BIG_ENDIAN 1
#endif

#define SPH_UPTR sph_u64

typedef unsigned int sph_u32;
typedef int sph_s32;
#ifndef __OPENCL_VERSION__
  typedef unsigned long long sph_u64;
  typedef long long sph_s64;
#else
  typedef unsigned long sph_u64;
  typedef long sph_s64;
#endif

#define SPH_64 1
#define SPH_64_TRUE 1

#define SPH_C32(x)    ((sph_u32)(x ## U))
#define SPH_T32(x) (as_uint(x))
#define SPH_ROTL32(x, n) rotate(as_uint(x), as_uint(n))
#define SPH_ROTR32(x, n)   SPH_ROTL32(x, (32 - (n)))

#define SPH_C64(x)    ((sph_u64)(x ## UL))
#define SPH_T64(x) (as_ulong(x))
#define SPH_ROTL64(x, n) rotate(as_ulong(x), (n) & 0xFFFFFFFFFFFFFFFFUL)
#define SPH_ROTR64(x, n)   SPH_ROTL64(x, (64 - (n)))

#define SPH_ECHO_64 1
#define SPH_KECCAK_64 1
#define SPH_JH_64 1
#define SPH_SIMD_NOCOPY 0
#define SPH_KECCAK_NOCOPY 0
#define SPH_SMALL_FOOTPRINT_GROESTL 0
#define SPH_GROESTL_BIG_ENDIAN 0
#define SPH_CUBEHASH_UNROLL 0

#ifndef SPH_COMPACT_BLAKE_64
  #define SPH_COMPACT_BLAKE_64 0
#endif
#ifndef SPH_LUFFA_PARALLEL
  #define SPH_LUFFA_PARALLEL 0
#endif
#ifndef SPH_KECCAK_UNROLL
  #define SPH_KECCAK_UNROLL 0
#endif
#ifndef SPH_HAMSI_EXPAND_BIG
  #define SPH_HAMSI_EXPAND_BIG 1
#endif

#include "blake.cl"
// #include "bmw.cl"
// #include "groestl.cl"
// #include "jh.cl"
// #include "keccak.cl"
// #include "skein.cl"
// #include "luffa.cl"
// #include "cubehash.cl"
// #include "shavite.cl"
// #include "simd.cl"
// #include "echo.cl"
// #include "hamsi.cl"
// #include "fugue.cl"
// #include "shabal.cl"
// #include "whirlpool.cl"
// #include "sha2.cl"


#define SWAP4(x) as_uint(as_uchar4(x).wzyx)
#define SWAP8(x) as_ulong(as_uchar8(x).s76543210)

#if SPH_BIG_ENDIAN
  #define DEC64E(x) (x)
  #define DEC64BE(x) (*(const __global sph_u64 *) (x));
  #define DEC32LE(x) SWAP4(*(const __global sph_u32 *) (x));
#else
  #define DEC64E(x) SWAP8(x)
  #define DEC64BE(x) SWAP8(*(const __global sph_u64 *) (x));
  #define DEC32LE(x) (*(const __global sph_u32 *) (x));
#endif

#define SHL(x, n) ((x) << (n))
#define SHR(x, n) ((x) >> (n))

#define CONST_EXP2  q[i+0] + SPH_ROTL64(q[i+1], 5)  + q[i+2] + SPH_ROTL64(q[i+3], 11) + \
                    q[i+4] + SPH_ROTL64(q[i+5], 27) + q[i+6] + SPH_ROTL64(q[i+7], 32) + \
                    q[i+8] + SPH_ROTL64(q[i+9], 37) + q[i+10] + SPH_ROTL64(q[i+11], 43) + \
                    q[i+12] + SPH_ROTL64(q[i+13], 53) + (SHR(q[i+14],1) ^ q[i+14]) + (SHR(q[i+15],2) ^ q[i+15])

typedef union {
  unsigned char h1[64];
  uint h4[16];
  ulong h8[8];
} hash_t;



void hash_x16rs_func_0(__global unsigned char* block, __global unsigned char* output)
{
    hash_t hashobj;
    hash_t* hash = &hashobj;

    // blake
    sph_u64 H0 = SPH_C64(0x6A09E667F3BCC908), H1 = SPH_C64(0xBB67AE8584CAA73B);
    sph_u64 H2 = SPH_C64(0x3C6EF372FE94F82B), H3 = SPH_C64(0xA54FF53A5F1D36F1);
    sph_u64 H4 = SPH_C64(0x510E527FADE682D1), H5 = SPH_C64(0x9B05688C2B3E6C1F);
    sph_u64 H6 = SPH_C64(0x1F83D9ABFB41BD6B), H7 = SPH_C64(0x5BE0CD19137E2179);
    sph_u64 S0 = 0, S1 = 0, S2 = 0, S3 = 0;
    // sph_u64 T0 = SPH_C64(0xFFFFFFFFFFFFFC00) + (80 << 3), T1 = 0xFFFFFFFFFFFFFFFF;
    sph_u64 T0 = SPH_C64(0x0000000000000100), T1 = 0;

    // if ((T0 = SPH_T64(T0 + 1024)) < 1024)
    // T1 = SPH_T64(T1 + 1);

    sph_u64 M0, M1, M2, M3, M4, M5, M6, M7;
    sph_u64 M8, M9, MA, MB, MC, MD, ME, MF;
    sph_u64 V0, V1, V2, V3, V4, V5, V6, V7;
    sph_u64 V8, V9, VA, VB, VC, VD, VE, VF;

    M0 = DEC64BE(block +  0);
    M1 = DEC64BE(block +  8);
    M2 = DEC64BE(block + 16);
    M3 = DEC64BE(block + 24);
    M4 = SPH_C64(0x8000000000000000);
    M5 = 0;
    M6 = 0;
    M7 = 0;
    M8 = 0;
    M9 = 0;
    MA = 0;
    MB = 0;
    MC = 0;
    MD = SPH_C64(0x0000000000000001);
    ME = 0;
    MF = SPH_C64(0x0000000000000100);


    ulong hhhh8[16];
    hhhh8[0] = M0;
    hhhh8[1] = M1;
    hhhh8[2] = M2;
    hhhh8[3] = M3;
    hhhh8[4] = M4;
    hhhh8[5] = M5;
    hhhh8[6] = M6;
    hhhh8[7] = M7;
    hhhh8[8] = M8;
    hhhh8[9] = M9;
    hhhh8[10] = MA;
    hhhh8[11] = MB;
    hhhh8[12] = MC;
    hhhh8[13] = MD;
    hhhh8[14] = ME;
    hhhh8[15] = MF;


    COMPRESS64;

    ulong sss8[8];
    sss8[0] = H0;
    sss8[1] = H1;
    sss8[2] = H2;
    sss8[3] = H3;
    sss8[4] = H4;
    sss8[5] = H5;
    sss8[6] = H6;
    sss8[7] = H7;


    hash->h8[0] = SWAP8(H0);
    hash->h8[1] = SWAP8(H1);
    hash->h8[2] = SWAP8(H2);
    hash->h8[3] = SWAP8(H3);
    hash->h8[4] = SWAP8(H4);
    hash->h8[5] = SWAP8(H5);
    hash->h8[6] = SWAP8(H6);
    hash->h8[7] = SWAP8(H7);

    int i;
    for(i=0; i<32; i++){
        output[i] = hash->h1[i];
    }
    // unsigned char *ttt = (unsigned char*)(&T0);
    // for(i=0; i<8; i++){
    //     output[i] = ttt[i];
    // }


    // barrier(CLK_GLOBAL_MEM_FENCE);

}



































































#endif // X16RX_CL