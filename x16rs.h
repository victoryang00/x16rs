#ifndef X16R_H
#define X16R_H

#ifdef __cplusplus
extern "C" {
#endif

#include <stdint.h>

/* sha3_256 function to do hash jobs */
void sha3_256(const char *input, const int in_size, char *output);

/* x16r_hash to test x16r algorithm */
void x16r_hash(const char* input, char* output);

/* add loopnum for x16r algorithm, loop number between 1 - 16 and will add up to 16 and never changed */
void x16rs_hash(const int loopnum, const char* input, char* output);

/* convent hash32 string to output16, which is diamond string */
void diamond_hash(const char* hash32, char* output16);

void x16rs_hash_sz(const char* input, char* output, int insize);

/* to mint diamond */
void miner_diamond_hash(const uint32_t hsstart, const uint32_t hsend, const int diamondnumber,
						const char* stop_mark1, const char* input32, const char* addr21, const char* extmsg32,
						char* nonce8, char* diamond16);

/* to mint hac coin */
void miner_x16rs_hash(const int loopnum, const int retmaxhash, const char* stop_mark1,
					  const uint32_t hsstart, const uint32_t hsend, const char* target_difficulty_hash32,
					  const char* input_stuff88, char* stopkind1, char* success1, char* nonce4, char* reshash32);

/* test function to test x16rs algorithm */
void test_print_x16rs(const char* input , char* output32x16);

#ifdef __cplusplus
}
#endif

#endif
