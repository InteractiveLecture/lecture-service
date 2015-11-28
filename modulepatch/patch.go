package modulepatch

import (
	"encoding/json"
	"io"
)

type OperationType string

const (
	ADD     OperationType = "add"
	REMOVE                = "remove"
	REPLACE               = "replace"
	COPY                  = "copy"
	MOVE                  = "move"
	TEST                  = "test"
)

type Operation struct {
	Type  OperationType `json:"op"`
	Path  string        `json:"path"`
	From  string        `json:"from"`
	Value interface{}   `json:"value"`
}

type Patch struct {
	Operations []Operation `json:"operations"`
	Version    int         `json:"version"`
	TopicID    string      `json:"topic_id"`
}

func Decode(reader io.Reader) (*Patch, error) {
	result := &Patch{}
	return result, json.NewDecoder(reader).Decode(result)
}
