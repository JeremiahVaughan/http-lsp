package rpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

type BaseMessage struct {
	Method string `json:"method"`
}

func EncodeMessage(msg any) ([]byte, error) {
	content, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("error, when encoding message to json for EncodeMessage(). Error: %v", err)
	}
	return []byte(fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(content), content)), nil
}

func DecodeMessage(msg []byte) (string, []byte, error) {
	header, content, found := bytes.Cut(msg, []byte{'\r', '\n', '\r', '\n'})
	if !found {
		return "", nil, errors.New("error, seperator not found for DecodeMessage()")
	}

	contentLengthBytes := header[len("Content-Length: "):]
	contentLength, err := strconv.Atoi(string(contentLengthBytes))
	if err != nil {
		return "", nil, fmt.Errorf("error, when parsing content length from message header. Error: %v", err)
	}

	var baseMessage BaseMessage
	actualContent := content[:contentLength]
	err = json.Unmarshal(actualContent, &baseMessage)
	if err != nil {
		return "", nil, fmt.Errorf("error, when decoding base message for DecodeMessage(). Error: %v", err)
	}

	return baseMessage.Method, actualContent, nil
}

func Split(data []byte, _ bool) (advance int, token []byte, err error) {
	header, content, found := bytes.Cut(data, []byte{'\r', '\n', '\r', '\n'})
	if !found {
		return 0, nil, nil
	}

	contentLengthBytes := header[len("Content-Length: "):]
	contentLength, err := strconv.Atoi(string(contentLengthBytes))
	if err != nil {
		return 0, nil, fmt.Errorf("error, when parsing content length from message header for Split(). Error: %v", err)
	}

	if len(content) != contentLength {
		return 0, nil, nil
	}

	totalLength := len(header) + 4 + contentLength
	return totalLength, data[:totalLength], nil
}
