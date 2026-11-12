package release

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strings"
)

// ChannelLatest is the channel of the newest build of main.
const ChannelLatest = "latest"

// Repo is the GitHub repository releases are published to.
const Repo = "dphaener/rhiza"

// maxManifestBytes and maxArchiveBytes bound what a resolver reads.
const (
	maxManifestBytes = 1 << 20
	maxArchiveBytes  = 512 << 20
)

// ChannelURLs are where each channel's manifest is published.
var ChannelURLs = map[string]string{
	ChannelRelease: "https://github.com/" + Repo + "/releases/latest/download/" + ManifestFile,
	ChannelLatest:  "https://github.com/" + Repo + "/releases/download/latest/" + ManifestFile,
}

// Resolver fetches channel manifests and the archives they name.
type Resolver struct {
	// URLs maps a channel to its manifest URL; nil means ChannelURLs.
	URLs map[string]string
	// Client does every request; nil means http.DefaultClient.
	Client *http.Client
	// OS and Arch pick the archive; empty means the running platform.
	OS, Arch string
}

func (r *Resolver) client() *http.Client {
	if r.Client != nil {
		return r.Client
	}
	return http.DefaultClient
}

func (r *Resolver) platform() (string, string) {
	os, arch := r.OS, r.Arch
	if os == "" {
		os = runtime.GOOS
	}
	if arch == "" {
		arch = runtime.GOARCH
	}
	return os, arch
}

func (r *Resolver) get(ctx context.Context, url string, limit int64) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("release: %w", err)
	}
	resp, err := r.client().Do(req)
	if err != nil {
		return nil, fmt.Errorf("release: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("release: GET %s: %s", url, resp.Status)
	}
	b, err := io.ReadAll(io.LimitReader(resp.Body, limit+1))
	if err != nil {
		return nil, fmt.Errorf("release: %w", err)
	}
	if int64(len(b)) > limit {
		return nil, fmt.Errorf("release: GET %s: larger than %d bytes", url, limit)
	}
	return b, nil
}

// Manifest fetches and validates the channel's manifest.
func (r *Resolver) Manifest(ctx context.Context, channel string) (*Manifest, error) {
	urls := r.URLs
	if urls == nil {
		urls = ChannelURLs
	}
	url, ok := urls[channel]
	if !ok {
		return nil, fmt.Errorf("release: unknown channel %q (want %s or %s)", channel, ChannelRelease, ChannelLatest)
	}
	b, err := r.get(ctx, url, maxManifestBytes)
	if err != nil {
		return nil, err
	}
	return ParseManifest(b, channel)
}

// ParseManifest decodes a manifest and checks it is complete and is the
// channel's own.
func ParseManifest(b []byte, channel string) (*Manifest, error) {
	var m Manifest
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, fmt.Errorf("release: malformed manifest: %w", err)
	}
	if m.Version == "" || m.Commit == "" {
		return nil, fmt.Errorf("release: malformed manifest: version and commit are required")
	}
	if m.Channel != channel {
		return nil, fmt.Errorf("release: manifest is for channel %q, not %q", m.Channel, channel)
	}
	for _, p := range m.Platforms {
		if p.OS == "" || p.Arch == "" || p.URL == "" || len(p.SHA256) != 64 {
			return nil, fmt.Errorf("release: malformed manifest: incomplete platform %s/%s", p.OS, p.Arch)
		}
	}
	return &m, nil
}

// Newest reports the newest version published on the channel.
func (r *Resolver) Newest(ctx context.Context, channel string) (string, error) {
	m, err := r.Manifest(ctx, channel)
	if err != nil {
		return "", err
	}
	return m.Version, nil
}

// Pick returns the manifest's entry for the resolver's platform.
func (r *Resolver) Pick(m *Manifest) (Platform, error) {
	os, arch := r.platform()
	for _, p := range m.Platforms {
		if p.OS == os && p.Arch == arch {
			return p, nil
		}
	}
	return Platform{}, fmt.Errorf("release: %s has no archive for %s/%s", m.Version, os, arch)
}

// Download fetches the channel's archive for this platform and returns it with
// its manifest, unless its sha256 differs from the manifest's.
func (r *Resolver) Download(ctx context.Context, channel string) ([]byte, *Manifest, error) {
	m, err := r.Manifest(ctx, channel)
	if err != nil {
		return nil, nil, err
	}
	p, err := r.Pick(m)
	if err != nil {
		return nil, nil, err
	}
	b, err := r.get(ctx, p.URL, maxArchiveBytes)
	if err != nil {
		return nil, nil, err
	}
	sum := sha256.Sum256(b)
	if got := hex.EncodeToString(sum[:]); got != strings.ToLower(p.SHA256) {
		return nil, nil, fmt.Errorf("release: %s sha256 is %s, manifest says %s", p.URL, got, p.SHA256)
	}
	return b, m, nil
}
