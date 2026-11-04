package api

import "testing"

func TestScanWorkflow(t *testing.T) {
	for _, tc := range []struct {
		name, text string
		want       workflowInfo
	}{
		{"block", "on:\n  pull_request:\n    branches: [ main ]\n  workflow_dispatch:\n    inputs:\n      sha:\n        required: true\nrun-name: CI ${{ inputs.sha }}\n",
			workflowInfo{pullRequest: true, dispatch: true, shaInput: true, runName: "CI ${{ inputs.sha }}"}},
		{"other input", "on:\n  pull_request:\n  workflow_dispatch:\n    inputs:\n      ref:\n        sha: x\n",
			workflowInfo{pullRequest: true, dispatch: true}},
		{"list", "on: [push, pull_request, workflow_dispatch]\n", workflowInfo{pullRequest: true, dispatch: true}},
		{"push only", "on:\n  push:\n", workflowInfo{}},
		{"flow map", "on: {pull_request: {}, workflow_dispatch: {inputs: {sha: {}}}}\n", workflowInfo{pullRequest: true, dispatch: true, shaInput: true}},
	} {
		if got := scanWorkflow(tc.text); got != tc.want {
			t.Errorf("%s: %+v, want %+v", tc.name, got, tc.want)
		}
	}
}
