package urlchecker

import (
	"testing"
	"time"
)

const (
	URL = "/LVE07/HEM/1949/06/25/LVG19490625-011.pdf"
)

func TestCheckTokenOk(t *testing.T) {

	timestamp := GetTimestamp()
	token := GenerateToken(URL, timestamp)
	err := Check(URL, timestamp, token)
	if err != nil {
		t.Fatalf("Error, this test should not have error but it gets %v\n", err)
	}
}

func TestCheckTokenOldTimestamp(t *testing.T) {
	timestamp := time.Now().Unix() - (max_seconds + 1000)
	token := GenerateToken(URL, timestamp)
	err := Check(URL, timestamp, token)

	if err == nil {
		t.Fatalf("Error, this test should have an error")
	}
	if err.Error() != url_too_old_error {
		t.Fatalf("Error should be %v and really is %v\n", url_too_old_error, err)
	}
}

func TestCheckTokenManipulatedTimestamp(t *testing.T) {
	timestamp := time.Now().Unix()
	token := GenerateToken(URL, timestamp)
	err := Check(URL, timestamp+1, token)
	if err == nil {
		t.Fatalf("Error, this test should have an error\n")
	}

	if err.Error() != token_mismatch_error {
		t.Fatalf("Error should be %v and really is %v\n", token_mismatch_error, err)
	}
}

func TestCheckTokenManipulatedUrlBack(t *testing.T) {
	timestamp := time.Now().Unix()
	token := GenerateToken(URL, timestamp)
	err := Check(URL+"aaa", timestamp, token)

	if err == nil {
		t.Fatalf("Error, this test should have an error\n")
	}

	if err.Error() != token_mismatch_error {
		t.Fatalf("Error should be %v and really is %v\n", token_mismatch_error, err)
	}
}

func TestCheckTokenManipulatedUrlFront(t *testing.T) {
	timestamp := time.Now().Unix()
	token := GenerateToken(URL, timestamp)
	err := Check("/aaa"+URL, timestamp, token)
	if err == nil {
		t.Fatalf("Error, this test should have an error\n")
	}

	if err.Error() != token_mismatch_error {
		t.Fatalf("Error should be %v and really is %v\n", token_mismatch_error, err)
	}
}
