package api

// This file holds guide and guide ls: the agent guide that ships in the
// binary, read from internal/guide. Neither needs a repository.

import (
	"fmt"
	"io"
	"reflect"

	"github.com/dphaener/genesis/internal/guide"
)

func init() {
	register(Operation{
		Name:     "guide",
		Summary:  "Print the agent guide: the genesis skill, the scrub skill or a topic",
		Request:  reflect.TypeFor[GuideRequest](),
		Response: reflect.TypeFor[GuideResponse](),
		Run:      runGuide,
	})
	register(Operation{
		Name:     "guide.ls",
		Summary:  "List the guide's skills and topics with their sizes",
		Request:  reflect.TypeFor[GuideLsRequest](),
		Response: reflect.TypeFor[GuideLsResponse](),
		Run:      runGuideLs,
	})
}

// GuideRequest asks for one text of the guide.
type GuideRequest struct {
	Name string `json:"name,omitempty" arg:"0" desc:"The text to print: genesis (the default), scrub or a topic name; genesis guide ls lists them"`
}

// GuideResponse is one text of the guide.
type GuideResponse struct {
	Name    string `json:"name" desc:"The text's name"`
	Version string `json:"version" desc:"The genesis version the text ships in"`
	Text    string `json:"text" desc:"The text, as markdown"`
}

// Human prints the text exactly.
func (r GuideResponse) Human(w io.Writer) {
	io.WriteString(w, r.Text)
}

// GuideLsRequest lists the guide's texts; it takes nothing.
type GuideLsRequest struct{}

// GuideLsResponse lists the guide's texts.
type GuideLsResponse struct {
	Texts []GuideTextInfo `json:"texts" desc:"The texts: the genesis skill, the scrub skill, then the topics by name"`
}

// GuideTextInfo describes one text of the guide.
type GuideTextInfo struct {
	Name        string `json:"name" desc:"What genesis guide takes"`
	Kind        string `json:"kind" enum:"skill,topic" desc:"skill or topic"`
	Tokens      int    `json:"tokens" desc:"The text's estimated size in tokens"`
	Description string `json:"description" desc:"When to load it"`
}

// Human prints one line per text: name, kind, tokens and description.
func (r GuideLsResponse) Human(w io.Writer) {
	width := 0
	for _, t := range r.Texts {
		width = max(width, len(t.Name))
	}
	for _, t := range r.Texts {
		fmt.Fprintf(w, "- %-*s  %-5s  %5d tokens  %s\n", width, t.Name, t.Kind, t.Tokens, t.Description)
	}
}

func runGuide(ctx *Ctx, req any) (any, error) {
	r := req.(GuideRequest)
	t, ok := guide.Lookup(r.Name)
	if !ok {
		return nil, NotFoundf("no guide text %q", r.Name).WithNext("genesis guide ls")
	}
	return GuideResponse{Name: t.Name, Version: ctx.Version, Text: t.Body}, nil
}

func runGuideLs(_ *Ctx, _ any) (any, error) {
	var resp GuideLsResponse
	for _, t := range guide.List() {
		kind := "topic"
		if t.Skill {
			kind = "skill"
		}
		resp.Texts = append(resp.Texts, GuideTextInfo{
			Name: t.Name, Kind: kind, Tokens: guide.Tokens(t.Body), Description: t.Description,
		})
	}
	return resp, nil
}
