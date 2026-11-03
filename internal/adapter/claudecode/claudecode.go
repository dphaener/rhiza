// Package claudecode is the adapter for Claude Code, the first built-in
// adapter and the default. It is the only package that names the agent.
package claudecode

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/dphaener/rhiza/internal/adapter"
)

// Name is the registry name of this adapter.
const Name = "claude-code"

// hookEvents are the events whose hooks the settings JSON installs.
var hookEvents = []string{"SessionStart", "UserPromptSubmit", "PreToolUse", "PostToolUse", "Notification", "Stop", "SessionEnd"}

var modes = []string{"acceptEdits", "auto", "bypassPermissions", "manual", "dontAsk", "plan"}

func init() { adapter.RegisterDefault(Adapter{}) }

// Adapter launches Claude Code.
type Adapter struct{}

// Name implements adapter.Adapter.
func (Adapter) Name() string { return Name }

// PermissionModes implements adapter.Adapter.
func (Adapter) PermissionModes() []string { return append([]string(nil), modes...) }

// Bypass implements adapter.Adapter.
func (Adapter) Bypass(mode string) bool { return mode == "bypassPermissions" }

// Command implements adapter.Adapter. The permission mode is always an
// argument and never part of the settings JSON. A model is passed as
// --model before the prompt; none leaves the flag out.
func (a Adapter) Command(l adapter.Launch) ([]string, error) {
	if err := adapter.CheckMode(a, l.PermissionMode); err != nil {
		return nil, err
	}
	if l.SessionID == "" {
		return nil, fmt.Errorf("a session ID is required")
	}
	settings, err := Settings(l.HookCommand)
	if err != nil {
		return nil, err
	}
	idFlag := "--session-id"
	if l.Resume {
		idFlag = "--resume"
	}
	argv := []string{"claude", idFlag, l.SessionID, "--permission-mode", l.PermissionMode, "--settings", settings}
	if l.Model != "" {
		argv = append(argv, "--model", l.Model)
	}
	if l.Prompt != "" {
		argv = append(argv, l.Prompt)
	}
	return argv, nil
}

type hookCmd struct {
	Type    string `json:"type"`
	Command string `json:"command"`
}

type hookGroup struct {
	Hooks []hookCmd `json:"hooks"`
}

type settings struct {
	Hooks map[string][]hookGroup `json:"hooks"`
}

// Settings returns the settings JSON: command hooks only, never permissions.
func Settings(hookCommand []string) (string, error) {
	if len(hookCommand) == 0 {
		return "", fmt.Errorf("a hook command is required")
	}
	s := settings{Hooks: map[string][]hookGroup{}}
	for _, ev := range hookEvents {
		cmd := shellJoin(append(append([]string(nil), hookCommand...), ev))
		s.Hooks[ev] = []hookGroup{{Hooks: []hookCmd{{Type: "command", Command: cmd}}}}
	}
	b, err := json.Marshal(s)
	return string(b), err
}

func shellJoin(args []string) string {
	out := make([]string, len(args))
	for i, a := range args {
		out[i] = shellQuote(a)
	}
	return strings.Join(out, " ")
}

func shellQuote(s string) string {
	if s != "" && strings.IndexFunc(s, func(r rune) bool {
		return !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || strings.ContainsRune("-_./:=@%+,", r))
	}) < 0 {
		return s
	}
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}

type payload struct {
	SessionID        string `json:"session_id"`
	NotificationType string `json:"notification_type"`
	Message          string `json:"message"`
	ToolName         string `json:"tool_name"`
	ToolInput        struct {
		Questions []struct {
			Question string `json:"question"`
		} `json:"questions"`
	} `json:"tool_input"`
}

// askTool is the tool the agent calls to put a question to the owner.
const askTool = "AskUserQuestion"

// sendTool is the tool the agent calls to send another session a message.
const sendTool = "SendMessage"

// questionText joins the questions of an ask call with " / ". Each is
// trimmed and has its control characters, newlines included, replaced by a
// space; empty questions are dropped.
func (p payload) questionText() string {
	var parts []string
	for _, q := range p.ToolInput.Questions {
		t := strings.TrimSpace(strings.Map(func(r rune) rune {
			if unicode.IsControl(r) {
				return ' '
			}
			return r
		}, q.Question))
		if t != "" {
			parts = append(parts, t)
		}
	}
	return strings.Join(parts, " / ")
}

// HookStatus implements adapter.Adapter.
func (Adapter) HookStatus(event string, raw []byte) (adapter.HookReport, bool, error) {
	var p payload
	if len(raw) > 0 {
		if err := json.Unmarshal(raw, &p); err != nil {
			return adapter.HookReport{}, false, fmt.Errorf("%s hook payload: %w", event, err)
		}
	}
	r := adapter.HookReport{SessionID: p.SessionID}
	switch event {
	case "SessionStart", "Stop":
		r.State = "idle"
	case "PreToolUse":
		if p.ToolName == askTool {
			r.State, r.Kind, r.Detail = "blocked", "question", p.questionText()
		} else {
			r.State = "working"
		}
	case "UserPromptSubmit", "PostToolUse":
		r.State = "working"
	case "SessionEnd":
		r.State = "done"
	case "Notification":
		switch p.NotificationType {
		case "permission_prompt":
			r.State, r.Kind, r.Detail = "blocked", "permission", p.Message
		case "idle_prompt":
			r.State = "idle"
		default:
			return r, false, nil
		}
	default:
		return adapter.HookReport{}, false, fmt.Errorf("unknown hook event %q", event)
	}
	return r, true, nil
}

// AskReply implements adapter.Adapter. A PreToolUse on the ask tool is the
// agent asking the owner; the reply denies that call with reason as the
// PreToolUse hook's permission decision.
func (Adapter) AskReply(event string, raw []byte, reason string) ([]byte, bool, error) {
	if event != "PreToolUse" {
		return nil, false, nil
	}
	var p payload
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, false, fmt.Errorf("%s hook payload: %w", event, err)
	}
	if p.ToolName != askTool {
		return nil, false, nil
	}
	type hookOutput struct {
		HookEventName            string `json:"hookEventName"`
		PermissionDecision       string `json:"permissionDecision"`
		PermissionDecisionReason string `json:"permissionDecisionReason"`
	}
	reply, err := json.Marshal(struct {
		HookSpecificOutput hookOutput `json:"hookSpecificOutput"`
	}{hookOutput{event, "deny", reason}})
	if err != nil {
		return nil, false, err
	}
	return reply, true, nil
}

// ContextReply implements adapter.Adapter. A PostToolUse hook adds its
// additionalContext to the agent's context at its next step.
func (Adapter) ContextReply(event, context string) ([]byte, bool, error) {
	if event != "PostToolUse" {
		return nil, false, nil
	}
	type hookOutput struct {
		HookEventName     string `json:"hookEventName"`
		AdditionalContext string `json:"additionalContext"`
	}
	reply, err := json.Marshal(struct {
		HookSpecificOutput hookOutput `json:"hookSpecificOutput"`
	}{hookOutput{event, context}})
	if err != nil {
		return nil, false, err
	}
	return reply, true, nil
}

// sendPayload is the part of a hook payload a send call's message is read
// from. ToolInput stays raw until the tool is known to be the send tool.
type sendPayload struct {
	SessionID string          `json:"session_id"`
	ToolName  string          `json:"tool_name"`
	ToolInput json.RawMessage `json:"tool_input"`
}

// sendInput is a send call's input. To is a pointer and Message raw so that
// absent, null and empty can be told apart.
type sendInput struct {
	To      *string         `json:"to"`
	Summary string          `json:"summary"`
	Message json.RawMessage `json:"message"`
}

// HookMessage implements adapter.Adapter. Only a PostToolUse event for the
// send tool carries a message, so the message is read once the call has run.
// A string message is its text, verbatim; any other JSON value is kept as
// its raw JSON and marked Structured.
func (Adapter) HookMessage(event string, raw []byte) (adapter.Message, bool, error) {
	var p sendPayload
	if len(raw) > 0 {
		if err := json.Unmarshal(raw, &p); err != nil {
			return adapter.Message{}, false, fmt.Errorf("%s hook payload: %w", event, err)
		}
	}
	m := adapter.Message{SessionID: p.SessionID}
	if event != "PostToolUse" || p.ToolName != sendTool {
		return m, false, nil
	}
	var in sendInput
	if len(p.ToolInput) > 0 {
		if err := json.Unmarshal(p.ToolInput, &in); err != nil {
			return adapter.Message{}, false, fmt.Errorf("%s hook payload: %s input: %w", event, sendTool, err)
		}
	}
	if in.To == nil || *in.To == "" {
		return adapter.Message{}, false, fmt.Errorf("%s hook payload: %s call has no recipient", event, sendTool)
	}
	msg := bytes.TrimSpace(in.Message)
	if len(msg) == 0 || string(msg) == "null" {
		return adapter.Message{}, false, fmt.Errorf("%s hook payload: %s call has no message", event, sendTool)
	}
	m.To, m.Summary = *in.To, in.Summary
	if msg[0] == '"' {
		if err := json.Unmarshal(msg, &m.Text); err != nil {
			return adapter.Message{}, false, fmt.Errorf("%s hook payload: %w", event, err)
		}
	} else {
		m.Text, m.Structured = string(msg), true
	}
	return m, true, nil
}

// TranscriptGlobs implements adapter.Adapter. Claude Code keeps a session's
// transcript at <config>/projects/<encoded cwd>/<session-id>.jsonl, where
// <config> is $CLAUDE_CONFIG_DIR when set and $HOME/.claude otherwise. Helper
// agents' transcripts sit one level deeper, under <session-id>/subagents/,
// so the patterns do not match them.
func (Adapter) TranscriptGlobs(sessionID string, getenv func(string) string) []string {
	if !safeSessionID(sessionID) || getenv == nil {
		return nil
	}
	var globs []string
	add := func(config string) {
		globs = append(globs, filepath.Join(globEscape(config), "projects", "*", sessionID+".jsonl"))
	}
	if dir := getenv("CLAUDE_CONFIG_DIR"); dir != "" {
		add(dir)
	}
	if home := getenv("HOME"); home != "" {
		add(filepath.Join(home, ".claude"))
	}
	return globs
}

// safeSessionID reports whether id can stand as one literal path element in
// a glob pattern: not empty, no separator, no parent reference and no
// pattern syntax.
func safeSessionID(id string) bool {
	return id != "" && !strings.ContainsAny(id, `/\*?[`) && !strings.Contains(id, "..")
}

// globEscape quotes the pattern syntax in a literal directory so that
// filepath.Glob matches it as written. Windows has no escape character in
// patterns, so there the directory is used as is.
func globEscape(dir string) string {
	if runtime.GOOS == "windows" {
		return dir
	}
	var b strings.Builder
	for _, r := range dir {
		if strings.ContainsRune(`*?[\`, r) {
			b.WriteByte('\\')
		}
		b.WriteRune(r)
	}
	return b.String()
}

// transcriptLine is the part of one transcript line that Answers reads.
type transcriptLine struct {
	Message struct {
		Content json.RawMessage `json:"content"`
	} `json:"message"`
	ToolUseResult json.RawMessage `json:"toolUseResult"`
}

type contentBlock struct {
	Type      string `json:"type"`
	ToolUseID string `json:"tool_use_id"`
}

// askResult is the toolUseResult of an answered AskUserQuestion call. A
// rejected call's toolUseResult is a string instead.
type askResult struct {
	Questions []struct {
		Question string `json:"question"`
		Header   string `json:"header"`
		Options  []struct {
			Label string `json:"label"`
		} `json:"options"`
		MultiSelect bool `json:"multiSelect"`
	} `json:"questions"`
	Answers     map[string]string `json:"answers"`
	Annotations map[string]struct {
		Notes string `json:"notes"`
	} `json:"annotations"`
}

// Answers implements adapter.Adapter. It reads the transcript one JSON line
// at a time with no line-length limit, and skips every line that is not an
// answered AskUserQuestion result, malformed lines included.
func (Adapter) Answers(transcript []byte) ([]adapter.Answer, error) {
	var out []adapter.Answer
	r := bufio.NewReader(bytes.NewReader(transcript))
	for {
		line, err := r.ReadBytes('\n')
		out = append(out, lineAnswers(line)...)
		if errors.Is(err, io.EOF) {
			return out, nil
		}
		if err != nil {
			return out, err
		}
	}
}

// lineAnswers returns the answered questions on one transcript line, or nil
// when the line holds none.
func lineAnswers(line []byte) []adapter.Answer {
	line = bytes.TrimSpace(line)
	if len(line) == 0 || line[0] != '{' {
		return nil
	}
	var l transcriptLine
	if json.Unmarshal(line, &l) != nil {
		return nil
	}
	res := bytes.TrimSpace(l.ToolUseResult)
	if len(res) == 0 || res[0] != '{' {
		return nil
	}
	var ask askResult
	if json.Unmarshal(res, &ask) != nil || len(ask.Answers) == 0 {
		return nil
	}
	id := toolUseID(l.Message.Content)
	var out []adapter.Answer
	for _, q := range ask.Questions {
		answer, ok := ask.Answers[q.Question]
		if !ok || answer == "" {
			continue
		}
		labels := make([]string, len(q.Options))
		for i, o := range q.Options {
			labels[i] = o.Label
		}
		a := adapter.Answer{
			ToolUseID:   id,
			Question:    q.Question,
			Header:      q.Header,
			Options:     labels,
			MultiSelect: q.MultiSelect,
			Notes:       ask.Annotations[q.Question].Notes,
		}
		a.Picked, a.Typed = splitAnswer(answer, labels, q.MultiSelect)
		out = append(out, a)
	}
	return out
}

// toolUseID returns the tool_use_id of the first tool_result block in a
// message's content, or "" when there is none.
func toolUseID(content json.RawMessage) string {
	var blocks []contentBlock
	if json.Unmarshal(content, &blocks) != nil {
		return ""
	}
	for _, b := range blocks {
		if b.Type == "tool_result" && b.ToolUseID != "" {
			return b.ToolUseID
		}
	}
	return ""
}

// limitLine is the part of one transcript line that UsageLimit reads.
// Claude Code writes the usage-limit message as a synthetic API-error line:
// isApiErrorMessage set, model "<synthetic>", the message in a text block.
type limitLine struct {
	Type       string `json:"type"`
	IsAPIError bool   `json:"isApiErrorMessage"`
	Message    struct {
		Model   string          `json:"model"`
		Content json.RawMessage `json:"content"`
	} `json:"message"`
}

// apiError reports whether the line is Claude Code's synthetic API-error
// message rather than the model's own turn.
func (l limitLine) apiError() bool {
	return l.IsAPIError || l.Message.Model == "<synthetic>"
}

// text returns the line's text blocks joined, or its content when that is a
// plain string. Tool calls and their inputs are left out.
func (l limitLine) text() string {
	var s string
	if json.Unmarshal(l.Message.Content, &s) == nil {
		return s
	}
	var blocks []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	}
	if json.Unmarshal(l.Message.Content, &blocks) != nil {
		return ""
	}
	var parts []string
	for _, b := range blocks {
		if b.Type == "text" {
			parts = append(parts, b.Text)
		}
	}
	return strings.Join(parts, "\n")
}

var limitReset = regexp.MustCompile(`(?i)usage limit reached\|(\d{9,})`)

// UsageLimit implements adapter.Adapter. The session is limited when the
// last assistant line of its transcript is the synthetic API-error message
// saying so in its text: "usage limit reached", with the reset as a Unix
// time after a bar, or "hit your limit". The phrases anywhere else, in the
// model's prose or a tool call's input, do not count. Any later assistant
// line means the session went on.
func (Adapter) UsageLimit(transcript []byte) (adapter.UsageLimit, bool) {
	var last limitLine
	for _, line := range bytes.Split(transcript, []byte("\n")) {
		line = bytes.TrimSpace(line)
		if len(line) == 0 || line[0] != '{' {
			continue
		}
		var l limitLine
		if json.Unmarshal(line, &l) != nil || l.Type != "assistant" {
			continue
		}
		last = l
	}
	if !last.apiError() {
		return adapter.UsageLimit{}, false
	}
	text := last.text()
	lower := strings.ToLower(text)
	if !strings.Contains(lower, "usage limit reached") && !strings.Contains(lower, "hit your limit") {
		return adapter.UsageLimit{}, false
	}
	var out adapter.UsageLimit
	if m := limitReset.FindStringSubmatch(text); m != nil {
		if sec, err := strconv.ParseInt(m[1], 10, 64); err == nil {
			out.ResetAt = time.Unix(sec, 0).UTC()
		}
	}
	return out, true
}

// PaneUsageLimit implements adapter.PaneLimiter: the pane's last lines show
// the usage-limit message the agent prints, "usage limit reached" or "hit
// your limit". Only the tail counts, so text scrolled far above does not.
func (Adapter) PaneUsageLimit(pane string) bool {
	lines := strings.Split(strings.TrimRight(pane, "\n \t"), "\n")
	tail := strings.ToLower(strings.Join(lines[max(0, len(lines)-15):], "\n"))
	return strings.Contains(tail, "usage limit reached") || strings.Contains(tail, "hit your limit")
}

// usageLine is the part of one transcript line that TokenUse reads. Usage
// is a pointer so that a line without usage can be told apart from one with
// zero tokens.
type usageLine struct {
	Type        string `json:"type"`
	IsSidechain bool   `json:"isSidechain"`
	Message     struct {
		ID    string `json:"id"`
		Usage *struct {
			Input         int64 `json:"input_tokens"`
			Output        int64 `json:"output_tokens"`
			CacheCreation int64 `json:"cache_creation_input_tokens"`
			CacheRead     int64 `json:"cache_read_input_tokens"`
		} `json:"usage"`
	} `json:"message"`
}

// TokenUse implements adapter.Adapter. It reads the transcript one JSON line
// at a time with no line-length limit, like Answers. Only assistant lines
// with usage count. Claude Code writes one line per content block of a
// message, each repeating the message's id and usage, so a message counts
// once per id; a line with no id counts on its own. Sidechain lines are the
// helper agents' and count toward Helpers, deduplicated apart from the
// session's.
func (Adapter) TokenUse(transcript []byte) (adapter.TokenUse, error) {
	var use adapter.TokenUse
	seen := [2]map[string]bool{{}, {}}
	r := bufio.NewReader(bytes.NewReader(transcript))
	for {
		line, err := r.ReadBytes('\n')
		addLineUsage(&use, &seen, line)
		if errors.Is(err, io.EOF) {
			return use, nil
		}
		if err != nil {
			return use, err
		}
	}
}

// addLineUsage adds one transcript line's usage to use, unless the line has
// none or its message was already counted in its bucket. seen holds the
// counted message ids, the session's first and the helpers' second.
func addLineUsage(use *adapter.TokenUse, seen *[2]map[string]bool, line []byte) {
	line = bytes.TrimSpace(line)
	if len(line) == 0 || line[0] != '{' {
		return
	}
	var l usageLine
	if json.Unmarshal(line, &l) != nil || l.Type != "assistant" || l.Message.Usage == nil {
		return
	}
	bucket, ids := &use.Session, seen[0]
	if l.IsSidechain {
		bucket, ids = &use.Helpers, seen[1]
	}
	if id := l.Message.ID; id != "" {
		if ids[id] {
			return
		}
		ids[id] = true
	}
	u := l.Message.Usage
	bucket.Input += u.Input
	bucket.Output += u.Output
	bucket.CacheCreation += u.CacheCreation
	bucket.CacheRead += u.CacheRead
}

// answerSep joins the labels of a multiSelect answer.
const answerSep = ", "

// splitAnswer separates an answer into the option labels the owner picked
// and the text the owner typed. A single-select answer is either exactly a
// label or typed text. A multiSelect answer is labels and typed text joined
// with ", "; since labels may contain ", " themselves, labels are matched
// whole at ", " boundaries, longest first, and never split. What no label
// claims is the typed text, each leftover piece trimmed and the pieces
// rejoined with ", ".
func splitAnswer(answer string, labels []string, multi bool) (picked []string, typed string) {
	for _, l := range labels {
		if l != "" && answer == l {
			return []string{l}, ""
		}
	}
	if !multi {
		return nil, answer
	}
	order := make([]string, 0, len(labels))
	for _, l := range labels {
		if l != "" {
			order = append(order, l)
		}
	}
	sort.SliceStable(order, func(i, j int) bool { return len(order[i]) > len(order[j]) })

	type span struct {
		start, end int
		label      string
	}
	var claimed []span
	free := func(start, end int) bool {
		for _, c := range claimed {
			if start < c.end && c.start < end {
				return false
			}
		}
		return true
	}
	atBoundary := func(start, end int) bool {
		return (start == 0 || strings.HasSuffix(answer[:start], answerSep)) &&
			(end == len(answer) || strings.HasPrefix(answer[end:], answerSep))
	}
	for _, l := range order {
		for from := 0; from <= len(answer)-len(l); {
			i := strings.Index(answer[from:], l)
			if i < 0 {
				break
			}
			start, end := from+i, from+i+len(l)
			if atBoundary(start, end) && free(start, end) {
				claimed = append(claimed, span{start, end, l})
				break
			}
			from = start + 1
		}
	}
	sort.Slice(claimed, func(i, j int) bool { return claimed[i].start < claimed[j].start })

	var rest []string
	keep := func(start, end int) {
		piece := answer[start:end]
		if start > 0 {
			piece = strings.TrimPrefix(piece, answerSep)
		}
		if end < len(answer) {
			piece = strings.TrimSuffix(piece, answerSep)
		}
		if piece = strings.TrimSpace(piece); piece != "" {
			rest = append(rest, piece)
		}
	}
	prev := 0
	for _, c := range claimed {
		picked = append(picked, c.label)
		keep(prev, c.start)
		prev = c.end
	}
	keep(prev, len(answer))
	return picked, strings.Join(rest, answerSep)
}

// TrustFile implements adapter.Truster. Claude Code keeps its global state
// at $CLAUDE_CONFIG_DIR/.claude.json when that is set and $HOME/.claude.json
// otherwise.
func (Adapter) TrustFile(getenv func(string) string) string {
	if getenv == nil {
		return ""
	}
	if dir := getenv("CLAUDE_CONFIG_DIR"); dir != "" {
		return filepath.Join(dir, ".claude.json")
	}
	if home := getenv("HOME"); home != "" {
		return filepath.Join(home, ".claude.json")
	}
	return ""
}

// Trust implements adapter.Truster: it sets
// projects[<dir>].hasTrustDialogAccepted to true for each dir and keeps
// every other key, so Claude Code never stops at its folder-trust prompt.
func (Adapter) Trust(current []byte, dirs []string) ([]byte, error) {
	top := map[string]json.RawMessage{}
	if len(bytes.TrimSpace(current)) > 0 {
		if err := json.Unmarshal(current, &top); err != nil {
			return nil, fmt.Errorf("read the Claude Code config: %w", err)
		}
	}
	projects := map[string]json.RawMessage{}
	if raw, ok := top["projects"]; ok && string(bytes.TrimSpace(raw)) != "null" {
		if err := json.Unmarshal(raw, &projects); err != nil {
			return nil, fmt.Errorf("read projects of the Claude Code config: %w", err)
		}
	}
	for _, dir := range dirs {
		if dir == "" {
			continue
		}
		proj := map[string]json.RawMessage{}
		if raw, ok := projects[dir]; ok && string(bytes.TrimSpace(raw)) != "null" {
			if err := json.Unmarshal(raw, &proj); err != nil {
				return nil, fmt.Errorf("read project %s of the Claude Code config: %w", dir, err)
			}
		}
		proj["hasTrustDialogAccepted"] = json.RawMessage("true")
		raw, err := json.Marshal(proj)
		if err != nil {
			return nil, err
		}
		projects[dir] = raw
	}
	raw, err := json.Marshal(projects)
	if err != nil {
		return nil, err
	}
	top["projects"] = raw
	out, err := json.MarshalIndent(top, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(out, '\n'), nil
}

// configDir returns the agent's config directory: $CLAUDE_CONFIG_DIR when
// set, else $HOME/.claude; "" when neither is known.
func configDir(getenv func(string) string) string {
	if getenv == nil {
		return ""
	}
	if dir := getenv("CLAUDE_CONFIG_DIR"); dir != "" {
		return dir
	}
	if home := getenv("HOME"); home != "" {
		return filepath.Join(home, ".claude")
	}
	return ""
}

// stampedUsage is a usageLine with its timestamp.
type stampedUsage struct {
	usageLine
	Timestamp string `json:"timestamp"`
}

// UsageRecords implements adapter.UsageReader. It walks every .jsonl file
// under <config>/projects, which holds sessions and, deeper, their helpers'
// transcripts, and reads nothing outside it. Input, output and
// cache-creation tokens count; cache reads never do. A message repeated
// across lines or files counts once, at its earliest timestamp, and a line
// with no id counts on its own.
func (Adapter) UsageRecords(getenv func(string) string, from, to time.Time) ([]adapter.UsageRecord, error) {
	config := configDir(getenv)
	if config == "" {
		return nil, nil
	}
	root := filepath.Join(config, "projects")
	if _, err := os.Stat(root); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	var recs []adapter.UsageRecord
	byID := map[string]int{}
	_ = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || !d.Type().IsRegular() || filepath.Ext(path) != ".jsonl" {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		for _, line := range bytes.Split(data, []byte("\n")) {
			line = bytes.TrimSpace(line)
			if len(line) == 0 || line[0] != '{' {
				continue
			}
			var l stampedUsage
			if json.Unmarshal(line, &l) != nil || l.Type != "assistant" || l.Message.Usage == nil {
				continue
			}
			at, err := time.Parse(time.RFC3339Nano, l.Timestamp)
			if err != nil {
				continue
			}
			u := l.Message.Usage
			rec := adapter.UsageRecord{At: at, Tokens: u.Input + u.Output + u.CacheCreation}
			if id := l.Message.ID; id != "" {
				if i, ok := byID[id]; ok {
					if at.Before(recs[i].At) {
						recs[i].At = at
					}
					continue
				}
				byID[id] = len(recs)
			}
			recs = append(recs, rec)
		}
		return nil
	})
	out := recs[:0]
	for _, r := range recs {
		if !r.At.Before(from) && r.At.Before(to) {
			out = append(out, r)
		}
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].At.Before(out[j].At) })
	return out, nil
}
