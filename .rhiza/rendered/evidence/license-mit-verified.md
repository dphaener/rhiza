# license-mit-verified

evidence · `01M6ZMV8JM36GFEYHGPAX060GS` · actor [claude-code-f28c75d7](../actor/claude-code-f28c75d7.md)

## result

PASS

## statement

At 5d28d354 (verifier, attempt 1): diffed root LICENSE against the choosealicense.com MIT template (the text GitHub's licensee matches); the only difference is the copyright line, filled as 'Copyright (c) 2026 Darin Haener'. LICENSE is the sole license file at the repo root, plain ASCII, so licensee's exact-match detection reports MIT. README.md line 31-33 '## License' states '[MIT](LICENSE).'. The leaf commit only adds the builder's actor record; no code changed. GitHub's live detection was not queried (no gh in verifier runs); the check rests on the exact template match.

## Links

### Out

- verifies [license-mit](../capability/rhiza/documentation/license-mit.md)

## History

- `01M6W5F0ZPYJ050A3826ARE92T` Pre-publish audit and MIT license: audit git history and .rhiza/ for secrets, personal data and private hostnames, add the MIT LICENSE, and hand the owner the go-ahead to make the repository public and tag v0.1.0
