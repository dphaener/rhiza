# readme-front-door-verified

evidence · `01M6Z9Y43XS2AEW1Y3FGZ2E49Q` · actor [claude-code-b5c3476b](../actor/claude-code-b5c3476b.md)

## result

PASS

## statement

At 4f10cea8 I read README.md and the diff from ddd443fa. The README explains what Rhiza is and why in a few paragraphs, covers install by go install (Go 1.25+, matching go.mod's go 1.25.14) or release binaries (linux/darwin, amd64/arm64 and checksums.txt, matching .goreleaser.yaml), links the quickstart, and ends with the MIT license. The concepts, agents and CLI reference pages are named as coming, not linked, which fixes attempt 1's dead-link finding. I pulled every link with grep: the relative targets are docs/quickstart.md and LICENSE, and test -e finds both (there are no anchors); the only other link is the GitHub releases URL. The docs-links-resolve check is not built at this commit, so I checked the links by hand. A scratch build runs `rhiza version`, `rhiza describe` and `rhiza --help` as the README says. The prose reads natural, with no AI tells I could find. The owner's sign-off is left to the signoff stage.

## Links

### Out

- verifies [readme-front-door](../capability/rhiza/documentation/readme-front-door.md)

## History

- `01M6W5EDEX8QAKNVPY5H6CH3Z9` README and quickstart: the README is Rhiza's front door (what and why, install, where next) and docs/quickstart.md walks the model loop without agents, run by a testscript
