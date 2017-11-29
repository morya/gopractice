package main

import "encoding/json"

type Counter struct {
	UnixTimeStamp int64  `json:"unixtimestamp"`
	Timestamp     string `json:"@timestamp"`

	SubmitRecv  int64 `json:"submitrecv"`
	SubmitAllow int64 `json:"submitallow"`
}

func fromLogStr(logline string) (*Counter, error) {
	c := &Counter{}

	err := json.Unmarshal([]byte(logline), c)
	if err != nil {
		return nil, err
	}

	return c, err
}
