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
	mkdir -p run_logs
	rm -f run_logs/summary.log
	@for i in $$(seq 0 9); do \
		echo "Running Dataset $$i..." | tee -a run_logs/summary.log; \
		./solution Dataset/$$i/input0.raw Dataset/$$i/input1.raw Dataset/$$i/input2.raw Dataset/$$i/input3.raw \
		          Dataset/$$i/output.raw Dataset/$$i/program_1_output.raw Dataset/$$i/program_2_output.raw \
		> run_logs/output_$$i.log 2>&1; \
		if [ $$? -eq 0 ]; then \
			echo "Dataset $$i: PASS" | tee -a run_logs/summary.log; \
		else \
			echo "Dataset $$i: FAIL" | tee -a run_logs/summary.log; \
			echo "Reason: See run_logs/output_$$i.log" >> run_logs/summary.log; \
		fi; \
		echo "" >> run_logs/summary.log; \
	done
	@echo "==== Run Summary ===="
	@cat run_logs/summary.log

clean: 
	rm -f parallel sequential solution
	cd ../helper_lib; make clean
	rm -rf run_logs

grade: run
	@echo "==== Grading Summary ===="
	@passed=$$(grep -c "PASS" run_logs/summary.log); \
	total=10; \
	echo "Score: $$passed/$$total"; \
	if [ $$passed -eq $$total ]; then \
		echo "All datasets passed "; \
	else \
		echo "$$((total - passed)) dataset(s) failed "; \
	fi
