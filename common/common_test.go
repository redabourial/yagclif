package common

import (
	"testing"
)

type stub struct {
	called        bool
	multipleCalls bool
	testCall      func(format string, args ...interface{})
}

func (s *stub) Fatalf(format string, args ...interface{}) {
	if !s.called {
		s.called = true
	} else {
		s.multipleCalls = true
	}
	if s.testCall != nil {
		s.testCall(format, args)
	}
}

func TestAssertEqual(t *testing.T) {
	mock1 := stub{
		called:        false,
		multipleCalls: false,
		testCall: func(format string, args ...interface{}) {
			msg := " hello world\r\n\tvalue : 1\r\n\texpected : 0\r\n"
			if format != msg {
				t.Fatalf("format is not as expected")
			}
			if len(args) > 1 {
				t.Fatalf("using args")
			}
		},
	}
	AssertEqual(1, 0, &mock1, "hello", "world")
	if !mock1.called {
		t.Fatalf("testing called when it should not")
	}
	mock2 := stub{
		called:        false,
		multipleCalls: false,
	}
	AssertEqual(1, 1, &mock2, "hello", "world")
	if mock2.called {
		t.Fatalf("testing not called when it should")
	}
}
