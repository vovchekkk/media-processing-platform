package domain

type TaskStatus string

const (
	StatusInProgress TaskStatus = "in_progress"
	StatusReady      TaskStatus = "ready"
	StatusFailed     TaskStatus = "failed"
)
