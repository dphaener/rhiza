# plan-review-check-rulings-qa-3

evidence · `01M7MTBDPM606938HFYYBV5H5H` · actor [claude-code-60928770](../actor/claude-code-60928770.md)

## result

PASS

## statement

QA round 3 at db1b1c925, scratch build (go build -trimpath, run by full path, outside PATH) driving autonomous plan runs by hand in scratch Rhiza roots with the plan-workflow-autonomous* fake rex/claude/gh: from the apply-rulings step --rulings exited 1 listing both questions (decision, finding's node export-csv/keep-flat); after linking only keep-flat decided-by it exited 1 listing only export-csv; after linking export-csv it exited 0 and --json listed both rulings with by IDs; after retiring export-csv (its only decided-by link) it exited 1 again; after adding export-csv-2 linked decided-by it exited 0. A clean review ran rulings visit 1 pass then gate. RHIZA_RUN unset, unknown or finished: refused, exit 2.

## Links

### Out

- verifies [plan-review-check-rulings](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-review-check-rulings.md)

## History

- `01M7MPAHJZY3RV8ANGGB660J52` Fix Plan run arbiter rulings are never delivered: a plan run whose last review round leaves a confirmed major puts it to the arbiter, the arbiter answers with a decision handle (seen: parity-seed-bridge-in-model in run 01M7GG7GXSA97HESRNFEBFCZ4T, split-check-whatif-once in run 01M7GG7HSBSB356BVGGQGH8YHB), then the run merges and accepts the plan with the question marked 'decision ..., not yet delivered' and no Decision node or plan change in the model. The accepted plan still carries the confirmed major. Expected: the ruling is written (Decision + reviser applies it) before merge, or the run fails instead of accepting.
