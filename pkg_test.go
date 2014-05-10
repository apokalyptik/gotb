package tb

import (
	"testing"
	"time"
)

func TestNoTokens(t *testing.T) {
	tb := new(TokenBucket)
	if tb.Get() {
		t.Errorf("Expected Get to return false, returned true")
	}
	if tb.Peek() {
		t.Errorf("Expected Peek to return false, returned true")
	}
}

func TestSomeTokens(t *testing.T) {
	tb := new(TokenBucket)
	tb.Set(1)
	if false == tb.Peek() {
		t.Errorf("Expected Peek with 1 token to return true, returned false")
	}
	if false == tb.Get() {
		t.Errorf("Expected Get with 1 token to return true, returned false")
	}
	if tb.Get() {
		t.Errorf("Expected Get to return false, returned true")
	}
	if tb.Peek() {
		t.Errorf("Expected Peek to return false, returned true")
	}
}

func TestAddingTokens(t *testing.T) {
	tb := new(TokenBucket)
	tb.Max(5)
	if tb.cur != 0 {
		t.Errorf("Tokens before starting should be 0, was %d", tb.cur)
	}
	err := tb.Start(time.Nanosecond)
	if err != nil {
		t.Errorf("Expected nil error, got %s", err.Error())
	}
	err = tb.Start(time.Nanosecond)
	if err == nil {
		t.Errorf("Expected an error, got nil")
	}
	time.Sleep(time.Millisecond)
	err = tb.Stop()
	if err != nil {
		t.Errorf("Expected nil error, got %s", err.Error())
	}
	err = tb.Stop()
	if err == nil {
		t.Errorf("Expected an error, got nil")
	}
	if tb.cur != 5 {
		t.Errorf("Expected 5 tokens after 10ms at a rate of 1 token/ms, got %d tokens", tb.cur)
	}
	if false == tb.Peek() {
		t.Errorf("Expected Peek with 5 tokens to return true, returned false")
	}
	for i := 0; i < 5; i++ {
		if false == tb.Get() {
			t.Errorf("Expected %d successfull Get()s", i)
		}
	}
	if tb.Get() {
		t.Errorf("Expected Get to return false, returned true")
	}
	if tb.Peek() {
		t.Errorf("Expected Peek to return false, returned true")
	}
	tb.Start(time.Nanosecond)
	time.Sleep(time.Millisecond)
	if false == tb.Peek() {
		t.Errorf("Expected Peek to return true, returned false")
	}

}
