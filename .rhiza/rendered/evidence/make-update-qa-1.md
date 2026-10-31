# make-update-qa-1

evidence · `01M6995GXH5V2GFBG0ZPX62987` · actor [claude-code-083d55ec](../actor/claude-code-083d55ec.md)

## result

PASS

## statement

At 6eb01e4, ran the real make update PREFIX=<mktemp>/bin with REX_CONFIG=<mktemp>/rexcfg/init.lua (holding one owner line) and the real rex on PATH. It ran go build -o <tmp>/bin/rhiza, then <tmp>/bin/rhiza rex install, with no cp. Output: 2 written, Rex reloaded, exit 0. The installed binary runs (rhiza version prints 6eb01e4, ad-hoc signature intact). init.lua kept the owner line and gained require("rhiza"), and rhiza.lua was written beside it. A second make update rebuilt and reported 0 written, 2 unchanged, Rex reloaded. With rex missing from PATH, make update failed (Error 1, exit 2) with 'rex config reload failed' and next: rex config reload. TestMakeUpdateDryRun PASSED.

## Links

### Out

- verifies [make-update](../capability/rhiza/rex-install-command/make-update.md)

## History

- `01M686RWSEKXXBCNC7RZMT07HB` make update rebuilds rhiza in place, rewrites the Rex script and reloads Rex in one command
