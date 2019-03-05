#ifndef X16R_H
#define X16R_H

#ifdef __cplusplus
extern "C" {
#endif

#include <stdint.h>

void x16r_hash(const char* input, char* output);
void x16rs_hash(const char* input, char* output);
void diamond_hash(const char* blkhash32, const char* addr21, const char* nonce8, char* output16);

void x16rs_hash_sz(const char* input, char* output, int insize);
void miner_diamond_hash(const char* input32, const char* addr21, char* nonce8, char* output16);

#ifdef __cplusplus
}
#endif

#endif
