// #ifndef X16RX_MAIN_CL
// #define X16RX_MAIN_CL



#include "x16rs.cl"



// x16rs hash 算法测试
__kernel void test_hash_x16rs(
   __global unsigned char* input,
   __global unsigned char* output)
{
    hash_x16rs_func_0(input, output);

    /*
    int i;
    for(i=0; i<32; i++){
        output[i] = input[i];
    }
    */
    
}


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