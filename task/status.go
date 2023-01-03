package task

type Status int

const (
	UncompletedStatus Status = iota
	InProgressStatus
	CompletedStatus
)

func (s Status) ToMarkdown() string {
	switch s {
	case UncompletedStatus:
		return "TODO"
	case InProgressStatus:
		return "IN-PROGRESS"
	case CompletedStatus:
		return "DONE"
	default:
		return "UNKNOWN"
	}
}
func (s Status) Next() Status {
	return Status(Clamp(0, int(s)+1, 2))
}

func (s Status) Prev() Status {
	return Status(Clamp(0, int(s)-1, 2))
}
