# docs-links-resolve-verified

evidence · `01M6ZAAH5RQYHVC4HC6XYE4Y1Z` · actor [claude-code-091083f5](../actor/claude-code-091083f5.md)

## result

PASS

## statement

Verified at 3c6b4c07 as the verifier. Ran go test -trimpath -count=1 -v ./internal/doclinks/...: TestBrokenLinksFail (fixture with a missing file, a missing anchor in another doc, a missing self-anchor and a missing sibling file; each is reported with file:line, while valid file, anchor, self, ../ and external links and fenced code are not flagged), TestCleanFixturePasses, and TestRepoDocsLinksResolve (Check over the real README.md and docs/, which includes docs/quickstart.md#..README.md#install and README's docs/quickstart.md and LICENSE links) all PASS. The test is a plain go test in internal/doclinks, so make check runs it through make test (test-rest covers every package but internal/api, and doclinks is not a budget test). Read the diff: the quickstart lost two links to concepts/why.md and agents.md, files that do not exist in docs/; the new check would otherwise fail on them. The quickstart testscript (go test -trimpath ./cmd/rhiza -run TestScripts/quickstart, with GOTMPDIR moved out of the repo's .git) still PASSes.

## Links

### Out

- verifies [docs-links-resolve](../capability/rhiza/documentation/docs-links-resolve.md)

## History

- `01M6W5EDEX8QAKNVPY5H6CH3Z9` README and quickstart: the README is Rhiza's front door (what and why, install, where next) and docs/quickstart.md walks the model loop without agents, run by a testscript
