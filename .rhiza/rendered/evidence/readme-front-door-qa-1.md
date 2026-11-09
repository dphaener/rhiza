# readme-front-door-qa-1

evidence · `01M6ZATPVAHJPFVZR1APDK9KRF` · actor [claude-code-d4acd8b3](../actor/claude-code-d4acd8b3.md)

## result

PASS

## statement

QA at b92987bf: read README.md. It covers what and why, install via go install github.com/dphaener/rhiza/cmd/rhiza@latest (matches go.mod's module and go 1.25) or release binaries for linux/darwin amd64/arm64 with checksums.txt (matches .goreleaser.yaml), a link to docs/quickstart.md, a note that the concepts, agents and reference pages are still to come, and an MIT LICENSE file. The built binary's version, --help, guide and describe commands that it points to all run. Its relative links resolve (doclinks test). The owner still has to sign off.

## Links

### Out

- verifies [readme-front-door](../capability/rhiza/documentation/readme-front-door.md)

## History

- `01M6W5EDEX8QAKNVPY5H6CH3Z9` README and quickstart: the README is Rhiza's front door (what and why, install, where next) and docs/quickstart.md walks the model loop without agents, run by a testscript
