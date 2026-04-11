package initcmd

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed templates
var templates embed.FS

// Run scaffolds a learning environment.
// Usage: acid init          → scaffold in current directory
//
//	acid init <dir>    → create <dir> and scaffold inside it
func Run(args []string) {
	targetDir := "."
	if len(args) > 0 {
		targetDir = args[0]
	}

	if targetDir != "." {
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "init: %v\n", err)
			os.Exit(1)
		}
	}

	var created []string

	flat := map[string]string{
		".env":               "templates/env.example",
		"agents.md":          "templates/agents.md",
		"learning_plan.md":   "templates/learning_plan.md",
		"docker-compose.yml": "templates/docker-compose.yml",
		"Makefile":           "templates/Makefile",
	}
	for dest, src := range flat {
		data, _ := templates.ReadFile(src)
		full := filepath.Join(targetDir, dest)
		if writeNew(full, data) {
			created = append(created, dest)
		}
	}

	_ = fs.WalkDir(templates, "templates/sequences", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, _ := filepath.Rel("templates/sequences", path)
		dest := filepath.Join(targetDir, "sequences", rel)
		_ = os.MkdirAll(filepath.Dir(dest), 0755)
		data, _ := templates.ReadFile(path)
		if writeNew(dest, data) {
			created = append(created, filepath.Join("sequences", rel))
		}
		return nil
	})

	printSummary(targetDir, created)
}

// writeNew writes content to path only if the file does not exist.
// Returns true if the file was written.
func writeNew(path string, content []byte) bool {
	if _, err := os.Stat(path); err == nil {
		fmt.Printf("  skip (exists)  %s\n", path)
		return false
	}
	if err := os.WriteFile(path, content, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "init: write %s: %v\n", path, err)
		os.Exit(1)
	}
	return true
}

func printSummary(dir string, created []string) {
	fmt.Printf("\nScaffolded learning environment in: %s\n\n", dir)
	fmt.Println("Created:")
	for _, f := range created {
		fmt.Printf("  %s\n", f)
	}
	fmt.Println()
	fmt.Println("Next steps:")
	step := 1
	if dir != "." {
		fmt.Printf("  %d. cd %s\n", step, dir)
		step++
	}
	fmt.Printf("  %d. Start a database:  make pg   (PostgreSQL) or  make mysql\n", step)
	step++
	fmt.Printf("  %d. Edit .env if you're using a custom database connection\n", step)
	step++
	fmt.Printf("  %d. Open two terminal panes:\n", step)
	fmt.Println("       FIRST    make serve")
	fmt.Println("       SECOND   claude --system-prompt agents.md")
	step++
	fmt.Printf("  %d. Say \"Let's start\" — the AI agent will guide you\n", step)
	step++
	fmt.Printf("  %d. Run scenarios manually at any time:\n", step)
	fmt.Println("       acid run -f sequences/lost_update.toml")
	step++
	fmt.Printf("  %d. Keep acid up to date:\n", step)
	fmt.Println("       make update")
	fmt.Println()
}
