package models

type Role string

const (
	SuperAdmin Role = "super_admin"
	GroupAdmin Role = "group_admin"
	Member     Role = "member"
)
