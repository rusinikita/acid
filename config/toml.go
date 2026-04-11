package config

import (
	"fmt"
	"strings"

	"github.com/BurntSushi/toml"

	"github.com/rusinikita/acid/call"
	"github.com/rusinikita/acid/sequence"
)

type fileStep struct {
	SQL   string `toml:"sql"`
	Cmd   string `toml:"cmd"`
	Trx   string `toml:"trx"`
	Setup bool   `toml:"setup"`
}

type file struct {
	Name          string     `toml:"name"`
	Description   string     `toml:"description"`
	LearningLinks []string   `toml:"learning_links"`
	DropTables    []string   `toml:"drop_tables"`
	Steps         []fileStep `toml:"steps"`
}

func Load(path string) (sequence.Sequence, error) {
	var f file
	if _, err := toml.DecodeFile(path, &f); err != nil {
		return sequence.Sequence{}, fmt.Errorf("decode %s: %w", path, err)
	}

	steps := make([]call.Step, 0, len(f.Steps))
	for i, s := range f.Steps {
		step, err := convertStep(s)
		if err != nil {
			return sequence.Sequence{}, fmt.Errorf("step %d: %w", i+1, err)
		}
		steps = append(steps, step)
	}

	return sequence.Sequence{
		Name:          f.Name,
		Description:   f.Description,
		LearningLinks: f.LearningLinks,
		DropTables:    f.DropTables,
		Calls:         steps,
	}, nil
}

func convertStep(s fileStep) (call.Step, error) {
	trx := call.TrxID(s.Trx)

	switch strings.ToLower(s.Cmd) {
	case "begin":
		return call.Begin(trx), nil
	case "commit":
		return call.Commit(trx), nil
	case "rollback":
		return call.Rollback(trx), nil
	case "":
		if s.Setup {
			return call.Setup(s.SQL), nil
		}
		if s.Trx != "" {
			return call.Call(s.SQL, trx), nil
		}
		return call.Call(s.SQL), nil
	default:
		return call.Step{}, fmt.Errorf("unknown cmd %q", s.Cmd)
	}
}
