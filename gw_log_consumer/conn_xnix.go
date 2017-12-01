// +build linux darwin

package main

import (
	"net"

	"github.com/morya/utils/log"
)

func (a *App) buildConn(address string) bool {
	c, err := net.Dial("unix", address)
	if err != nil {
		log.InfoErrorf(err, "dial unix socket failed")
		return false
	}

	a.sock = c
	return true
}
