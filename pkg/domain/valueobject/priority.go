package valueobject

type Priority int

const (
	Low Priority = iota + 1
	Medium
	High
)
