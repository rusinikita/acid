package sequence

import "github.com/rusinikita/acid/call"

type Sequence struct {
	Calls         []call.Step
	Description   string
	LearningLinks []string
}
