# web-fonts-embedded

capability · `01M6PFJWR6X3ZNEVCZ118QA4SX` · parent [web-monitor](../web-monitor.md)

## acceptance

A web package test GETs each font from the handler with a font content type, asserts each licence file is embedded, and asserts the rendered layout references no http(s) URL outside the page's own origin.

## statement

Petrona, Instrument Sans and JetBrains Mono are embedded in the binary as woff2 files with their OFL licences and served under /static/, so pages load nothing from the network.

## Links

### Out

- decided-by [web-ui-design-basics](../../../decision/web-ui-design-basics.md)

### In

- blocks [claude-design-output](../../../question/claude-design-output.md)
- verifies [web-fonts-embedded-verified](../../../evidence/web-fonts-embedded-verified.md)
- verifies [web-fonts-embedded-qa-1](../../../evidence/web-fonts-embedded-qa-1.md)
- verifies [pr-316-web-ui-design-system](../../../evidence/pr-316-web-ui-design-system.md)

## History

- `01M6PFFFZNADHH2WGMB5MES8QT` Web design system: the Claude Design system built from the owner's basics as one embedded stylesheet, embedded fonts and the shared layout and partials every page renders with
