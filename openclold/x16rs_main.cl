// #ifndef X16RX_MAIN_CL
// #define X16RX_MAIN_CL

// #include "x16rs.cl"
#include "test.cl"

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