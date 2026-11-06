# root-daemon-round-robin-verified

evidence · `01M6SBY5H79TCP8PVKDH1B8TSR` · actor [claude-code-347ccaf4](../actor/claude-code-347ccaf4.md)

## result

PASS

## statement

Verifier claude-code-347ccaf4 at 158ebf7c ran go test -trimpath -count=1 -v -run RoundRobin ./internal/daemon/ (both PASS) plus the whole internal/daemon package, go vet and gofmt (clean). Read the diff: RoundRobin(projects, running, maxRuns) gives one turn per project in order, each taking that project's next Ready change set in its roadmap order, skipping a project with nothing left or at Running+taken >= its Max (0 = no cap), stopping when slots are full or no progress. Tests assert maxRuns 3 over A(3), B(2), C(1) picks A1,B1,C1; A max 1 with maxRuns 4 picks A1,B1,C1,B2 (A one slot, rest to B and C); also maxRuns 6, already-running slots and A already at its max.

## Links

### Out

- verifies [root-daemon-round-robin](../capability/rhiza/orchestration/daemon/root-daemon-round-robin.md)

## History

- `01M6RN3RCGAZ9RNTWTZPKFM61M` Root daemon: one daemon on the root builds every registered project round-robin under one machine-wide run cap and usage governor, with machine settings in the root config and project settings in each repo, under one service
