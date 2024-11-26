package rpc

import "testing"

type EncodingExample struct {
	Testing bool `json:"testing"`
}

func TestEncode(t *testing.T) {
	expected := "Content-Length: 16\r\n\r\n{\"testing\":true}"
	actual, err := EncodeMessage(EncodingExample{Testing: true})
	if err != nil {
		t.Errorf("error, did not expect error but got. Error: %v", err)
	}
	if expected != string(actual) {
		t.Errorf("error, expected: %s does not equal actual: %s", expected, actual)
	}
}

func TestDecode(t *testing.T) {
	incommingMessage := "Content-Length: 15\r\n\r\n{\"Method\":\"hi\"}"
	method, content, err := DecodeMessage([]byte(incommingMessage))
	if err != nil {
		t.Errorf("error, did not expected error but got. Error %v", err)
	}

	if len(content) != 15 {
		t.Errorf("error, incorrect content length, got: %d, but expected: %d", len(content), 15)
	}

	if method != "hi" {
		t.Errorf("error, expected: %s does not equal actual: %s", "hi", method)
	}
}
