package model

type Role int // enum in golang

const (
	MEMBER Role = iota
	ADMIN
)

func (r Role) String() string {
	return []string{"MEMBER", "ADMIN"}[r]
}
