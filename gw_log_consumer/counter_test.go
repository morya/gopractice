package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestCounter(t *testing.T) {
	var c = &Counter{
		UnixTimeStamp: 15,
		SubmitRecv:    100,
		SubmitAllow:   55,
	}
	data, _ := json.Marshal(c)
	fmt.Printf("%s\n", data)
}
