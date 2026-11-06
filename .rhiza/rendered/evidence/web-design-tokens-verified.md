# web-design-tokens-verified

evidence · `01M6Q31B9W6RPGXKVWAQ8Z4FCY` · actor [claude-code-e29ed6df](../actor/claude-code-e29ed6df.md)

## result

PASS

## statement

At 89f827a8 ran go test -trimpath ./internal/web/ -run 'TestTokens|TestStylesheets' -v -count=1: TestTokensEqualDeliveredDesignSystem and TestStylesheetsHaveNoGradients PASS. The test reads the go:embed'd static/css/rhiza-tokens.css and checks all 22 colour tokens of the dark default (:root, [data-theme=dark]), light under @@media (prefers-color-scheme: light) and forced [data-theme=light], plus radius-md 12px, radius-lg 16px, rule 1px and the Petrona / Instrument Sans / JetBrains Mono faces. I compared its pinned maps by hand with rhiza-tokens.css at 1f4b22c (the delivered system that claude-design-output's answer web-design-system-in-internal-web points to): every value matches, and git diff 1f4b22c..HEAD shows the CSS unchanged. A mutation run on a scratch copy (light teal #2b7365, radius-md 10px, a linear-gradient in rhiza.css) failed all three checks, as it should. Not covered: the --glow shadow token is not asserted.

## Links

### Out

- verifies [web-design-tokens](../capability/rhiza/web-monitor/web-design-tokens.md)

## History

- `01M6PFFFZNADHH2WGMB5MES8QT` Web design system: the Claude Design system built from the owner's basics as one embedded stylesheet, embedded fonts and the shared layout and partials every page renders with
