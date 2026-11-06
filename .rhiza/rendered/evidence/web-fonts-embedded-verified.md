# web-fonts-embedded-verified

evidence · `01M6Q32XXZVGB1G8G1PTF49HE2` · actor [claude-code-ea82f104](../actor/claude-code-ea82f104.md)

## result

PASS

## statement

At c8d607e2 ran `go test -trimpath -count=1 -run 'TestFonts|TestPagesLoad' -v ./internal/web/`: TestFontsServedFromHandler and TestPagesLoadNothingFromTheNetwork PASS. Read internal/web/fonts_test.go: it GETs /static/fonts/{Petrona,InstrumentSans,JetBrainsMono}-Variable.woff2 from Handler(nil), asserts 200, a font/woff2 content type and the wOF2 magic, reads each OFL-*.txt from the embed FS asserting 'SIL OPEN FONT LICENSE', renders / and fails on any src/href http(s) URL, and checks both CSS files have no url()/@@import to http(s) or protocol-relative URLs. Confirmed by hand: `file` reports all three as WOFF2, the OFL texts are present, web.go embeds static via go:embed, rhiza-tokens.css @@font-face rules use ../fonts/*.woff2, and layout.html references only /-rooted paths.

## Links

### Out

- verifies [web-fonts-embedded](../capability/rhiza/web-monitor/web-fonts-embedded.md)

## History

- `01M6PFFFZNADHH2WGMB5MES8QT` Web design system: the Claude Design system built from the owner's basics as one embedded stylesheet, embedded fonts and the shared layout and partials every page renders with
