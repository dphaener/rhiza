# web-spacing-rules-verified

evidence · `01M7B4P0B86J9E00N14CKE9NNW` · actor [claude-code-24dbb4b1](../actor/claude-code-24dbb4b1.md)

## result

PASS

## statement

At cd0c52dab re-ran the Checks block: go test -trimpath ./internal/web -run 'TestPartialsGolden|TestSpacingRules' (exit 0), go test -trimpath ./internal/web (exit 0), rhiza ship generate (exit 0), rhiza render (exit 0). Read the diff: rhiza.css sets .rz-pagehead > div gap var(--space-3), .rz-panel-head > div gap var(--space-1), .rz-panel-head margin-bottom var(--space-4), .rz-stack/.rz-grid gap var(--space-5), .rz-main gap var(--space-6), .rz-panel padding var(--space-4) with var(--space-5) in @@media (min-width: 721px); inline margins removed from change.html and daemon.html, goldens regenerated; TestSpacingRules parses the stylesheet and checks each value, heading margins and inline template margins. Mutation in a scratch copy: .rz-stack gap at var(--space-4) (16px) fails TestSpacingRules; a margin-top on .rz-pagehead h1 and an inline style margin in change.html each fail it. No TODOs, nothing outside internal/web and the leaf's .rhiza records changed.

## Links

### Out

- verifies [web-spacing-rules](../capability/rhiza/web-monitor/web-spacing-rules.md)

## History

- `01M79NKERTK62YZM9C42BJ0JW3` The web UI's frame: one header and navigation for every screen with a real project picker (search, each project's state and needs-me count) in place of the header's project list, and shared components on one spacing and type scale
