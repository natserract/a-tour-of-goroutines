
.PHONY: run_debug
run_debug:
	DEBUG_ALL=true k6 run script.js > output.txt 2>&1

.PHONY: run_timed
run_timed:
	k6 run --duration 10s script.js

.PHONY: run
run:
	k6 run script.js

.PHONY: runLoadTest
runLoadTest:
	LOAD_TEST=true k6 run script.js
