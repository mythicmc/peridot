package deploy

type PrepareUpdateError struct {
	Name string
	Type string
}

func (e PrepareUpdateError) Error() string {
	return "failed to prepare " + e.Type + " update for " + e.Name
}
