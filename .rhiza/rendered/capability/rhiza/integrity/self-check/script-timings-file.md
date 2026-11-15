# script-timings-file

capability · `01M7EQE65YZAMRMTAKKJN4A7A6` · parent [self-check](../self-check.md)

## acceptance

go test -trimpath ./tools/scripttimings shows: a go test -json fixture of TestScripts subtests, and a gh run view --log fixture of script-time lines from four scripts jobs, each merge into the previous file as sorted name<TAB>seconds lines, a script that failed or did not run keeping its previous time and a script no longer under cmd/rhiza/testdata/script dropped; failure case: input with no passing TestScripts subtest and no script-time line exits non-zero and leaves the file unchanged, as does a malformed previous file. The builder checks in the file with a line for every cmd/rhiza/testdata/script/*.txtar at 1 second (@script-timings-seed-then-ci). Refresh (orchestrator, before merge): make script-timings RUN=<this PR's CI run> rewrites it from the scripts jobs, recorded after the virtual clock, and that commit lands with the change set. Touches tools/scripttimings/, Makefile, cmd/rhiza/testdata/script-timings.tsv. Must not add a module dependency.

## statement

cmd/rhiza/testdata/script-timings.tsv, checked in, holds each script's recorded wall time (name, tab, seconds, sorted by name), and make script-timings rewrites it through a small Go tool under tools/scripttimings from the script-time lines of a CI run's scripts jobs (RUN=<run id>, read with gh run view --log) or from a go test -json file of TestScripts (JSON=<file>).

## Links

### Out

- decided-by [fast-suite-proposal-approved](../../../../decision/fast-suite-proposal-approved.md)
- decided-by [script-timings-seed-then-ci](../../../../decision/script-timings-seed-then-ci.md)

### In

- depends-on [scripts-longest-first](scripts-longest-first.md)
- depends-on [script-shards-weighted](script-shards-weighted.md)
- depends-on [ci-script-timings](ci-script-timings.md)
- verifies [script-timings-file-verified](../../../../evidence/script-timings-file-verified.md)
- verifies [script-timings-file-qa-1](../../../../evidence/script-timings-file-qa-1.md)

## History

- `01M7EQ3K7PW2XYR684CXSX74TG` Parallel scripts, longest first, weighted CI shards: cmd/rhiza's scripts run at a -parallel bounded by cores and memory, start longest first, and CI's script shards split by a checked-in timing file instead of round-robin by name. Evidence: scripts use ~1.2 cores at -parallel 4; at -parallel 12 they take 381 s instead of 639 s with 0 failures; CI shards took 61-239 s.
