#ifndef X16R_H
#define X16R_H

#ifdef __cplusplus
extern "C" {
#endif

#include <stdint.h>

void sha3_256(const char *input, const int in_size, char *output);

void x16r_hash(const char* input, char* output);
void x16rs_hash(const char* input, char* output);
void diamond_hash(const char* blkhash32, const char* addr21, const char* nonce8, char* output16);

void x16rs_hash_sz(const char* input, char* output, int insize);
void miner_diamond_hash(const char* input32, const char* addr21, char* nonce8, char* output16);

void miner_x16rs_hash_v1(const char* stop_mark1, const char* target_difficulty_hash32, const char* input_stuff88, char* nonce4);

#ifdef __cplusplus
}
#endif

#endif
