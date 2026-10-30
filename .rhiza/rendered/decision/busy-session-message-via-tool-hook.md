# busy-session-message-via-tool-hook

decision · `01M65059WCV14J8PHDR2HBA4DJ`

## reasoning

Read internal/api/session_message.go: today the pointer is typed only while the session is idle at its prompt, and the command waits otherwise, so a busy session never receives it mid-task. Typing into a busy terminal would mix the pointer into whatever the agent is doing or into the operator's input, and the file-plus-pointer design already keeps the terminal free of stored content. Read internal/adapter/claudecode/claudecode.go: the adapter already installs post-tool and pre-tool command hooks and writes hook JSON output (hookSpecificOutput with a permission decision), so returning context from a hook fits the existing hook path with no new channel. Chosen over a longer wait because the change intent is delivery without waiting for the prompt, and over typing mid-task because it corrupts terminal input. Not verified: that the agent injects hook additional context from a post-tool hook as described; that needs a check before build.

## statement

A message to a busy session is written to its message file as today, and the pointer to that file is returned to the agent as hook context at its next tool-call boundary, so it arrives mid-task without typing into the terminal. The typed pointer at the prompt stays as the fallback when the session goes idle with no further tool call.

## wording

A message to a busy session is written to its message file as today, and its pointer is returned to the agent as hook context at the next tool-call boundary. The typed pointer at the prompt stays as the fallback when the session goes idle with no further tool call.

## Links

### Out

- answers [question-c315a867-1](../question/question-c315a867-1.md)
- derived-from [busy-session-arbiter-ruling](../source/busy-session-arbiter-ruling.md)

### In

- decided-by [busy-session-hook-delivery](../capability/rhiza/orchestration/step-messages/busy-session-hook-delivery.md)
- decided-by [idle-typed-pointer-kept](../capability/rhiza/orchestration/step-messages/idle-typed-pointer-kept.md)

## History

- `01M61GW821FWSX05A9FC231N65` rhiza session message reaches a busy session without waiting for its prompt, so the owner or arbiter can steer a session mid-task.
