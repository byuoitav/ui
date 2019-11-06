package bff

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Message map[string]json.RawMessage

func ErrorMessage(format string, a ...interface{}) Message {
	return StringMessage("error", format, a...)
}

func StringMessage(key string, format string, a ...interface{}) Message {
	m := make(map[string]json.RawMessage)
	m[key] = []byte(strconv.Quote(fmt.Sprintf(format, a...)))
	return m
}

func JSONMessage(key string, val interface{}) (Message, error) {
	data, err := json.Marshal(val)
	if err != nil {
		return nil, err
	}

	m := make(map[string]json.RawMessage)
	m[key] = data
	return m, nil
}
