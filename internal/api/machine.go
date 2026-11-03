package api

// This file holds machine.show: the machine's trust and the config file it
// came from. It reads the machine config and writes nothing; rhiza has no
// operation that writes trust, the file is edited by hand.

import (
	"fmt"
	"io"
	"os"
	"reflect"

	"github.com/dphaener/rhiza/internal/machine"
)

func init() {
	register(Operation{
		Name:     "machine.show",
		Summary:  "Show the machine's trust and the machine config file it came from, daily with no machine config; writes nothing",
		Request:  reflect.TypeFor[MachineShowRequest](),
		Response: reflect.TypeFor[MachineShowResponse](),
		Run:      runMachineShow,
	})
}

// MachineShowRequest takes nothing.
type MachineShowRequest struct{}

// MachineShowResponse is the machine's trust.
type MachineShowResponse struct {
	Trust string `json:"trust" enum:"daily,dedicated" desc:"The machine's trust: dedicated when its config says so, else daily"`
	Path  string `json:"path" desc:"The machine config file, whether or not it exists"`
	File  bool   `json:"file" desc:"Whether the machine config file exists; false means trust is daily by default"`
}

// Human writes the trust and the file.
func (r MachineShowResponse) Human(w io.Writer) {
	if !r.File {
		fmt.Fprintf(w, "trust: %s (no machine config; it would be %s)\n", r.Trust, r.Path)
		return
	}
	fmt.Fprintf(w, "trust: %s (from %s)\n", r.Trust, r.Path)
}

func runMachineShow(ctx *Ctx, _ any) (any, error) {
	if ctx.DryRun {
		return nil, Usagef("machine show does not take --dry-run: it writes nothing").WithNext("rhiza machine show")
	}
	path := machine.Path(ctxGetenv(ctx), ctx.Home)
	c, err := machine.Read(path)
	if err != nil {
		return nil, Errorf("%v", err).WithNext("rhiza machine show")
	}
	return MachineShowResponse{Trust: c.Trust, Path: path, File: machineFileExists(path)}, nil
}

func machineFileExists(path string) bool {
	if path == "" {
		return false
	}
	_, err := os.Stat(path)
	return err == nil
}
