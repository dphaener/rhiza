package main

import (
	"fmt"
	"io"

	"github.com/dphaener/rhiza/internal/api"
	"github.com/dphaener/rhiza/internal/canon"
)

// writeResponse prints a successful response: canonical JSON with --json,
// else its markdown from Human.
func writeResponse(stdout io.Writer, resp any, asJSON bool) error {
	if h, ok := resp.(api.Humaner); ok && !asJSON {
		h.Human(stdout)
		return nil
	}
	b, err := canon.Marshal(resp)
	if err != nil {
		return fmt.Errorf("encode response: %w", err)
	}
	_, err = stdout.Write(b)
	return err
}

// writeError prints err and returns its exit code. With --json, stdout gets
// the canonical {"error":{...}} and nothing else; without, stderr gets
// "rhiza: <message>" and, if there is one, "next: <command>".
func writeError(stdout, stderr io.Writer, err error, asJSON bool) int {
	if asJSON {
		b, merr := canon.Marshal(api.ErrorJSON(err))
		if merr != nil {
			fmt.Fprintf(stderr, "rhiza: %v\n", err)
			return 1
		}
		stdout.Write(b)
		return api.ExitCode(err)
	}
	e := api.AsError(err)
	fmt.Fprintf(stderr, "rhiza: %s\n", e.Message)
	if e.Next != "" {
		fmt.Fprintf(stderr, "next: %s\n", e.Next)
	}
	return api.ExitCode(err)
}
