// Package projectctx is the project a Rhiza TUI app is in: the registered
// project whose folder, or <project>.<slug> sibling worktree, holds the
// directory the app was opened in. At the root, or outside any project,
// there is none, and the header says so.
package projectctx

import "github.com/dphaener/rhiza/internal/rootcfg"

// NoProject is the header's text when the app is in no project.
const NoProject = "No project"

// Current is the name of the project that dir belongs to. It reports false
// at the root itself, outside any root and when the root cannot be read.
func Current(dir string) (string, bool) {
	root, ok := rootcfg.Find(dir)
	if !ok {
		return "", false
	}
	cfg, err := rootcfg.Load(root)
	if err != nil {
		return "", false
	}
	return cfg.ProjectOf(root, dir)
}

// Header is the one line an app shows above its content: the project's
// name, or NoProject.
func Header(dir string) string {
	if name, ok := Current(dir); ok {
		return "Project: " + name
	}
	return NoProject
}
