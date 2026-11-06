# web-design-tokens

capability · `01M6PFJWQ0ZX6WPEX4S94N8ZNW` · parent [web-monitor](../web-monitor.md)

## acceptance

A web package test parses the embedded stylesheet and asserts the colour tokens of both themes equal the token sheet of the delivered Claude Design system (where claude-design-output points), and that no gradient appears in it.

## statement

One embedded stylesheet carries the Claude Design system as CSS custom properties: the dark soil theme by default and the light paper theme under prefers-color-scheme, the gold, teal and rose accents, the Petrona, Instrument Sans and JetBrains Mono type scale, 12-16px radii and 1px rules, with no gradients.

## Links

### Out

- decided-by [web-ui-design-basics](../../../decision/web-ui-design-basics.md)
- decided-by [web-ui-design-from-claude-design](../../../decision/web-ui-design-from-claude-design.md)

### In

- blocks [claude-design-output](../../../question/claude-design-output.md)
- verifies [web-design-tokens-verified](../../../evidence/web-design-tokens-verified.md)
- verifies [web-design-tokens-qa-1](../../../evidence/web-design-tokens-qa-1.md)
- verifies [pr-316-web-ui-design-system](../../../evidence/pr-316-web-ui-design-system.md)

## History

- `01M6PFFFZNADHH2WGMB5MES8QT` Web design system: the Claude Design system built from the owner's basics as one embedded stylesheet, embedded fonts and the shared layout and partials every page renders with
