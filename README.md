# Rhiza

Rhiza (it's said rise-ah) keeps a living model of your product inside your repo: what it is, why it's that way, and what's planned next. You change the model through change sets, small reviewed steps that land in git alongside the code, so the reasons behind a decision don't live in somebody's head or a chat thread that scrolled away.

It's built for working with coding agents. An agent that can read the model knows what the product is supposed to do and which constraints it can't break, and it can propose changes to the model the same way you do. You can also ignore agents entirely: the CLI works on plain files, with no daemon, and that's how you'll start.

Rhiza was called Genesis until milestone 2. Its older logs keep the names they were written with.

## Install

With Go (1.25 or newer):

```bash
go install github.com/dphaener/rhiza/cmd/rhiza@latest
```

Or install a release binary with the script, which checks the archive's sha256 against the release manifest (Linux or macOS, amd64 or arm64):

```bash
curl -fsSL https://raw.githubusercontent.com/dphaener/rhiza/main/install.sh | sh
```

It installs into `~/.local/bin`; pass `--dir DIR` to change that. There are two channels, each with a manifest listing every archive and its sha256:

- `release` (the default): the newest release, from `https://github.com/dphaener/rhiza/releases/latest/download/manifest.json`.
- `latest`: the newest build of `main`, from `https://github.com/dphaener/rhiza/releases/download/latest/manifest.json`. Pick it with `curl -fsSL https://raw.githubusercontent.com/dphaener/rhiza/main/install.sh | sh -s -- --channel latest`.

Or grab a binary for Linux or macOS (amd64 or arm64) from the [releases page](https://github.com/dphaener/rhiza/releases), unpack it and put `rhiza` on your `$PATH`. Each release has a `checksums.txt` to verify against.

To verify a download by hand, compare `sha256sum` (or `shasum -a 256` on macOS) of the archive with its entry in the channel's manifest or `checksums.txt`. Release and `latest` archives also carry a build provenance attestation: `gh attestation verify rhiza_<version>_<os>_<arch>.tar.gz --repo dphaener/rhiza` checks that GitHub Actions built it from this repository.

Then `rhiza version` should print what you installed.

## Where next

Start with the [quickstart](docs/quickstart.md). It sets Rhiza up in a repo and walks the model loop by hand, no agents needed.

After that, read what you need:

- [Concepts](docs/concepts/model.md): the model, [change sets](docs/concepts/change-sets.md), and [why Rhiza works this way](docs/concepts/why.md).
- [Working with coding agents](docs/agents.md): how an agent reads the model and hands its work back.
- [Runs and the daemon](docs/runs-and-daemon.md): letting Rhiza start agents for you and keep several going while you're away. This is the part that needs Rex and Claude Code; the model itself doesn't.
- [Writing a custom workflow](docs/custom-workflows.md): changing what a run does.
- [Skipping the full CI round for model-only changes](docs/ci-since-last-green.md): a CI recipe for any project.
- Reference, generated from the code: the [CLI](docs/reference/cli.md), [config](docs/reference/config.md) and [workflow file](docs/reference/workflow-file.md).

## Development

`make check` runs everything CI runs. `make install` builds the `rhiza` binary into `~/bin` (set `PREFIX` to change it). `make update` does that, then runs the new binary's `rex install` to rewrite the Rex script and reload Rex. The Go module is `github.com/dphaener/rhiza`. Rhiza keeps its own log in `.rhiza/`, and the environment variables it reads start with `RHIZA_`.

## License

[MIT](LICENSE).
