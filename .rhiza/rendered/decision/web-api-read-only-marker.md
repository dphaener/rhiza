# web-api-read-only-marker

decision · `01M6PMGNFR7EQWGKHF0EKREYBG`

## quote

Yes a simple marker will do.

## statement

An internal/api operation is exposed by the web UI's /api/<op> only when the registry marks it ReadOnly: a new explicit flag on the operation, set only on operations that change nothing (not the model log, settings, files, Rex or GitHub). /api/ refuses every operation without the marker, so even a hand-typed URL cannot change anything. Write (which only means appends to the log) is not used for this

## Links

### Out

- answers [review-01m6pf30-r3-f1](../question/review-01m6pf30-r3-f1.md)
- derived-from [web-ui-milestone-interview](../source/web-ui-milestone-interview.md)

### In

- decided-by [web-api-read-only](../capability/rhiza/web-monitor/web-api-read-only.md)

## History

- `01M6PFEXSHWC6TR4F3GFKHV09R` Web server: the daemon serves a read-only web UI on 127.0.0.1:7787 (rhiza.web.port, or off), with read-only internal/api operations as JSON, server-sent events and embedded templates and assets
