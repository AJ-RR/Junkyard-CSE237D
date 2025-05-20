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
		LDFLAGS += -lOpenCL
	else
		LDFLAGS += -L/usr/local/cuda/lib64 -lOpenCL
		INCFLAGS += -I/usr/local/cuda/include
	endif
else
	LDFLAGS += -lOpenCL
endif

all: solution

solution: ../helper_lib/helper_lib.a main.c
	$(CC) $(CFLAGS) -o $@ $^ $(INCFLAGS) $(LDFLAGS)

run: solution
	mkdir -p run_logs
	rm -f run_logs/summary.log results.json

	@start_time=$$(date +%s); \
	echo "{" > results.json; \
	echo '  "tests": [' >> results.json; \
	passed=0; \
	for i in $$(seq 0 9); do \
		echo "Running Dataset $$i..." | tee -a run_logs/summary.log; \
		./solution Dataset/$$i/input0.raw Dataset/$$i/input1.raw Dataset/$$i/input2.raw Dataset/$$i/input3.raw \
		          Dataset/$$i/output.raw Dataset/$$i/program_1_output.raw Dataset/$$i/program_2_output.raw \
		> run_logs/output_$$i.log 2>&1; \
		if [ $$? -eq 0 ]; then \
			status="PASS"; \
			score=10; \
			passed=$$((passed+1)); \
		else \
			status="FAIL"; \
			score=0; \
		fi; \
		echo "Dataset $$i: $$status" | tee -a run_logs/summary.log; \
		echo "Reason: See run_logs/output_$$i.log" >> run_logs/summary.log; \
		echo '    {' >> results.json; \
		echo '      "score": '"$$score"',' >> results.json; \
		echo '      "max_score": 10.0,' >> results.json; \
		echo '      "status": "'"$$status"'",' >> results.json; \
		echo '      "name": "Test '"$$((i+1))"'",' >> results.json; \
		echo '      "name_format": "text",' >> results.json; \
		echo '      "visibility": "visible"' >> results.json; \
		echo -n '    }' >> results.json; \
		if [ $$i -lt 9 ]; then echo "," >> results.json; else echo "" >> results.json; fi; \
	done; \
	end_time=$$(date +%s); \
	exec_time=$$((end_time - start_time)); \
	echo '  ],' >> results.json; \
	echo '  "score": '"$$((passed * 10))"',' >> results.json; \
	echo '  "n_time": '"$$exec_time" >> results.json; \
	echo "}" >> results.json; \
	echo "==== Run Summary ===="; \
	cat run_logs/summary.log

clean:
	rm -f parallel sequential solution
	cd ../helper_lib; make clean
	rm -rf run_logs
	rm -f results.json

grade: run
	@echo "==== Grading Summary ===="
	@passed=$$(grep -c "PASS" run_logs/summary.log); \
	total=10; \
	echo "Score: $$passed/$$total"; \
	if [ $$passed -eq $$total ]; then \
		echo "All datasets passed"; \
	else \
		echo "$$((total - passed)) dataset(s) failed"; \
	fi

score:
	@passed=$$(grep -c "PASS" run_logs/summary.log); \
	total=10; \
	echo "Score: $$passed/$$total"
