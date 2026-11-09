# publish-audit

capability · `01M6W5H42YADPD9H92ZZFBH54D` · parent [documentation](../documentation.md)

## acceptance

an Evidence node records the scanner (gitleaks over full history plus a review of Sources and session notes), its clean result or each finding with its disposition; the owner reads it before flipping visibility

## statement

Before publishing, an audit of the full git history and .rhiza/ (log, Sources, rendered) for secrets, tokens, personal data and private hostnames fixes what the current tree holds and records the result as Evidence, listing anything only a history rewrite would remove

## Links

### Out

- decided-by [m30-audit-then-owner-flips](../../../decision/m30-audit-then-owner-flips.md)
- decided-by [m30-go-public](../../../decision/m30-go-public.md)

### In

- verifies [publish-audit-report](../../../evidence/publish-audit-report.md)
- verifies [publish-audit-verified](../../../evidence/publish-audit-verified.md)
- verifies [publish-audit-qa-1](../../../evidence/publish-audit-qa-1.md)
- verifies [publish-audit-verified-2](../../../evidence/publish-audit-verified-2.md)
- verifies [pr-414-pre-publish-audit-and-mit](../../../evidence/pr-414-pre-publish-audit-and-mit.md)

## History

- `01M6W5F0ZPYJ050A3826ARE92T` Pre-publish audit and MIT license: audit git history and .rhiza/ for secrets, personal data and private hostnames, add the MIT LICENSE, and hand the owner the go-ahead to make the repository public and tag v0.1.0
