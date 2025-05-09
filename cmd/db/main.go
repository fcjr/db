package main

import (
	"cmp"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/fcjr/db/internal/server"
)

func main() {
	ctx := context.Background()

	if err := run(ctx, os.Args, os.Getenv, os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string, getenv func(string) string, stdin io.Reader, stdout, stderr io.Writer) error {
	srv, err := server.New()
	if err != nil {
		return err
	}

	addr := cmp.Or(getenv("LISTEN_ADDR"), ":4000")
	return srv.ListenAndServe(ctx, addr)
}
