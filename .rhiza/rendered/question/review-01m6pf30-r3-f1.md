# review-01m6pf30-r3-f1

question · `01M6PGBJM14G29CMMYB19YSYRS`

## statement

Plan review round 3 left finding f1 (contradiction on web-api-read-only): The acceptance has a test 'walk the api registry asserting every allowlisted operation is a read', but the registry has no read marker. Operation.Write (internal/api/api.go) only records whether an operation appends to the log. daemon.set, daemon.start, render and ship.merge all have Write false yet write config, files or GitHub, so a Write-based assertion would pass for them. The leaf needs to name what makes an operation a read (an explicit allowlist checked against a named set, or a new registry flag) so the check can actually fail. The critic confirmed it: internal/api Operation has no read marker: Write only means 'appends to the log', and daemon.set, daemon.start, render and ship.merge all leave Write false while writing git config, files or GitHub. The acceptance's 'every allowlisted operation is a read' has no defined test, so a Write-based check would pass for those writes; the leaf must name what makes an operation a read. How should it be settled?

## Links

### Out

- blocks [web-api-read-only](../capability/rhiza/web-monitor/web-api-read-only.md)

### In

- answers [web-api-read-only-marker](../decision/web-api-read-only-marker.md)

## History

- `01M6PFEXSHWC6TR4F3GFKHV09R` Web server: the daemon serves a read-only web UI on 127.0.0.1:7787 (rhiza.web.port, or off), with read-only internal/api operations as JSON, server-sent events and embedded templates and assets
