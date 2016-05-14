package errored

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestErrorStringFormat(t *testing.T) {
	refStr := "error string"
	e := Errorf("%s", refStr)
	e.SetDebug(true)

	fileName := "errored_test.go"
	lineNum := 13 // line number where error was formed
	funcName := "github.com/contiv/errored.TestErrorStringFormat"

	expectedStr := fmt.Sprintf("%s [%s %s %d]", refStr, funcName, fileName, lineNum)

	if e.Error() != expectedStr {
		t.Fatalf("error string mismatch. Expected: %q, got %q", expectedStr,
			e.Error())
	}
}

func getError(msg string) *Error {
	return Errorf(msg)
}

func TestErrorStackTrace(t *testing.T) {
	msg := "an error"
	e := getError(msg)
	e.SetTrace(false)
	e.SetDebug(true)

	if e.desc != msg {
		t.Fatal("Description did not match provided")
	}

	fileName := "errored_test.go"
	lineNum := 29 // line number where error was formed
	funcName := "github.com/contiv/errored.getError"

	expectedStr := fmt.Sprintf("%s [%s %s %d]", msg, funcName, fileName, lineNum)

	if e.Error() != expectedStr {
		t.Fatalf("Error message yielded an incorrect result with trace unset: %s %s", e.Error(), expectedStr)
	}

	e.SetTrace(false)
	e.SetDebug(false)
	if e.Error() != "an error" {
		t.Fatalf("Error message did yielded stack trace with trace unset: %q", e.Error())
	}

	e.SetTrace(true)
	if e.Error() == "an error\n" {
		t.Fatalf("Error message did not yield stack trace with trace set: %v", e.Error())
	}

	lines := strings.Split(e.Error(), "\n")

	if len(lines) != 6 {
		t.Fatalf("Stack trace yielded incorrect count: %d", len(lines))
	}
}

func TestErrorCombined(t *testing.T) {
	e := getError("one")
	e2 := getError("two")
	newErr := e.Combine(e2)
	if newErr.Error() != "one: two" {
		t.Fatalf("Errors did not combine in description: %v", newErr.Error())
	}

	if !reflect.DeepEqual(e.stack[0], newErr.stack[0]) {
		t.Fatalf("First stack was not equivalent: %v %v", e.stack, newErr.stack[0])
	}

	if !reflect.DeepEqual(e2.stack[0], newErr.stack[1]) {
		t.Fatalf("Second stack was not equivalent: %v %v", e.stack, newErr.stack[0])
	}

	if !newErr.Contains(e) || !newErr.Contains(e) {
		t.Fatal("Could not find original errors in combined error")
	}

	AlwaysDebug = true
	AlwaysTrace = true
	defer func() {
		AlwaysTrace = false
		AlwaysDebug = false
	}()

	if !newErr.Contains(e) || !newErr.Contains(e) {
		t.Fatal("Could not find original errors in combined error")
	}

	AlwaysTrace = false
	AlwaysDebug = false

	err := errors.New("my error")
	newErr = e.Combine(err)
	if newErr.Error() != "one: my error" {
		t.Fatalf("Could not combine error type into *Error: %v %v", e, err)
	}

	newErr = e.Combine(nil)
	if !reflect.DeepEqual(e, newErr) {
		t.Fatalf("Nil handling was inappropriate during combine")
  }

	if !newErr.ContainsFunc(func(err error) bool {
		return err.Error() == "one"
	}) {
		t.Fatal("ContainsFunc did not yield a true value")
	}
}

func TestAlways(t *testing.T) {
	defer func() {
		AlwaysDebug = false
		AlwaysTrace = false
	}()

	e := getError("one")
	if e.Error() != "one" {
		t.Fatalf("one != one, or so the test says. %q", e.Error())
	}

	AlwaysDebug = true
	if !strings.Contains(e.Error(), "[") {
		t.Fatalf("Debug output was not provided in error: %q", e.Error())
	}

	AlwaysTrace = true
	if !strings.Contains(e.Error(), "\n") {
		t.Fatalf("Trace output was not provided in error: %q", e.Error())
	}
}

func TestCode(t *testing.T) {
	e := Errorf("error")
	e.Code = 100
	if e.Error() != "100 error" {
		t.Fatalf("Error code output did not equal expectation: %v", e.Error())
	}
}
