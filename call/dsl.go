package call

func Setup(code string) Step {
	return Step{
		Code:      code,
		TestSetup: true,
	}
}

func Call(code string, id ...TrxID) Step {
	call := Step{
		Code: code,
	}

	if len(id) > 0 {
		call.Trx = id[0]
	}

	if len(id) > 1 {
		panic("more than 1 trx id not supported: " + call.Code)
	}

	return call
}

func Begin(id TrxID) Step {
	return Step{
		Trx:        id,
		TrxCommand: TrxBegin,
	}
}

func Commit(id TrxID) Step {
	return Step{
		Trx:        id,
		TrxCommand: TrxCommit,
	}
}

func Rollback(id TrxID) Step {
	return Step{
		Trx:        id,
		TrxCommand: TrxRollback,
	}
}
