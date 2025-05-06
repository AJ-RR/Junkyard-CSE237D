CC       = gcc
CFLAGS   = -g -Wall
INCFLAGS := -I../helper_lib
LDFLAGS  := ../helper_lib/helper_lib.a -lm

UNAME_O := $(shell uname -o)
UNAME_M := $(shell uname -m)

ifeq ($(UNAME_O), Darwin)
	LDFLAGS += -framework OpenCL
else ifeq ($(UNAME_O), GNU/Linux)
	ifeq ($(UNAME_M), aarch64)
		# Assume ARM with OpenCL (like Raspberry Pi or Jetson)
		LDFLAGS += -lOpenCL
	else
		# Assume x86_64 with NVIDIA GPU
		LDFLAGS += -L/usr/local/cuda/lib64 -lOpenCL
		INCFLAGS += -I/usr/local/cuda/include
	endif
else # Android or unknown
	LDFLAGS += -lOpenCL
endif

all: solution

solution: ../helper_lib/helper_lib.a main.c
	$(CC) $(CFLAGS) -o $@ $^ $(INCFLAGS) $(LDFLAGS)

../helper_lib/helper_lib.a: 
	cd ../helper_lib; make

run: solution
	./solution Dataset/0/input0.raw Dataset/0/input1.raw Dataset/0/input2.raw Dataset/0/input3.raw Dataset/0/output.raw Dataset/0/program_1_output.raw Dataset/0/program_2_output.raw
	./solution Dataset/1/input0.raw Dataset/1/input1.raw Dataset/1/input2.raw Dataset/1/input3.raw Dataset/1/output.raw Dataset/1/program_1_output.raw Dataset/1/program_2_output.raw
	./solution Dataset/2/input0.raw Dataset/2/input1.raw Dataset/2/input2.raw Dataset/2/input3.raw Dataset/2/output.raw Dataset/2/program_1_output.raw Dataset/2/program_2_output.raw
	./solution Dataset/3/input0.raw Dataset/3/input1.raw Dataset/3/input2.raw Dataset/3/input3.raw Dataset/3/output.raw Dataset/3/program_1_output.raw Dataset/3/program_2_output.raw
	./solution Dataset/4/input0.raw Dataset/4/input1.raw Dataset/4/input2.raw Dataset/4/input3.raw Dataset/4/output.raw Dataset/4/program_1_output.raw Dataset/4/program_2_output.raw
	./solution Dataset/5/input0.raw Dataset/5/input1.raw Dataset/5/input2.raw Dataset/5/input3.raw Dataset/5/output.raw Dataset/5/program_1_output.raw Dataset/5/program_2_output.raw
	./solution Dataset/6/input0.raw Dataset/6/input1.raw Dataset/6/input2.raw Dataset/6/input3.raw Dataset/6/output.raw Dataset/6/program_1_output.raw Dataset/6/program_2_output.raw
	./solution Dataset/7/input0.raw Dataset/7/input1.raw Dataset/7/input2.raw Dataset/7/input3.raw Dataset/7/output.raw Dataset/7/program_1_output.raw Dataset/7/program_2_output.raw
	./solution Dataset/8/input0.raw Dataset/8/input1.raw Dataset/8/input2.raw Dataset/8/input3.raw Dataset/8/output.raw Dataset/8/program_1_output.raw Dataset/8/program_2_output.raw
	./solution Dataset/9/input0.raw Dataset/9/input1.raw Dataset/9/input2.raw Dataset/9/input3.raw Dataset/9/output.raw Dataset/9/program_1_output.raw Dataset/9/program_2_output.raw

clean: 
	rm -f parallel sequential solution
	cd ../helper_lib; make clean