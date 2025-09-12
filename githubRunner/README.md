#

## Setting Up ARC
Set up the Github Actions Controller Runner with

```
./setup.sh
```
The setup wizard will prompt with a couple of steps to ensure the setup is done properly

1. The URL for the repo or organization the runner will be assigned to. This will usually be in the form of `https://Github.com/organization` or `https://github.com/user/repo`
2. A personal access token (classic) for access to the org/repo. Ensure the token has scope `admin:org` or `repo` enabled
3. The minimum # of runners that will be running at once (Default: 1). Increasing this value will decrease the cold start time at the cost of computation
4. The maximum # of runners that can be running at once (Default: 5). Increasing this value will reduce peak latency at the cost of computation
5. Limits for CPU and Memory for each runner (TODO)

After the setup wizard has completed, the controller should show up with `kubectl get pods -n arc-systems` and any runners should show with `kubectl get autoscalingrunnersets -n arc-runners`. The runners should also show up in the org/repo settings under Actions > Runners

## Test Creation
By default, Github writes all the test cases to the `.github/classroom.yml` file which is included in all the student repos. Unfortunately, this method makes it very simple for a perceptive student to find the tests and code specifically for them.

To mitigate this problem, I built a program to build the test cases into a binary which Github can then run directly to produce the student's score.

The `testsCreator.go` file takes in a JSON file of test cases and produces a binary for each test. The script can be run with

`go run testsCreator.go <tests.json> <output prefix> <command> [args...] [--strict]`

Where `tests.json` includes the test cases in a structure like

```
[
    {
        "input": "inputString",
        "output": "outputString",
        "hidden": false
    },
    {
        "input": "hiddenTestInput",
        "output": "hiddenTestOutput",
        "hidden": true
    },
    ...
]

```

The program will then produce binaries `prefix_test1` `prefix_test2`, etc for all the tests, which can then be included in the starter repo and configured to run through the github classroom autograder.

**Note:** The binaries by default are set to cross-compile to an ARM64 Linux-compatible executable. Therefore, there may be issues if the autograding hardware is of a different OS/architecture.

**TODO:** Currently the script only handles a single input and a single output. Multiple sequential input is still to be implemented.

**DEMO:** Create tests with 

* `go run testsCreator.go tests.json python python demo.py`(Python)
* `go run testsCreator.go tests.json exec ./demo` (Golang)

then test the included programs with `tests/python_test1`or 
`tests/exec_test1`