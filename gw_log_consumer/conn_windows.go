package main

import (
	"net"
	"time"

	"github.com/morya/utils/log"
)

type FakeAddr struct{}

func (f *FakeAddr) Network() string {
	return ""
}

func (f *FakeAddr) String() string {
	return ""
}

type FakeConn struct {
}

func (c *FakeConn) Read(p []byte) (n int, err error) {
	return 0, nil
}

func (c *FakeConn) Write(p []byte) (n int, err error) {
	if len(p) > 1 {
		log.Infof("fakeconn %s", p)
	}
	return len(n), nil
}
func (c *FakeConn) Close() error {
	return nil
}

func (c *FakeConn) LocalAddr() net.Addr {
	return &FakeAddr{}
}

func (c *FakeConn) RemoteAddr() net.Addr {
	return &FakeAddr{}
}

func (c *FakeConn) SetDeadline(t time.Time) error {
	return nil
}
func (c *FakeConn) SetReadDeadline(t time.Time) error {
	return nil
}
func (c *FakeConn) SetWriteDeadline(t time.Time) error {
	return nil
}

func (a *App) buildConn(address string) bool {
	a.sock = &FakeConn{}
	return true
}
