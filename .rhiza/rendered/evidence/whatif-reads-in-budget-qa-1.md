# whatif-reads-in-budget-qa-1

evidence · `01M6CQG056PEC18R6JAYR7CPNV` · actor [claude-code-07a1d7e1](../actor/claude-code-07a1d7e1.md)

## result

PASS

## statement

QA at eaf04bd9 with the built binary (go build ./cmd/rhiza) against the merge-base 3ef3d9c5 binary. (1) Scratch repo: 270 as-is nodes, accepted change set A (add under cap-2, set, link, change-set serves, retire) and draft B depending on A (add under A's node, set over A's set, rename handle, mv to new parent, mv to root, retire node with children, retire constraint, question blocks, decided-by, unlink, intent child). 50 reads (show by handle/ID/prefix/retired/renamed/link, ls with kind/parent/retired/limit, tree incl. --retired/--depth/--kind, why, log, roadmap, ready, search, slice) x what-if:A and what-if:B, --json: all 100 outputs and exit codes byte-identical between old and new; B's view shows its writes (tree has cache-evict and cap-four, not cap-1; cap-12 under cap-9; cap-3 statement from B). (2) 50,000-node model from internal/gen.Default via the CLI: 26 reads in what-if of the draft and of the accepted change set identical between binaries. Wall-clock CLI medians of 7 runs in what-if:<draft>, process start included: show 33.1, ls 33.2, tree 35.3, why 33.8, log 34.4, roadmap 37.1, ready 37.1 ms (old binary 226-290 ms); as-is 27-32 ms. go test ./internal/api -run 'TestReadBudget|TestReadBackendsAgree' and ./internal/index -run TestPropertyOverlayAgrees passed.

## Links

### Out

- verifies [whatif-reads-in-budget](../capability/rhiza/views/fast-reads/whatif-reads-in-budget.md)

## History

- `01M5SWTFB3ZAX4YKX4KA2ZVSZR` An incremental what-if folds one change set over the indexed as-is, so what-if reads meet the 50 ms target too.
