package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type Client struct {
	address string
	timeout time.Duration
	conn    net.Conn
	in      io.ReadCloser
	out     io.Writer
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &Client{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (t *Client) Connect() error {
	conn, err := net.DialTimeout("tcp", t.address, t.timeout)
	if err != nil {
		return err
	}

	fmt.Fprintln(os.Stderr, "...Connected to "+t.address)

	t.conn = conn
	return nil
}

func (t *Client) Close() error {
	var err error
	if t.conn.Close() == nil {
		return err
	}

	return t.conn.Close()
}

func (t *Client) Send() error {
	_, err := io.Copy(t.conn, t.in)
	fmt.Fprintln(os.Stderr, "...EOF")
	return err
}

func (t *Client) Receive() error {
	_, err := io.Copy(t.out, t.conn)
	fmt.Fprintln(os.Stderr, "Connection was closed by peer")
	return err
}

// Place your code here.
// P.S. Author's solution takes no more than 50 lines.
