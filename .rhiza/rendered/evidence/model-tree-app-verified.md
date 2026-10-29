# model-tree-app-verified

evidence · `01M6351WN1JEQZ9966DEQ6RVC1` · actor [claude-code-6eb717bd](../actor/claude-code-6eb717bd.md)

## result

PASS

## statement

Verifier at a3a2efd: wrote gen.Default() (the 50,000-node model) to a scratch repo and built rhiza from this commit. `rhiza tui app model --once` equals `rhiza tree --view as-is` + blank line + `rhiza roadmap` byte for byte (cmp: 15,938 tree lines, 37 roadmap lines); with the deep-subtree ref (the tree read the api budget D14 times) `rhiza tui app model <ref> --once` equals `rhiza tree <ref> --view as-is` + roadmap likewise. Whole-process wall time for the subtree+roadmap app read is 0.03-0.04 s over 5 runs; TestGeneratedModel's in-process median is 10 ms (< 50 ms). Note: with no ref the whole 50k tree takes about 0.2 s, the same as `rhiza tree` itself, and is outside the budgeted read as D14 defines it. Unit tests (tab switch, scroll, quit, error shown), tui-model-app.txtar (no socket, no daemon lock) and make check all pass (exit 0). Diff read: in-process api operations, no daemon or socket.

## Links

### Out

- verifies [model-tree-app](../capability/rhiza/model-tree-app.md)

## History

- `01M62JMT63RY271EAF14R8P4VG` A Rex app showing the as-is capability tree and the roadmap
