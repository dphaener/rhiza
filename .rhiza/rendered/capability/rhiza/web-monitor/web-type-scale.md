# web-type-scale

capability · `01M7B1K2MDGPBDPTQ7EFZYG1NH` · parent [web-monitor](../web-monitor.md)

## acceptance

go test -trimpath ./internal/web -run TestTypeScale passes: it parses static/css/rhiza-tokens.css and finds exactly the size tokens 11, 12, 13, 14, 15, 16, 17, 20, 22, 26, 28, 32 and 36px and the line-height tokens 1, 1.6, 16, 18, 20, 22, 26, 30, 34, 36 and 40px; it fails on any font-size, line-height or font shorthand in static/css/rhiza.css or a template's inline style whose size or line height is not var() of one of those tokens; and it reads every font declaration in pages/*.html and fails if one sets a size or line height the tokens lack, so the scale cannot drop a size an approved design uses (failure cases: a 28px/32px stat value fails, and a 10.5px row head fails).

## statement

Every web page sets its text from one type scale, defined once as size and line-height tokens in static/css/rhiza-tokens.css and taken from what the design pages (pages/*.html) set. Display face: page title 36/40 on desktop and 32/36 on a phone, quote 26/34 on desktop and 22/30 on a phone, detail and section title 28/34, panel title 20/26, step and empty-state title 17/22, wordmark 20/1 on desktop and 22/1 on a phone. Body face: lead 17/1.6, project name 16/22, body 15/22, rows and tables 14/20, meta 13/18, label 12/16, buttons and controls 15/1 and 14/1, and 13/20 for inline evidence. Mono face: ids and numbers 13/20, mono meta 13/18, times 12/20 and 12/18, crumbs 12/16, counts and pills 12/1, kickers, row heads and severities 11/16 and 11/20, the wordmark's tag 11/1. So the sizes are 11, 12, 13, 14, 15, 16, 17, 20, 22, 26, 28, 32 and 36px, and the line heights 1, 1.6, 16, 18, 20, 22, 26, 30, 34, 36 and 40px. static/css/rhiza.css's sizes outside them move to the step of the same role: 10.5px kickers, row heads, severities and wordmark tag to 11px, 11.5px pills to 12px, 12.5px row steps, times and meta to 12px (12/20 and 12/18), 13.5px history rows to 14/20, the 28/32 stat value to 28/34, and the phone page and detail titles 30/34 and 24/30 to 32/36 and 28/34. static/css/rhiza.css and the templates set no size or line height outside the tokens.

## Links

### Out

- decided-by [web-spacing-scale](../../../decision/web-spacing-scale.md)

### In

- verifies [web-type-scale-verified](../../../evidence/web-type-scale-verified.md)
- verifies [web-type-scale-qa-1](../../../evidence/web-type-scale-qa-1.md)
- verifies [web-type-scale-qa-2](../../../evidence/web-type-scale-qa-2.md)
- verifies [pr-571-the-web-ui-s-frame](../../../evidence/pr-571-the-web-ui-s-frame.md)

## History

- `01M79NKERTK62YZM9C42BJ0JW3` The web UI's frame: one header and navigation for every screen with a real project picker (search, each project's state and needs-me count) in place of the header's project list, and shared components on one spacing and type scale
