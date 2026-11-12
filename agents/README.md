# agents/

Instructions for the agents that develop Rhiza. This is the index: every
doc in this directory is listed below, one line each, saying which agent
reads it and when. Repo-wide guidance for every agent is in
[AGENTS.md](../AGENTS.md).

## Index

- [orchestrator.md](orchestrator.md): the orchestrator agent reads it before running Rhiza's own development (watch loop, main tree, load, PRs, owner).
- [developing.md](developing.md): any agent changing Rhiza's code or running its tests reads it first (testscripts, scratch binaries, `RHIZA_ROOT`, `GOTMPDIR`, the root daemon, `rhiza render`).
