// #ifndef X16RX_MAIN_CL
// #define X16RX_MAIN_CL



#include "x16rs.cl"



// x16rs hash 算法测试
__kernel void test_hash_x16rs(
   __global unsigned char* input,
   __global unsigned char* output)
{
    hash_x16rs_func_0(input, output);
    hash_x16rs_func_1(input, output);
    hash_x16rs_func_2(input, output);
    hash_x16rs_func_3(input, output);
    hash_x16rs_func_4(input, output);
    hash_x16rs_func_5(input, output);
    hash_x16rs_func_6(input, output);
    hash_x16rs_func_7(input, output);
    hash_x16rs_func_8(input, output);
    hash_x16rs_func_9(input, output);
    hash_x16rs_func_10(input, output);
    // hash_x16rs_func_11(output, output);
    // hash_x16rs_func_12(output, output);
    // hash_x16rs_func_13(output, output);
    // hash_x16rs_func_14(output, output);
    // hash_x16rs_func_15(output, output);
    
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