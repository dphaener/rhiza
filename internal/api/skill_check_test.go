package api

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/check"
	"github.com/dphaener/rhiza/internal/workflow"
)

func TestSkillCheckReport(t *testing.T) {
	gates := func(results ...string) []SkillGate {
		names := []string{gateSessionDone, gateTrailers, gateCheck}
		var out []SkillGate
		for i, r := range results {
			out = append(out, SkillGate{Name: names[i], Result: r, Detail: names[i] + " detail"})
		}
		return out
	}
	cases := []struct {
		name   string
		resp   SkillCheckResponse
		exit   int
		header string
		lines  []string
	}{
		{
			name:   "AllPass",
			resp:   SkillCheckResponse{Invoke: "/b {leaf}", Role: "builder", Agent: "claude", Passed: true, Gates: gates("pass", "pass", "pass"), Run: "succeeded"},
			exit:   0,
			header: "skill check /b {leaf} as builder passes: 3 of 3 gates pass",
			lines:  []string{"  pass session-done session-done detail", "  pass trailers     trailers detail", "  pass check        check detail", "  run: succeeded"},
		},
		{
			name:   "TrailersFail",
			resp:   SkillCheckResponse{Invoke: "/b", Role: "builder", Agent: "claude", Gates: gates("pass", "fail", "pass"), Run: "failed: leaf"},
			exit:   1,
			header: "skill check /b as builder fails: 2 of 3 gates pass",
			lines:  []string{"  fail trailers     trailers detail"},
		},
		{
			name:   "Kept",
			resp:   SkillCheckResponse{Invoke: "/b", Role: "builder", Gates: gates("fail", "pass", "pass"), Scratch: "/tmp/x/repo"},
			exit:   1,
			header: "skill check /b as builder fails: 2 of 3 gates pass",
			lines:  []string{"  fail session-done session-done detail", "  scratch repo kept: /tmp/x/repo"},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.resp.ExitCode(); got != c.exit {
				t.Errorf("ExitCode() = %d, want %d", got, c.exit)
			}
			text := human(t, c.resp)
			lines := strings.Split(strings.TrimSuffix(text, "\n"), "\n")
			if lines[0] != c.header {
				t.Errorf("header = %q, want %q", lines[0], c.header)
			}
			for _, want := range c.lines {
				if !strings.Contains(text, want+"\n") {
					t.Errorf("text lacks line %q:\n%s", want, text)
				}
			}
			b, err := json.Marshal(c.resp)
			if err != nil {
				t.Fatal(err)
			}
			var back struct {
				Passed bool `json:"passed"`
				Gates  []struct {
					Name, Result string
				} `json:"gates"`
				Scratch *string `json:"scratch"`
			}
			if err := json.Unmarshal(b, &back); err != nil {
				t.Fatal(err)
			}
			if back.Passed != c.resp.Passed || len(back.Gates) != len(c.resp.Gates) {
				t.Fatalf("json = %s", b)
			}
			for i, g := range back.Gates {
				if g.Name != c.resp.Gates[i].Name || g.Result != c.resp.Gates[i].Result {
					t.Errorf("json gate %d = %+v, want %+v", i, g, c.resp.Gates[i])
				}
			}
			if (back.Scratch != nil) != (c.resp.Scratch != "") {
				t.Errorf("json scratch = %v, want it only when kept: %s", back.Scratch, b)
			}
		})
	}
}

func TestSkillCheckTrailersGate(t *testing.T) {
	const change = "01CHANGE"
	full := check.Commit{Hash: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", Changes: []string{change}, Actors: []string{"claude-code-1"}, Leaves: []string{skillCheckLeaf}}
	cases := []struct {
		name    string
		commits []check.Commit
		why     string
		want    string
		detail  string
	}{
		{name: "None", why: "the builder made no commit", want: "fail", detail: "the builder made no commit"},
		{name: "All three", commits: []check.Commit{full}, want: "pass", detail: "commit aaaaaaaaaaaa carries"},
		{name: "No trailers", commits: []check.Commit{{Hash: "b"}}, want: "fail", detail: "none of the builder's 1 commit(s) carries all of Change: 01CHANGE, Actor: and Leaf: greeting-file"},
		{name: "No actor", commits: []check.Commit{{Hash: "c", Changes: []string{change}, Leaves: []string{skillCheckLeaf}}}, want: "fail", detail: "none of the builder's"},
		{name: "Other change", commits: []check.Commit{{Hash: "d", Changes: []string{"01OTHER"}, Actors: []string{"a"}, Leaves: []string{skillCheckLeaf}}}, want: "fail", detail: "none of the builder's"},
		{name: "Other leaf", commits: []check.Commit{{Hash: "e", Changes: []string{change}, Actors: []string{"a"}, Leaves: []string{"other"}}}, want: "fail", detail: "none of the builder's"},
		{name: "One of two", commits: []check.Commit{{Hash: "f"}, full}, want: "pass", detail: "commit aaaaaaaaaaaa carries"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			g := trailersGate(c.commits, c.why, change)
			if g.Name != gateTrailers || g.Result != c.want || !strings.Contains(g.Detail, c.detail) {
				t.Errorf("trailersGate = %+v, want %s with %q", g, c.want, c.detail)
			}
		})
	}
}

func TestSkillCheckDefinition(t *testing.T) {
	invoke := `/my-builder "quoted" \ {leaf}`
	data, err := skillCheckDefinition(invoke)
	if err != nil {
		t.Fatal(err)
	}
	name, problems := workflow.Validate(skillCheckWorkflow, data, nil)
	if len(problems) > 0 {
		t.Fatalf("fixture workflow has problems %v:\n%s", problems, data)
	}
	def, err := workflow.ParseWith(name, data, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(def.Order) != 1 {
		t.Fatalf("steps = %v, want one", def.Order)
	}
	s := def.Steps[def.Order[0]]
	if s.Kind != workflow.KindLeaves {
		t.Fatalf("kind = %q, want leaves", s.Kind)
	}
	if got := s.Role(workflow.RoleBuilder).RenderPrompt(workflow.LeafVars{Leaf: skillCheckLeaf}); got != `/my-builder "quoted" \ greeting-file` {
		t.Errorf("builder prompt = %q", got)
	}
}

func TestSkillCheckRefusals(t *testing.T) {
	home := t.TempDir()
	ctx := testCtx(t.TempDir())
	ctx.Home = home
	ctx.Getenv = func(k string) string {
		if k == "HOME" {
			return home
		}
		return ""
	}
	cases := []struct {
		name, req, want string
	}{
		{"No invoke", `{"permission_mode":"manual"}`, "skill check needs an invoke"},
		{"Unknown role", `{"invoke":"Build {leaf}","role":"planner","permission_mode":"manual"}`, `unknown role "planner": want builder`},
		{"No permission mode", `{"invoke":"Build {leaf}"}`, "needs --permission-mode"},
		{"Skill missing", `{"invoke":"/rhiza-builder Build {leaf}","permission_mode":"manual"}`, "start command /rhiza-builder unknown to claude: skill not installed; run rhiza guide install --agent claude"},
		{"No programs", `{"invoke":"Build {leaf}","permission_mode":"manual"}`, "this caller runs no programs"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := callOp(t, ctx, "skill.check", c.req)
			if err == nil || !strings.Contains(AsError(err).Message, c.want) {
				t.Fatalf("err = %v, want %q", err, c.want)
			}
		})
	}
}
