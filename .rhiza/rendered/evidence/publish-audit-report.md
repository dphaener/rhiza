# publish-audit-report

evidence · `01M6ZQ504B8KYS3VAZA4FBP90T` · actor [claude-code-411ae041](../actor/claude-code-411ae041.md)

## result

PASS

## statement

Pre-publish audit at 5d28d354 (main plus this branch), amended on attempt 2. SCANNER: gitleaks 8.18.4 with default rules over full history (`gitleaks detect --log-opts=--all`, 887 commits, 6.9s): no leaks found; the attempt-1 verifier re-ran it over history (888 commits) and the tree (--no-git) with the same clean result. MANUAL SWEEP of HEAD (git grep over tree, .rhiza/log, .rhiza/rendered incl. Sources and session notes) for password/secret/api key/private key/ghp_/github_pat_/sk-/sk-ant-/AKIA/xox tokens, emails (`git grep -hoE '<email regex>' HEAD | sort | uniq -c`), IPs, home paths, private hostnames, phone numbers; and history for deleted .env/.pem/.key/.db/credentials files and blobs over 1 MB: no secrets or credentials found (only hits: a GITHUB_TOKEN secret reference in release.yml, the word 'secret' in scrub docs and a test fixture SECRET-RAW; the two 1 MB blobs are describe-golden.txtar). FINDINGS: (1) docs/rex-upstream/notification-banner.md named the owner's private tailnet HTTPS hostname: FIXED in commit 36850ef3 (now 'a tailnet HTTPS name'); it remains in earlier history. (2) Two of the owner's personal email addresses (the gmail one and the hey.com full-name one): ACCEPTED, no action; they are the git author/committer identity of 913 commits and appear in Actor records in .rhiza/log, so they are public once the repo is, removable only by a history rewrite. (3) A third owner address (the short hey.com one) was used as a sample email in internal/textref/textref.go (doc comment) and internal/textref/textref_test.go (two test inputs): FIXED in the publish-audit attempt-2 commit, replaced by ada@example.com (go test ./internal/textref passes). Capability text-ref-syntax's acceptance named it as well: FIXED going forward in review cycle 2 by a set of that acceptance in this change set (now ada@example.com, matching the code), which clears it from the current render of text-ref-syntax. It remains in 8 .rhiza/log entries of the text-references change set 01M6HNNWFR1J2AHQRNKFTPME46, in the 4 evidence renders derived from them (text-ref-syntax-qa-2, text-ref-syntax-verified, -verified-2, -verified-3; evidence statements record what was run and are not rewritten), and in older history of the two Go files: ACCEPTED as a public contact address; only a history rewrite removes it from .rhiza/log (append-only) and git history. (4) A Tailscale CGNAT address and the owner's machine name in three session notes (and their log entries and renders): ACCEPTED as low risk (reachable only inside the owner's tailnet); not edited because .rhiza/log is append-only and renders derive from it. (5) The owner's macOS username in absolute home paths in 18 files (log entries, renders, evidence): ACCEPTED, discloses a username only. (6) Other emails, IPs and hosts found (test@rhiza.invalid, me@x.org and x@*.example in test fixtures, a 192.168.x test fixture, daemon.home file name, example.com/example.test/invalid addresses) are test fixtures or file names, not private hosts. ONLY A HISTORY REWRITE WOULD REMOVE: the owner emails in every commit's author/committer fields; the tailnet hostname in older versions of notification-banner.md; the third owner address in older versions of textref.go/textref_test.go and in the text-references change set's .rhiza/log entries; and findings 4 and 5 in .rhiza/log (append-only) and its history. Recommendation: none of these is a credential; the owner decides whether to publish as is. The owner reads this before flipping visibility.

## Links

### Out

- verifies [publish-audit](../capability/rhiza/documentation/publish-audit.md)

## History

- `01M6W5F0ZPYJ050A3826ARE92T` Pre-publish audit and MIT license: audit git history and .rhiza/ for secrets, personal data and private hostnames, add the MIT LICENSE, and hand the owner the go-ahead to make the repository public and tag v0.1.0
