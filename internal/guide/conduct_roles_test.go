package guide

import (
	"slices"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/workflow"
)

// conductRoles are the roles a conductor spawns, in the topic's order.
var conductRoles = []string{
	"builder", "verifier", "check-plan", "check-code", "check-qa",
	"critic", "fixer", "gate-fix", "ci-fix", "conflict-fix",
}

// listedRoles returns the roles the Roles section of body lists, one per
// item of the form "- `<role>` (`rhiza guide <topic>`)", and the items
// whose topic is not the role's own name.
func listedRoles(body string) (roles []string, mismatched []string) {
	for _, l := range strings.Split(sectionOf(body, "## Roles"), "\n") {
		rest, ok := strings.CutPrefix(l, "- `")
		if !ok {
			continue
		}
		role, rest, ok := strings.Cut(rest, "`")
		if !ok {
			continue
		}
		roles = append(roles, role)
		if !strings.HasPrefix(rest, " (`rhiza guide "+role+"`)") {
			mismatched = append(mismatched, role)
		}
	}
	return roles, mismatched
}

// roleLine returns the Roles item of body for role with its indented
// continuation lines.
func roleLine(body, role string) string {
	out := ""
	for _, l := range strings.Split(sectionOf(body, "## Roles"), "\n") {
		switch {
		case strings.HasPrefix(l, "- `"+role+"`"):
			out = l
		case out != "" && strings.HasPrefix(l, "  "):
			out += "\n" + l
		case out != "":
			return out
		}
	}
	return out
}

// rolesWithoutTopic returns each role body lists that no guide topic is
// named for.
func rolesWithoutTopic(body string) []string {
	roles, _ := listedRoles(body)
	var out []string
	for _, r := range roles {
		if tx, ok := Lookup(r); !ok || tx.Skill {
			out = append(out, r)
		}
	}
	return out
}

// checksWithoutRole returns each check the conductor steps of def declare
// that body lists no check-<check> role for.
func checksWithoutRole(body string, def *workflow.Definition) []string {
	roles, _ := listedRoles(body)
	var out []string
	for _, s := range def.Steps {
		if !s.Conducts() {
			continue
		}
		for _, c := range s.Checks {
			if !slices.Contains(roles, "check-"+c.Name) && !slices.Contains(out, c.Name) {
				out = append(out, c.Name)
			}
		}
	}
	slices.Sort(out)
	return out
}

// shipChange parses the built-in ship-change workflow.
func shipChange(t *testing.T) *workflow.Definition {
	t.Helper()
	data, ok := workflow.Builtin("ship-change")
	if !ok {
		t.Fatal("no built-in ship-change")
	}
	def, err := workflow.Parse("ship-change", data)
	if err != nil {
		t.Fatal(err)
	}
	return def
}

func TestConductRoles(t *testing.T) {
	c, ok := Lookup("conduct-roles")
	if !ok || c.Skill {
		t.Fatalf("conduct-roles: found %v, skill %v; want a topic", ok, c.Skill)
	}
	if n := Tokens(c.Body); n > TopicCap {
		t.Errorf("conduct-roles costs %d tokens, cap %d", n, TopicCap)
	}
	roles, mismatched := listedRoles(c.Body)
	if !slices.Equal(roles, conductRoles) {
		t.Errorf("conduct-roles lists roles %q, want %q", roles, conductRoles)
	}
	for _, r := range mismatched {
		t.Errorf("conduct-roles: role %s does not name rhiza guide %s", r, r)
	}
	for _, r := range rolesWithoutTopic(c.Body) {
		t.Errorf("conduct-roles: role %s has no guide topic", r)
	}
	if !strings.Contains(sectionOf(c.Body, "## Roles"), "the skill of the same name, `/rhiza-<role>`") {
		t.Errorf("conduct-roles: Roles does not name each role's skill")
	}
	def := shipChange(t)
	if s, ok := def.Steps["conduct"]; !ok || !s.Conducts() {
		t.Fatalf("ship-change's conduct step is no conductor step")
	}
	for _, n := range checksWithoutRole(c.Body, def) {
		t.Errorf("conduct-roles: ship-change declares check %s, but no check-%s role is listed", n, n)
	}
	checkPhrases(t, "conduct-roles", map[string][]string{
		"flags": {
			"`--brief <file>` always", "`--leaf <leaf>` for a builder and a verifier",
			"`--check <check>` for a critic", "`--model <model>` optional",
		},
		"check roles": {
			"`[[step.<name>.check]]`", "conductor step", "the arbiter",
			"`RHIZA_CHECK`", "`RHIZA_CYCLE`",
			"refused with the declared names", "name one of its checks: plan, code, qa",
			"`compliance`, `quality` and `model` belong to the classic workflows",
		},
	})
	if strings.Contains(c.Body, "ship rebase --no-push") {
		t.Errorf("conduct-roles still names rhiza ship rebase --no-push")
	}
	if l := roleLine(c.Body, "conflict-fix"); !strings.Contains(l, "a code conflict the ship steps leave in\n  progress") {
		t.Errorf("conduct-roles: conflict-fix entry %q names no code conflict the ship steps leave in progress", l)
	}
	if strings.Contains(c.Body+c.Description, "Claude") {
		t.Errorf("conduct-roles names an agent product")
	}
}

// The conduct topic points to conduct-roles in Read first and stays
// under its cap.
func TestConductReadsRoles(t *testing.T) {
	c, _ := Lookup("conduct")
	if !strings.Contains(sectionOf(c.Body, "## Read first"), "`rhiza guide conduct-roles`") {
		t.Errorf("conduct: Read first does not name rhiza guide conduct-roles")
	}
	if n := Tokens(c.Body); n > TopicCap {
		t.Errorf("conduct costs %d tokens, cap %d", n, TopicCap)
	}
}

// A role with no guide topic, and a declared check with no check-<check>
// role, are each reported by name.
func TestConductRolesCatchesDrift(t *testing.T) {
	c, _ := Lookup("conduct-roles")
	broken := strings.Replace(c.Body, "- `ci-fix` (`rhiza guide ci-fix`)",
		"- `check-scope` (`rhiza guide check-scope`): never.\n- `ci-fix` (`rhiza guide ci-fix`)", 1)
	if got := rolesWithoutTopic(broken); !slices.Equal(got, []string{"check-scope"}) {
		t.Errorf("rolesWithoutTopic = %q, want [check-scope]", got)
	}
	def := shipChange(t)
	conduct := def.Steps["conduct"]
	conduct.Checks = append(slices.Clone(conduct.Checks), workflow.Check{Name: "lint"})
	def.Steps["conduct"] = conduct
	if got := checksWithoutRole(c.Body, def); !slices.Equal(got, []string{"lint"}) {
		t.Errorf("checksWithoutRole = %q, want [lint]", got)
	}
	noQA := strings.Replace(c.Body, "- `check-qa` (`rhiza guide check-qa`)", "- `check-xx` (`rhiza guide check-xx`)", 1)
	if got := checksWithoutRole(noQA, shipChange(t)); !slices.Equal(got, []string{"qa"}) {
		t.Errorf("checksWithoutRole without check-qa = %q, want [qa]", got)
	}
}
