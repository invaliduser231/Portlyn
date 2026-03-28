package main

import (
	"encoding/json"
	"fmt"
	"os"

	"portlyn/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	issues := cfg.ValidationIssues()
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(map[string]any{
		"ok":     cfg.Validate() == nil,
		"issues": issues,
	})
	for _, issue := range issues {
		if issue.Level == "error" {
			os.Exit(1)
		}
	}
}
