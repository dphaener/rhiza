# guide-contract-verified

evidence · `01M7AR6AEB4MWG9DVHTMB0533M` · actor [claude-code-e748be38](../actor/claude-code-e748be38.md)

## result

PASS

## statement

At d4010285f (attempt 2): go test -trimpath ./internal/guide -run TestContract exit 0 (TestContract and TestContractRefFailureNamesTopic pass); go test -trimpath ./internal/guide/... exit 0; rhiza ship generate and rhiza render exit 0 with a clean tree. In a scratch copy, removing {cycle} from contract.md makes TestContract exit 1 with 'contract lacks {cycle}', and removing the shared contract-ref line fails with every role topic named. A built binary's rhiza guide contract prints 40 lines listing the six obligations each with its gate, all invoke placeholders ({change},{slug},{title},{run},{step},{brief},{slice},{handback},{worktree},{leaf},{attempt},{check},{cycle}) matching internal/workflow/workflow.go plus the new file placeholders, RHIZA_BRIEF/SLICE/HANDBACK beside the other RHIZA_ variables, and brief.md, slice.json, handback.md; all 17 role topics print 'rhiza guide contract'. Attempt-1 finding f1 (missing placeholders) is fixed and tested. Role topic edits are wording only (to stay under the cap); no obligation changed.

## Links

### Out

- verifies [guide-contract](../capability/rhiza/guide/agent-guide/guide-contract.md)

## History

- `01M79KGHYF3272SZ9H9Y12FWF4` Every agent session a run starts gets one gate-checked contract and its context in files: brief.md, slice.json and, on a hand-back, handback.md, with templated invoke and handback text on agent steps
