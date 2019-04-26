// #ifndef X16RX_MAIN_CL
// #define X16RX_MAIN_CL


#include "sha3_256.cl"
#include "x16rs.cl"


void hash_x16rs_choice_step(hash_t* stephash){

    uint8_t algo = stephash->h4[7] % 16;

    switch (algo) {
        case 0 : hash_x16rs_func_0  ( stephash ); break;
        case 1 : hash_x16rs_func_1  ( stephash ); break;
        case 2 : hash_x16rs_func_2  ( stephash ); break;
        case 3 : hash_x16rs_func_3  ( stephash ); break;
        case 4 : hash_x16rs_func_4  ( stephash ); break;
        case 5 : hash_x16rs_func_5  ( stephash ); break;
        case 6 : hash_x16rs_func_6  ( stephash ); break;
        case 7 : hash_x16rs_func_7  ( stephash ); break;
        case 8 : hash_x16rs_func_8  ( stephash ); break;
        case 9 : hash_x16rs_func_9  ( stephash ); break;
        case 10: hash_x16rs_func_10 ( stephash ); break;
        case 11: hash_x16rs_func_11 ( stephash ); break;
        case 12: hash_x16rs_func_12 ( stephash ); break;
        case 13: hash_x16rs_func_13 ( stephash ); break;
        case 14: hash_x16rs_func_14 ( stephash ); break;
        case 15: hash_x16rs_func_15 ( stephash ); break;
    }
    

}



// x16rs hash 算法
__kernel void hash_x16rs(
   __global unsigned char* input,
   __global unsigned char* output)
{
    ////////////////////////////////////
    hash_t hhh;
    // hash_x16rs_func_0 (&hhh ); 
    // hash_x16rs_func_1 (&hhh ); 
    // hash_x16rs_func_2 (&hhh ); 
    // hash_x16rs_func_3 (&hhh ); 
    // hash_x16rs_func_4 (&hhh ); 
    // hash_x16rs_func_5 (&hhh ); 
    // hash_x16rs_func_6 (&hhh ); 
    // hash_x16rs_func_7 (&hhh ); 
    // hash_x16rs_func_8 (&hhh ); 
    // hash_x16rs_func_9 (&hhh ); 
    // hash_x16rs_func_10 (&hhh ); 
    // hash_x16rs_func_11 (&hhh ); 
    // hash_x16rs_func_12 (&hhh ); 
    // hash_x16rs_func_13 (&hhh ); 
    // hash_x16rs_func_14 (&hhh ); 
    // hash_x16rs_func_15 (&hhh ); 
    ////////////////////////////////////

    hash_t hsobj ;
    for(int i = 0; i < 32; i++){
        hsobj.h1[i] = input[i];
    }

    // 计算
    // for(int i=0; i<1; i++){
    hash_x16rs_choice_step(&hsobj);
    // }

    // 结果
    for(int i=0; i<32; i++){
        output[i] = hsobj.h1[i];
    }

}



// x16rs hash miner 算法
__kernel void miner_do_hash_x16rs(
   __global unsigned char* target_difficulty_hash_32,
   __global unsigned char* input_stuff_89,
   const unsigned int   base_start,
   __global unsigned char* output_nonce_4)
{

    int nonce = 23645 + get_global_id(0);
    char *nonce_ptr = &nonce;

    // stuff
    unsigned char base_stuff[89];
    for(int i=0; i<89; i++){
        base_stuff[i] = input_stuff_89[i];
    }
    base_stuff[79] = nonce_ptr[0];
    base_stuff[80] = nonce_ptr[1];
    base_stuff[81] = nonce_ptr[2];
    base_stuff[82] = nonce_ptr[3];

    // hash x16rs
    unsigned char hash1[64];
    sha3_256_hash(base_stuff, 89, hash1);

    hash_x16rs_choice_step(hash1);

    // miner check
    char is_ok = 1;
    for(int i=0; i<32; i++){
        unsigned char a1 = hash1[i];
        unsigned char a2 = target_difficulty_hash_32[i];
        if( a1 > a2 ){
            is_ok = 0;
            break;
        }else if( a1 < a2 ){
            is_ok = 1;
            break;
        }
    }

    // copy set
    if(1){
        output_nonce_4[0] = nonce_ptr[0];
        output_nonce_4[1] = nonce_ptr[1];
        output_nonce_4[2] = nonce_ptr[2];
        output_nonce_4[3] = nonce_ptr[3];
    }


}



///////////////////////////////////////////////////////////



// sha3 hash 算法
__kernel void hash_sha3(
   __global unsigned char* input,
   __global unsigned char* output)
{
    hash_t hs0;
    for(int i = 0; i < 32; i++)
        hs0.h1[i] = input[i];

    sha3_256_hash(hs0.h1, 32, hs0.h1);

    // 结果
    for(int i=0; i<32; i++){
        output[i] = hs0.h1[i];
    }


}


//////////////////////////////////////////////////////////////////////////





// x16rs hash 算法测试
__kernel void test_hash_x16rs(
   __global unsigned char* input,
   __global unsigned char* output)
{
    hash_t hs0 ;
    hash_t hs1 ;
    hash_t hs2 ;
    hash_t hs3 ;
    hash_t hs4 ;
    hash_t hs5 ;
    hash_t hs6 ;
    hash_t hs7 ;
    hash_t hs8 ;
    hash_t hs9 ;
    hash_t hs10;
    hash_t hs11;
    hash_t hs12;
    hash_t hs13;
    hash_t hs14;
    hash_t hs15;

    for(int i = 0; i < 32; i++)
        hs0.h1[i] = input[i];



    hash_t hhh;
    hash_x16rs_func_0 (&hhh ); 
    hash_x16rs_func_1 (&hhh ); 
    hash_x16rs_func_2 (&hhh ); 
    hash_x16rs_func_3 (&hhh ); 
    hash_x16rs_func_4 (&hhh ); 
    hash_x16rs_func_5 (&hhh ); 
    hash_x16rs_func_6 (&hhh ); 
    hash_x16rs_func_7 (&hhh ); 
    hash_x16rs_func_8 (&hhh ); 
    hash_x16rs_func_9 (&hhh ); 
    hash_x16rs_func_10 (&hhh ); 
    hash_x16rs_func_11 (&hhh ); 
    hash_x16rs_func_12 (&hhh ); 
    hash_x16rs_func_13 (&hhh ); 
    hash_x16rs_func_14 (&hhh ); 
    hash_x16rs_func_15 (&hhh ); 




    hash_x16rs_func_0 (&hs0 ); 
    for(int i = 0; i < 32; i++) hs1 .h1[i] = hs0 .h1[i];
    hash_x16rs_func_1 (&hs1 ); 
    for(int i = 0; i < 32; i++) hs2 .h1[i] = hs1 .h1[i];
    hash_x16rs_func_2 (&hs2 ); 
    for(int i = 0; i < 32; i++) hs3 .h1[i] = hs2 .h1[i];
    hash_x16rs_func_3 (&hs3 ); 
    for(int i = 0; i < 32; i++) hs4 .h1[i] = hs3 .h1[i];
    hash_x16rs_func_4 (&hs4 ); 
    for(int i = 0; i < 32; i++) hs5 .h1[i] = hs4 .h1[i];
    hash_x16rs_func_5 (&hs5 ); 
    for(int i = 0; i < 32; i++) hs6 .h1[i] = hs5 .h1[i];
    hash_x16rs_func_6 (&hs6 ); 
    for(int i = 0; i < 32; i++) hs7 .h1[i] = hs6 .h1[i];
    hash_x16rs_func_7 (&hs7 ); 
    for(int i = 0; i < 32; i++) hs8 .h1[i] = hs7 .h1[i];
    hash_x16rs_func_8 (&hs8 ); 
    for(int i = 0; i < 32; i++) hs9 .h1[i] = hs8 .h1[i];
    hash_x16rs_func_9 (&hs9 ); 
    for(int i = 0; i < 32; i++) hs10.h1[i] = hs9 .h1[i];
    hash_x16rs_func_10(&hs10); 
    for(int i = 0; i < 32; i++) hs11.h1[i] = hs10.h1[i];
    hash_x16rs_func_11(&hs11); 
    for(int i = 0; i < 32; i++) hs12.h1[i] = hs11.h1[i];
    hash_x16rs_func_12(&hs12); 
    for(int i = 0; i < 32; i++) hs13.h1[i] = hs12.h1[i];
    hash_x16rs_func_13(&hs13); 
    for(int i = 0; i < 32; i++) hs14.h1[i] = hs13.h1[i];
    hash_x16rs_func_14(&hs14); 
    for(int i = 0; i < 32; i++) hs15.h1[i] = hs14.h1[i];
    hash_x16rs_func_15(&hs15);

    // 结果
    for(int i=0; i<32; i++){
        output[i] = hs15.h1[i];
    }

}

/*
// x16rs hash 算法测试
__kernel void test_hash_x16rs_old(
   __global unsigned char* input,
   __global unsigned char* output)
{

    unsigned char iiipppttt[64];

    for(int i=0; i<32; i++){
        iiipppttt[i] = input[i];
    }

    unsigned char oootttppp[64];

    unsigned char innnn[64];
    unsigned char otttt[64];

    hash_x16rs_func_0 (iiipppttt, oootttppp);
    hash_x16rs_func_1 (oootttppp, oootttppp);
    hash_x16rs_func_2 (oootttppp, oootttppp);
    hash_x16rs_func_3 (oootttppp, oootttppp);
    hash_x16rs_func_4 (oootttppp, oootttppp);
    hash_x16rs_func_5 (oootttppp, oootttppp);
    // hash_x16rs_func_6 (oootttppp, oootttppp);
    // hash_x16rs_func_7 (oootttppp, oootttppp);
    // hash_x16rs_func_8 (oootttppp, oootttppp);
    // hash_x16rs_func_9 (oootttppp, oootttppp);
    // hash_x16rs_func_10(oootttppp, oootttppp);
    hash_x16rs_func_11(oootttppp, oootttppp);
    hash_x16rs_func_12(oootttppp, oootttppp);
    hash_x16rs_func_13(oootttppp, oootttppp);
    hash_x16rs_func_14(oootttppp, oootttppp);
    hash_x16rs_func_15(oootttppp, oootttppp);
    
    for(int i=0; i<32; i++){
        output[i] = oootttppp[i];
    }

}

*/



// 矩阵算法测试
__kernel void square(
   __global float* input,
   __global float* output,
   const unsigned int count)
{
   int i = get_global_id(0);
   if(i < count)
       output[i] = input[i] * input[i];
}


// #endif // X16RX_MAIN_CL