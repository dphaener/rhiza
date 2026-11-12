// Command release writes manifest.json for a GoReleaser dist folder. It fails
// unless all four platforms are present and match checksums.txt.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dphaener/rhiza/internal/release"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) > 0 && args[0] == "version" {
		v, err := gitLatestVersion(".")
		if err != nil {
			return err
		}
		fmt.Println(v)
		return nil
	}
	if len(args) > 0 && args[0] == "breaking" {
		tag, bs, err := breakingAgainstLastTag(".")
		if err != nil {
			return err
		}
		for _, b := range bs {
			fmt.Println(b)
		}
		fmt.Fprintf(os.Stderr, "%d breaking difference(s) since %s\n", len(bs), tag)
		return nil
	}
	fs := flag.NewFlagSet("release", flag.ContinueOnError)
	var o release.Options
	fs.StringVar(&o.Dist, "dist", "dist", "GoReleaser dist folder")
	fs.StringVar(&o.Project, "project", "rhiza", "archive name prefix")
	fs.StringVar(&o.Tag, "tag", "", "git tag, e.g. v0.1.0")
	fs.StringVar(&o.Version, "version", "", "archive version (default: tag without v)")
	fs.StringVar(&o.Commit, "commit", "", "commit sha")
	fs.StringVar(&o.Repo, "repo", "dphaener/rhiza", "GitHub owner/name")
	fs.StringVar(&o.Channel, "channel", "", "manifest channel (default: release)")
	fs.StringVar(&o.ReleaseTag, "release-tag", "", "GitHub release the archives are on (default: tag)")
	out := fs.String("out", "", "output file (default: <dist>/manifest.json)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	m, err := release.Build(o)
	if err != nil {
		return err
	}
	b, err := m.Marshal()
	if err != nil {
		return err
	}
	if *out == "" {
		*out = filepath.Join(o.Dist, release.ManifestFile)
	}
	return os.WriteFile(*out, b, 0o644)
}
