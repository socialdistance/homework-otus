package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/pflag"
)

var (
	ErrorHost      = errors.New("error host")
	ErrorPort      = errors.New("error port")
	ErrorConnected = errors.New("error connected")
)

func main() {
	var timeout time.Duration

	pflag.DurationVarP(&timeout, "timeout", "t", 10*time.Second, "")
	pflag.Parse()

	host := pflag.Arg(0)
	if host == "" {
		fmt.Println(ErrorHost)
		os.Exit(1)
	}

	port := pflag.Arg(1)
	if port == "" {
		fmt.Println(ErrorPort)
		os.Exit(1)
	}

	addr := net.JoinHostPort(host, port)
	clt := NewTelnetClient(addr, timeout, os.Stdin, os.Stdout)
	if err := clt.Connect(); err != nil {
		fmt.Println(ErrorConnected)
		os.Exit(1)
	}

	ctx, cnlFunc := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer clt.Close()

	go func() {
		defer cnlFunc()
		if err := clt.Send(); err != nil {
			fmt.Fprint(os.Stderr, err)
		}
	}()

	go func() {
		defer cnlFunc()
		if err := clt.Receive(); err != nil {
			fmt.Fprint(os.Stderr, err)
		}
	}()

	<-ctx.Done()
}
