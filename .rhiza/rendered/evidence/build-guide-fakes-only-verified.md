# build-guide-fakes-only-verified

evidence · `01M65FXRQ9A5FGQVDH3J1YJDX2` · actor [claude-code-3bd44bfe](../actor/claude-code-3bd44bfe.md)

## result

PASS

## statement

At 6f27dd0 I read the diff: internal/guide/text/build.md's Verifier section now says 'Test runs with fakes only, never real Rex, agents or `gh`', and TestBuildText in internal/guide/build_test.go pins that phrase under a 'fakes-only' key. go test ./internal/guide -run 'TestBuildText|Cap' -count=1 -v passed. TestCaps showed build at 2344 bytes, 586 tokens, under TopicCap 800. The rest of go test ./... passed except internal/tui/modelapp TestGeneratedModel, a read-timing check (median ~63ms, want <50ms) that this text-only change cannot affect. make self-check flagged a stale render for the builder's landed actor record, which is outside this leaf.

## Links

### Out

- verifies [build-guide-fakes-only](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-guide/build-guide-fakes-only.md)

## History

- `01M61GW73H05YHE4DK0CXFD1KH` No agent session starts a run: rhiza run start refuses inside an agent session (RHIZA_SESSION set), naming the session, and the build and review guides tell verifiers and checks to test runs with fakes only, never real Rex, agents or gh.
