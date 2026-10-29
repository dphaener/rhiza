# quick-launch-layer-qa-1

evidence · `01M62YY120BBSCE04K8DX0S7SD` · actor [claude-code-655e5878](../actor/claude-code-655e5878.md)

## result

FAIL

## statement

QA at 4a03cf7 in real Rex: rex new a throwaway session running the built rhiza tui launch. The layer lists ready, roadmap, show, slice, why with summaries. Pressing down+enter on roadmap failed: 'rhiza: tui launch: rex session.view: stream is not attached to session ... (conflict, 409)' and no modal opened. cmd/rhiza/tui.go dials Rex but never calls session.attach, which rexuikit surface.New documents as required. A scratch build adding rex.Attach(ctx, session) after Dial opened a 'roadmap' block showing the roadmap app, so the fix is that one call. The fake-Rex unit tests pass because the fake does not enforce attachment.

## Links

### Out

- verifies [quick-launch-layer](../capability/rhiza/quick-launch-layer.md)

## History

- `01M62JMSW4FH8RAWSQXHC7ZPND` Rex quick-launch layer and install command, with roadmap, ready, show, why and slice as apps
