// Package release builds the channel manifest published beside a release's
// archives: the channel, version, commit and, per platform, the archive URL and
// sha256, so a client can find and verify the newest release.
package release

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ChannelRelease is the channel of an official v* release.
const ChannelRelease = "release"

// ChecksumsFile is the checksums file GoReleaser writes into dist.
const ChecksumsFile = "checksums.txt"

// ManifestFile is the manifest's file name, and the release asset's name.
const ManifestFile = "manifest.json"

// Platforms are the os/arch pairs every release must carry.
var Platforms = []Platform{
	{OS: "linux", Arch: "amd64"},
	{OS: "linux", Arch: "arm64"},
	{OS: "darwin", Arch: "amd64"},
	{OS: "darwin", Arch: "arm64"},
}

// Platform is one archive of a release.
type Platform struct {
	OS     string `json:"os"`
	Arch   string `json:"arch"`
	URL    string `json:"url,omitempty"`
	SHA256 string `json:"sha256,omitempty"`
}

// Manifest names one release and where to download and verify it.
type Manifest struct {
	Channel   string     `json:"channel"`
	Version   string     `json:"version"`
	Commit    string     `json:"commit"`
	Platforms []Platform `json:"platforms"`
}

// Options say what to build the manifest from.
type Options struct {
	Dist    string // folder with the archives and checksums.txt
	Project string // archive name prefix, e.g. rhiza
	Tag     string // git tag, e.g. v0.1.0; it is the manifest's version
	Version string // archive version, the tag without its v; defaults to it
	Commit  string // full commit sha the release was built from
	Repo    string // owner/name on GitHub
	Channel string // defaults to ChannelRelease
}

// Build reads dist and returns the manifest, or an error unless all four
// platforms are present and every archive's sha256 matches checksums.txt.
func Build(o Options) (*Manifest, error) {
	if o.Tag == "" || o.Commit == "" || o.Repo == "" || o.Project == "" || o.Dist == "" {
		return nil, fmt.Errorf("release: dist, project, tag, commit and repo are all required")
	}
	if o.Channel == "" {
		o.Channel = ChannelRelease
	}
	if o.Version == "" {
		o.Version = strings.TrimPrefix(o.Tag, "v")
	}
	f, err := os.Open(filepath.Join(o.Dist, ChecksumsFile))
	if err != nil {
		return nil, fmt.Errorf("release: %w", err)
	}
	defer f.Close()
	sums, err := ParseChecksums(f)
	if err != nil {
		return nil, err
	}
	m := &Manifest{Channel: o.Channel, Version: o.Tag, Commit: o.Commit}
	var missing []string
	for _, p := range Platforms {
		name := fmt.Sprintf("%s_%s_%s_%s.tar.gz", o.Project, o.Version, p.OS, p.Arch)
		want, ok := sums[name]
		if !ok {
			missing = append(missing, p.OS+"/"+p.Arch)
			continue
		}
		got, err := fileSHA256(filepath.Join(o.Dist, name))
		if err != nil {
			return nil, fmt.Errorf("release: %w", err)
		}
		if got != want {
			return nil, fmt.Errorf("release: %s sha256 is %s but %s says %s", name, got, ChecksumsFile, want)
		}
		p.SHA256 = want
		p.URL = fmt.Sprintf("https://github.com/%s/releases/download/%s/%s", o.Repo, o.Tag, name)
		m.Platforms = append(m.Platforms, p)
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("release: missing platforms: %s", strings.Join(missing, ", "))
	}
	return m, nil
}

// ParseChecksums reads a sha256sum-style file into file name to lowercase hex.
func ParseChecksums(r io.Reader) (map[string]string, error) {
	out := map[string]string{}
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 2 || len(fields[0]) != 64 {
			return nil, fmt.Errorf("release: bad checksums line %q", line)
		}
		out[strings.TrimPrefix(fields[1], "*")] = strings.ToLower(fields[0])
	}
	return out, sc.Err()
}

// Marshal returns the manifest's JSON, indented, ending in a newline.
func (m *Manifest) Marshal() ([]byte, error) {
	sort.SliceStable(m.Platforms, func(i, j int) bool {
		a, b := m.Platforms[i], m.Platforms[j]
		return a.OS+"/"+a.Arch < b.OS+"/"+b.Arch
	})
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(b, '\n'), nil
}

func fileSHA256(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
