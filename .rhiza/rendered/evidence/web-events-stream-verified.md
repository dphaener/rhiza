# web-events-stream-verified

evidence · `01M6PVXPG7S5XSJ6NDAMD3Q1DQ` · actor [claude-code-7d5e40c8](../actor/claude-code-7d5e40c8.md)

## result

PASS

## statement

Verifier at d8816f86 (attempt 2). Read the diff: dirFingerprints now fingerprints each change-set directory by its entries' count and last name, which fixes attempt 1's finding f1 (an entry appended to an existing change set sent no change event). TestEventsNameWhatChanged (internal/web/events_test.go) drives web.Events over api.WebEventSources with a locked recorder in a fixture repo. It appends a log entry to a new change set and then to an existing one, moves a run journal, writes driver log lines, changes the daemon's state file, and adds then removes a disk-low inbox item. It reads a matching named event for each (change ×2 for the same ID, run, log, daemon, inbox arrive and leave) and checks the stream holds no stray lines. Ran: go test -trimpath -count=1 ./internal/web/... ok; go test -trimpath -race -count=5 -run TestEvents ./internal/web/ ok (no races); go test -trimpath -count=1 -run 'Web|Inbox' ./internal/api/ ok; go vet and gofmt -l on internal/web and internal/api were clean. TestEventsRefuseWrites gives 405 for POST, PUT and DELETE.

## Links

### Out

- verifies [web-events-stream](../capability/rhiza/web-monitor/web-events-stream.md)

## History

- `01M6PFEXSHWC6TR4F3GFKHV09R` Web server: the daemon serves a read-only web UI on 127.0.0.1:7787 (rhiza.web.port, or off), with read-only internal/api operations as JSON, server-sent events and embedded templates and assets
