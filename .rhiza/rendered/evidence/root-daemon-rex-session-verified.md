# root-daemon-rex-session-verified

evidence · `01M6SF1F2Z2CPRQ29YW1G6DHAP` · actor [claude-code-1b83b35c](../actor/claude-code-1b83b35c.md)

## result

PASS

## statement

At dc20d2c6 ran go test -trimpath -count=1 ./cmd/rhiza -run TestScripts/daemon-rex-session (PASS), plus the daemon*, run-open*, inbox-open* and root* scripts and go test ./internal/api/... (all PASS). The script builds a root with rhiza root init, registers alpha with rhiza project add and folder with rhiza project add folder --name other, unsets REX_SESSION and fakes rex: rhiza daemon --once creates rex session new rhiza-alpha and rhiza-other exactly once each and never rhiza-folder; each detached run driver line carries REX_SESSION=rhiza-alpha / REX_SESSION=rhiza-other; with both sessions listed by rex session ls a second --once creates none; rhiza run open --run in the renamed project, for a live (locked) run with a recorded tab and no session, issues rex block inspect -s rhiza-other and rex focus -s rhiza-other. Diff read: repoSessionName prefers Ctx.ProjectName (set in newDaemonProj) then the root registry's ProjectOf name, falling back to the folder; run open falls back to it when no session is recorded.

## Links

### Out

- verifies [root-daemon-rex-session](../capability/rhiza/orchestration/daemon/root-daemon-rex-session.md)

## History

- `01M6RN3RCGAZ9RNTWTZPKFM61M` Root daemon: one daemon on the root builds every registered project round-robin under one machine-wide run cap and usage governor, with machine settings in the root config and project settings in each repo, under one service
