package kerr

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestNewKError(t *testing.T) {
	err := New(1000, "test error")

	if err.Code() != 1000 {
		t.Errorf("expected code 1000, got %d", err.Code())
	}

	if err.Error() != "test error" {
		t.Errorf("expected msg 'test error', got '%s'", err.Error())
	}

	if err.Display() != "test error" {
		t.Errorf("expected display default to msg, got '%s'", err.Display())
	}
}

func TestWithDisplay(t *testing.T) {
	err := New(1001, "internal error")
	newErr := err.WithDisplay("对外提示信息")

	if newErr.Display() != "对外提示信息" {
		t.Errorf("expected display '对外提示信息', got '%s'", newErr.Display())
	}

	if newErr.Code() != 1001 {
		t.Errorf("code should be unchanged, got %d", newErr.Code())
	}
}

func TestWrapUnwrap(t *testing.T) {
	baseErr := errors.New("base error")
	err := New(2000, "kerr error").Wrap(baseErr)

	if !errors.Is(err, baseErr) {
		t.Errorf("expected errors.Is to detect base error")
	}

	unwrapped := err.Unwrap()
	if unwrapped != baseErr {
		t.Errorf("expected Unwrap() to return baseErr, got %v", unwrapped)
	}
}

func TestIs(t *testing.T) {
	err1 := New(3000, "err1")
	err2 := New(3000, "err2")

	if !err1.Is(err2) {
		t.Errorf("expected err1.Is(err2) to be true")
	}

	err3 := New(3001, "err3")
	if err1.Is(err3) {
		t.Errorf("expected err1.Is(err3) to be false")
	}
}

func TestWithStack(t *testing.T) {
	err := New(4000, "stack error")
	errWithStack := err.WithStack()

	stackStr := errWithStack.Stack()
	if stackStr == "" || !strings.Contains(stackStr, "TestWithStack") {
		t.Errorf("stack trace not captured correctly: %s", stackStr)
	}
}

func TestJSONSerialization(t *testing.T) {
	baseErr := errors.New("base error")
	err := New(5000, "json error").Wrap(baseErr).WithStack()

	data, marshalErr := json.Marshal(err)
	if marshalErr != nil {
		t.Fatalf("MarshalJSON failed: %v", marshalErr)
	}

	var result map[string]interface{}
	if jsonErr := json.Unmarshal(data, &result); jsonErr != nil {
		t.Fatalf("Unmarshal failed: %v", jsonErr)
	}

	if result["code"] != float64(5000) {
		t.Errorf("expected code 5000, got %v", result["code"])
	}

	if result["message"] == "" {
		t.Errorf("expected message field to be non-empty")
	}

	jsonStr := err.ToJSON()
	if !strings.Contains(jsonStr, "json error") {
		t.Errorf("ToJSON output missing message: %s", jsonStr)
	}
}

func TestFormat(t *testing.T) {
	err := New(6000, "format test").WithStack()
	output := fmt.Sprintf("%+v", err)

	if !strings.Contains(output, "format test") || !strings.Contains(output, "TestFormat") {
		t.Errorf("expected output to contain message and stack, got: %s", output)
	}

}

func TestNestedWrap(t *testing.T) {
	err1 := New(7000, "level 1").WithStack()
	err2 := err1.Wrap(New(7001, "level 2")).WithStack()
	err3 := err2.Wrap(errors.New("level 3")).WithStack()

	if !errors.Is(err3, err1) {
		t.Errorf("expected err3 to contain err1")
	}

	stack := err3.Stack()
	if stack == "" {
		t.Errorf("stack should be captured for nested wrap")
	}
}

func TestChainWithDisplay(t *testing.T) {
	err := New(8000, "original")
	err2 := err.WithDisplay("对外提示").Wrap(errors.New("cause"))

	if err2.Display() != "对外提示" {
		t.Errorf("expected display '对外提示', got '%s'", err2.Display())
	}
	if !strings.Contains(err2.Error(), "cause") {
		t.Errorf("expected Error() to contain cause, got '%s'", err2.Error())
	}
}
