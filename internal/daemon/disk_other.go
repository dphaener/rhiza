//go:build !unix

package daemon

// FreeBytes is not known off unix: the disk guard does not hold anything.
func FreeBytes(dir string) (uint64, error) {
	return ^uint64(0), nil
}
