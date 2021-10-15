# tests

_module_ A basic testing framework. Compatible to TAP13.
It doesn't provide any expectations or matchers.
Instead there is only a fail function.

```
import tests

tests.test "my test case", { fail =>
when False, fail "when should not trigger when False"
}

when tests.enabled, tests.runTests
```

By default, tests are enabled, when `$LITHIA_TESTS` is set.

### Discussion

In case you really need expectations, writing a wrapper around tests should be possible.

## TestCase

_data_ 

### Properties

- `title` - The title of the test case for the logs.
- `impl fail` - The implementation of the test case.
Calls fail with a String, when failing

## TestSummary

_data_ The prinatble summary of all tests.

### Properties

- `ok` - How many tests have been ok.
- `notOk` - How many tests have been not ok.
- `failedTests` - List of failed test numbers.


## runTestCase

_func_ `runTestCase summary, testCase`


## runTests

_func_ `runTests `

Runs all test cases, that have been buffered by now.
## test

_func_ `test case, function`

Adds a new test case to the queue and will be executed once `runTests` has been called.

```
import tests

tests.test "my test case", { fail =>
when False, fail "when should not trigger when False"
}
```

