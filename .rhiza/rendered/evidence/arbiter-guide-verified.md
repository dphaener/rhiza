# arbiter-guide-verified

evidence · `01M5XCZZWHE5FDS228ATJ6HDVQ` · actor [claude-code-54c93d5c-verifier](../actor/claude-code-54c93d5c-verifier.md)

## result

PASS on 2026-10-26 at commit 9b213db (code as of 8b72ad5) on branch cs/autonomous-mode-and-the, re-run with -count=1 by an independent verifier (claude-code-54c93d5c-verifier)

## statement

Testscript guide checks rhiza guide ls lists eight texts including '- arbiter topic N tokens', and rhiza guide arbiter prints '# Arbiter', the apply recipe with wording=(assumed, would have escalated), reasoning=, the link {"op":"link","from":"<handle>","type":"answers","to":"<question id>"} and rhiza session done. internal/guide TestCaps keeps every text, arbiter included (386 tokens), within its cap, and the topic list test includes arbiter.

## tests

go test -count=1 ./internal/guide; go test -count=1 ./cmd/rhiza -run 'TestScripts/guide$'

## Links

### Out

- verifies [arbiter-guide](../capability/rhiza/workflows/arbiter-sessions/arbiter-guide.md)

## History

- `01M5WPTVT93DTC076WJNB7E2KX` Autonomous mode and the arbiter: a run can skip the owner's gates and questions, recording them as rule autonomous and (assumed, would have escalated), and the engine calls a fresh arbiter session at a decision point; agent permissions are never bypassed.
