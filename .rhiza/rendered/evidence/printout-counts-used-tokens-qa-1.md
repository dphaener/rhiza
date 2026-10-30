# printout-counts-used-tokens-qa-1

evidence · `01M65C37J0P2YP6EZK9BQ8XX8S` · actor [claude-code-a8fce9bb](../actor/claude-code-a8fce9bb.md)

## result

PASS

## statement

QA at 6b63cac with a binary built from this checkout (go build -o $TMPDIR/rhiza-qa ./cmd/rhiza). On the live review run 01M65B5R7VG5QDFMFFVT012N3R, the installed rhiza 2b1628d printed 'no token use recorded' for all five running review sessions; the new binary listed each with real counts read from its Claude transcript, and Totals rose from input 156 to 196 (156+8+8+10+6+8). run status --json listed 9 of 9 sessions with a tokens object. In a scratch repo, a two-step run with a fake claude transcript (7/11/13/17) printed both sessions and Totals 14/22/26/34 in its final printout. After every tokens entry was stripped from the journal, run status still showed both sessions' counts and the same Totals, and the journal's checksum was unchanged after the reads (nothing journaled).

## Links

### Out

- verifies [printout-counts-used-tokens](../capability/rhiza/runtime/session-journal/run-status-tokens/printout-counts-used-tokens.md)

## History

- `01M61GW8RKSBM4BT1B8B3FR5JH` A run's final printout and run status show each session's token use, never 'no token use recorded' for sessions that used tokens.
