package api

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/canon"
)

func callDescribe(t *testing.T, req string) (DescribeResponse, error) {
	t.Helper()
	o, ok := Lookup("describe")
	if !ok {
		t.Fatal("describe is not registered")
	}
	r, err := DecodeRequest(o, []byte(req))
	if err != nil {
		t.Fatalf("DecodeRequest(%s): %v", req, err)
	}
	resp, err := o.Run(&Ctx{}, r)
	if err != nil {
		return DescribeResponse{}, err
	}
	return resp.(DescribeResponse), nil
}

func TestDescribeAll(t *testing.T) {
	resp, err := callDescribe(t, `{}`)
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Operations) != len(Operations()) {
		t.Fatalf("got %d operations, want %d", len(resp.Operations), len(Operations()))
	}
	var d *OperationInfo
	for i, o := range resp.Operations {
		if o.Name == "describe" {
			d = &resp.Operations[i]
		}
	}
	if d == nil {
		t.Fatal("describe does not describe itself")
	}
	if d.Write || d.Summary == "" {
		t.Errorf("describe info = %+v", d)
	}
	want, err := canon.Marshal(Schema(reflect.TypeFor[DescribeRequest]()))
	if err != nil {
		t.Fatal(err)
	}
	got, err := canon.Marshal(d.Request)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Errorf("request schema =\n%s\nwant\n%s", got, want)
	}
	out, err := canon.Marshal(resp)
	if err != nil {
		t.Fatal(err)
	}
	var shape struct {
		Operations []map[string]json.RawMessage `json:"operations"`
	}
	if err := json.Unmarshal(out, &shape); err != nil {
		t.Fatal(err)
	}
	for _, k := range []string{"name", "summary", "write", "request", "response"} {
		if _, ok := shape.Operations[0][k]; !ok {
			t.Errorf("operation JSON lacks %q: %s", k, out)
		}
	}
}

func TestDescribeOne(t *testing.T) {
	resp, err := callDescribe(t, `{"op":"describe"}`)
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Operations) != 1 || resp.Operations[0].Name != "describe" {
		t.Errorf("got %+v", resp.Operations)
	}
}

func TestDescribeUnknown(t *testing.T) {
	_, err := callDescribe(t, `{"op":"nope"}`)
	if ExitCode(err) != 4 {
		t.Fatalf("err = %v, exit %d, want 4", err, ExitCode(err))
	}
	if AsError(err).Next != "rhiza describe" {
		t.Errorf("next = %q", AsError(err).Next)
	}
}

func TestDescribeHuman(t *testing.T) {
	resp, err := callDescribe(t, `{}`)
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	resp.Human(&buf)
	out := buf.String()
	for _, s := range []string{"## describe", "rhiza describe [<op>]", "`<op>`"} {
		if !strings.Contains(out, s) {
			t.Errorf("human output lacks %q:\n%s", s, out)
		}
	}
}

// Every response renders for people; --json alone is not enough.
func TestEveryResponseIsHumaner(t *testing.T) {
	h := reflect.TypeFor[Humaner]()
	for _, o := range Operations() {
		if !o.Response.Implements(h) {
			t.Errorf("%s: response %s has no Human(io.Writer) method", o.Name, o.Response)
		}
		Params(o.Request) // panics on a request breaking the tag contract
		Schema(o.Request)
		Schema(o.Response)
	}
}
