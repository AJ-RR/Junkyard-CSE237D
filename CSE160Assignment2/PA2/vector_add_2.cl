__kernel void vectorAdd(__global const int *a, __global const int *b,
                        __global int *result, const unsigned int size) {
  //@@ Insert code to implement vector addition here
  int idx = get_global_id(0);

  if (idx < size) {
    result[idx] = a[idx] + b[idx];
  }
}
