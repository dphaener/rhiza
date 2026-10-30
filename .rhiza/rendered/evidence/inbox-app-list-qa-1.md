# inbox-app-list-qa-1

evidence · `01M6752NYN7NQ19JWBA94AJ2QK` · actor [claude-code-558adfdf](../actor/claude-code-558adfdf.md)

## result

PASS

## statement

Built the binary at 8f9c3b9 with go build, outside PATH. In a scratch repo (RHIZA_ROOT, RHIZA_CHANGE, RHIZA_RUN and RHIZA_ACTOR unset) with two accepted change sets, a question run, failed runs (one on an unknown change set) and two real runs waiting at sign-off (started with rhiza run start, fake rex on PATH): rhiza tui app inbox --once matched rhiza inbox --json byte for byte with cmp; change_title appears for known change sets and is left out for the unknown one. I ran rhiza tui app inbox in a real Rex window and drove it with rex send and rex capture. It grouped Questions, Sign-offs and Failed runs, each row with the change set title (the raw ID when the title is unknown) and step, then the question or failure detail. A run journal written from outside appeared within 4s with no key press. A fresh repo showed Nothing waits on you. rhiza rex install into a scratch config dir wrote rhiza_app_inbox (title inbox) and the inbox entry in the app map.

## Links

### Out

- verifies [inbox-app-list](../capability/rhiza/inbox-app-list.md)

## History

- `01M6732DJSQ14B2WTZYZ9MPRA5` An inbox app lists what waits on the owner, questions, sign-offs and failed runs, and acts on each in place: answer, approve, request changes, retry, abandon or open its tab
