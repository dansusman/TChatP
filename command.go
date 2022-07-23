package main

type commandUUID int

const (
	NAME commandUUID = iota
	JOIN
	GROUPS
	MSG
	ABORT
	HELP
	SWITCH
	LEAVE
)

type command struct {
	uuid   commandUUID
	args   []string
	client *client
}
