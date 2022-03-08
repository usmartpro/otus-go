package main

import (
	"errors"
	"io"
	"net"
	"time"
)

var ErrConnectionError = errors.New("connection error")

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type Client struct {
	conn    net.Conn
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &Client{address: address, timeout: timeout, in: in, out: out}
}

func (t *Client) Connect() error {
	conn, err := net.DialTimeout("tcp", t.address, t.timeout)
	if err != nil {
		return ErrConnectionError
	}
	t.conn = conn
	return nil
}

func (t *Client) Send() error {
	if t.conn == nil {
		return ErrConnectionError
	}
	_, err := io.Copy(t.conn, t.in)
	return err
}

func (t *Client) Receive() error {
	if t.conn == nil {
		return ErrConnectionError
	}
	_, err := io.Copy(t.out, t.conn)
	return err
}

func (t *Client) Close() error {
	if t.conn != nil {
		return t.conn.Close()
	}

	return nil
}
