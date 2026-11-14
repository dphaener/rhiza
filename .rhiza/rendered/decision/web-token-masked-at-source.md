# web-token-masked-at-source

decision · `01M7B52TYHMGX5ZFARNQK3TSHC`

## reasoning

The finding is confirmed in the code: stateRows() (internal/api/daemon_settings.go) copies every setting Value into daemon state, and daemon.status (ReadOnly, internal/api/daemon.go ~2057) returns those Values, served without authentication by ReadOnlyHandler (web_api.go). @web-token already says the token's value is never shown, so the intent is clear and only its coverage is incomplete. Masking per operation would leave the token in the on-disk daemon state and invite the next read-only op to leak it again; masking once where rows are built (settings rows and stateRows) closes daemon.settings, daemon.status and the state together, and needs no new scope beyond the change set's intent of a required web token. The acceptance must name daemon.status so a builder who masks only daemon.settings fails.

## statement

The web token's value is masked where setting rows are built, not per operation: rhiza.web.token's row carries only whether it is set, in daemon.settings, in the rows the daemon's poll copies into its state (conf.stateRows()) and so in daemon.status's Settings, and in every read-only operation served unauthenticated at GET /api/<op>. The daemon's web handler reads the token itself from the setting, never from state or a response. @{01M79Q8PX68N5R6XWD5MHV6FKJ}'s acceptance gains a check that daemon.status and an unauthenticated GET /api/daemon.status through WebProjectHandler do not return the token's value.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m7b3wjc309vgwck1ef6ms5gb-r3-f1](../question/review-01m7b3wjc309vgwck1ef6ms5gb-r3-f1.md)

## History

- `01M79NNP2BRD9R8XG15XCYACZ4` Safe actions in the web UI behind a required web token: answer an inbox item, approve a sign-off, pause or resume the daemon or a project; nothing else
