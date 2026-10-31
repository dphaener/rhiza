package quiettest

import (
	"fmt"
	"slices"
	"testing"
	"time"
)

// TestTestProcesses reads ps output as Linux (procps) and darwin print it
// for the same go test run: this process is pid 300 under the go command,
// pid 200.
func TestTestProcesses(t *testing.T) {
	const ps = `    1     0 /sbin/init splash
  100     1 zsh
  200   100 go test ./...
  300   200 /tmp/go-build1/b001/api.test -test.paniconexit0 -test.timeout=10m0s
  301   200 /tmp/go-build1/b002/slice.test -test.paniconexit0 -test.timeout=10m0s
  302   200 /usr/local/go/pkg/tool/linux_amd64/compile -o /tmp/go-build1/b003/_pkg_.a -trimpath /tmp/go-build1/b003=> -p index
  303   200 /usr/local/go/pkg/tool/linux_amd64/vet -atomic -bool -buildtags
  304   200 /usr/local/go/pkg/tool/linux_amd64/link -o /tmp/go-build1/b004/gen.test
  305   301 /tmp/go-build1/b002/slice.test -test.run=^TestChild$
  306   100 tail -5
  307   200 git status
  308   200 /private/var/folders/x/T/go-build2/b005/changeset.test -test.paniconexit0
bad line
`
	got := testProcesses(ps, 300, 200)
	want := []string{"slice.test", "compile", "vet", "link", "changeset.test"}
	if !slices.Equal(got, want) {
		t.Errorf("testProcesses = %q, want %q", got, want)
	}
	if got := testProcesses(ps, 300, 999); got != nil {
		t.Errorf("testProcesses under another parent = %q, want none", got)
	}
	if got := testProcesses("", 300, 200); got != nil {
		t.Errorf("testProcesses of nothing = %q, want none", got)
	}
}

func TestQuietString(t *testing.T) {
	cases := []struct {
		q    Quiet
		want string
	}{
		{Quiet{waited: 1750 * time.Millisecond}, "alone after 1.8s"},
		{Quiet{waited: 41340 * time.Millisecond, last: []string{"slice.test"}}, "alone after 41.3s (last beside slice.test)"},
		{Quiet{waited: 3 * time.Minute, last: []string{"index.test", "slice.test"}, running: []string{"slice.test"}}, "beside slice.test after 3m0s"},
		{Quiet{err: fmt.Errorf("ps: exec: no ps")}, "other packages unknown (ps: exec: no ps)"},
	}
	for _, c := range cases {
		if got := c.q.String(); got != c.want {
			t.Errorf("%+v.String() = %q, want %q", c.q, got, c.want)
		}
	}
}
