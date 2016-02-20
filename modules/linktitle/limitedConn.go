package linktitle

import (
	"fmt"
	"io"
	"net"
	"time"
)

// This is a Conn that will error out and close the connection after a given byte limit
type limitedConn struct {
	c net.Conn
	r io.LimitedReader
}

// Return a limitedConn, to be plugged into places that expect a regular net.Conn
func NewLimitedConn(network string, address string, limit int64) (net.Conn, error) {
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}

	return &limitedConn{conn, io.LimitedReader{conn, limit}}, nil
}

// Do reading through the LimitedReader
func (c *limitedConn) Read(buf []byte) (int, error) {
	n, err := c.r.Read(buf)
	if err != nil {
		if c.r.N == 0 { // If the LimitedReader reached its limit
			fmt.Println("Connection limit exceeded.")
			c.Close()
		}
	}
	return n, err
}

// Pass other stuff to the regular Conn's functions

func (c *limitedConn) Close() error {
	return c.c.Close()
}

func (c *limitedConn) Write(buf []byte) (int, error) {
	return c.c.Write(buf)
}

func (c *limitedConn) LocalAddr() net.Addr {
	return c.c.LocalAddr()
}

func (c *limitedConn) RemoteAddr() net.Addr {
	return c.c.RemoteAddr()
}

func (c *limitedConn) SetDeadline(t time.Time) error {
	return c.c.SetDeadline(t)
}

func (c *limitedConn) SetWriteDeadline(t time.Time) error {
	return c.c.SetWriteDeadline(t)
}

func (c *limitedConn) SetReadDeadline(t time.Time) error {
	return c.c.SetReadDeadline(t)
}
