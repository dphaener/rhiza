package daemon

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/dphaener/rhiza/internal/release"
)

// maxBinaryBytes bounds the executable read out of a release archive.
const maxBinaryBytes = 512 << 20

// ChannelUpdate brings the daemon's own build up to date with a release
// channel: when the channel's manifest names a newer version than the running
// one it downloads the archive, checks its sha256 and replaces the running
// executable, whose change the daemon's re-exec then takes over. It runs no
// git and no make.
type ChannelUpdate struct {
	// Channel is release.ChannelLatest or release.ChannelRelease.
	Channel string
	// Resolver reads the channel's manifest and archive.
	Resolver *release.Resolver
	// Current is the running version.
	Current string
	// Exe is the executable to replace.
	Exe string
	// Timeout bounds one check and download; zero means a minute.
	Timeout time.Duration
	// Log reports what was done or skipped.
	Log func(format string, args ...any)
	// Now is the clock; nil checks on every Run.
	Now func() time.Time
	// Every is the least time between two checks.
	Every time.Duration

	checked time.Time
	// tried is the version an install was last attempted for, so a failed
	// one is not retried until the channel moves again.
	tried string
}

func (u *ChannelUpdate) logf(format string, args ...any) {
	if u.Log != nil {
		u.Log(format, args...)
	}
}

// Run does one check and reports whether the executable was replaced. A
// failing fetch, download, checksum or install is logged and skipped; a
// version that failed is not tried again until the channel moves. A Run within
// Every of the last check does nothing.
func (u *ChannelUpdate) Run() bool {
	if u == nil || u.Resolver == nil || u.Exe == "" {
		return false
	}
	if u.Now != nil {
		now := u.Now()
		if !u.checked.IsZero() && now.Sub(u.checked) < u.Every {
			return false
		}
		u.checked = now
	}
	timeout := u.Timeout
	if timeout == 0 {
		timeout = time.Minute
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	m, err := u.Resolver.Manifest(ctx, u.Channel)
	if err != nil {
		u.logf("auto-update %s: %v", u.Channel, err)
		return false
	}
	if !release.Newer(m.Version, u.Current) || m.Version == u.tried {
		return false
	}
	u.tried = m.Version
	archive, _, err := u.Resolver.Download(ctx, u.Channel)
	if err != nil {
		u.logf("auto-update %s: skipped until the channel moves again: %v", u.Channel, err)
		return false
	}
	if err := installBinary(u.Exe, archive); err != nil {
		u.logf("auto-update %s: skipped until the channel moves again: install: %v", u.Channel, err)
		return false
	}
	u.logf("auto-update %s: installed %s (was %s)", u.Channel, m.Version, u.Current)
	return true
}

// installBinary replaces exe with the rhiza executable of a release tar.gz.
// It writes beside exe and renames over it, so the running image is never
// half written.
func installBinary(exe string, archive []byte) error {
	bin, err := extractBinary(archive)
	if err != nil {
		return err
	}
	mode := os.FileMode(0o755)
	if fi, err := os.Stat(exe); err == nil {
		mode = fi.Mode().Perm() | 0o100
	}
	tmp, err := os.CreateTemp(filepath.Dir(exe), ".rhiza-update-*")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.Write(bin); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Chmod(mode); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), exe)
}

func extractBinary(archive []byte) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewReader(archive))
	if err != nil {
		return nil, fmt.Errorf("archive: %w", err)
	}
	tr := tar.NewReader(gz)
	for {
		h, err := tr.Next()
		if err == io.EOF {
			return nil, fmt.Errorf("archive has no rhiza executable")
		}
		if err != nil {
			return nil, fmt.Errorf("archive: %w", err)
		}
		if h.Typeflag != tar.TypeReg || path.Base(h.Name) != "rhiza" {
			continue
		}
		b, err := io.ReadAll(io.LimitReader(tr, maxBinaryBytes+1))
		if err != nil {
			return nil, fmt.Errorf("archive: %w", err)
		}
		if len(b) == 0 || len(b) > maxBinaryBytes {
			return nil, fmt.Errorf("archive: rhiza executable has a bad size")
		}
		return b, nil
	}
}
