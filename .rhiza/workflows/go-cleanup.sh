#!/bin/sh
# Cleanup phase of Rhiza's own workflows: trims what Go and rhiza's tests
# leave behind. GOCACHE and TMPDIR (default /tmp) place the folders; both
# are read from the run's environment.
#   - Go build cache entries unused for three hours (go refreshes an
#     entry's mtime on use; its own trim waits five days),
#   - go-build* temp folders older than two hours,
#   - rhiza-* test temp folders older than an hour.
cache=${GOCACHE:-$(go env GOCACHE 2>/dev/null)}
tmp=${TMPDIR:-/tmp}
if [ -n "$cache" ] && [ "$cache" != off ] && [ -d "$cache" ]; then
	find "$cache" -mindepth 2 -maxdepth 2 -type f -path '*/[0-9a-f][0-9a-f]/*' -mmin +180 -delete
fi
if [ -d "$tmp" ]; then
	find "$tmp" -mindepth 1 -maxdepth 1 -type d -name 'go-build*' -mmin +120 -exec rm -rf {} +
	find "$tmp" -mindepth 1 -maxdepth 1 -type d -name 'rhiza-*' -mmin +60 -exec rm -rf {} +
fi
exit 0
