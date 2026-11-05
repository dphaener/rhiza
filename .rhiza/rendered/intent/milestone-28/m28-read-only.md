# m28-read-only

intent · `01M6PFEK8PKRGYC4C6WZS9BY5Y` · parent [milestone-28](../milestone-28.md)

## statement

Nothing the web UI serves writes: only allowlisted read operations of internal/api are reachable, it binds only to 127.0.0.1, and turning it off (rhiza.web.port off) leaves the daemon as it was.

## History

- `01M6PF31YCJ041HB94MDWQR331` I want to build out the first pass of a web UI that can be used to monitor the progress of Rhiza. This will run as part of the daemon process and be served on localhost. The user will be responsible for any port forwarding (if it's running on this VM for instance). First pass will be read only. For monitoring purposes only and we will build from there. The first requirement will be to obtain a design. You will generate a prompt that I will then feed to Claude design to generate a design system.
