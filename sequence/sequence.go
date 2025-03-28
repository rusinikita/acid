package sequence

import "github.com/rusinikita/acid/call"

type Sequence struct {
	Name          string
	Description   string
	Calls         []call.Step
	LearningLinks []string
}
