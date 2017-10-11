package main

type msgPrefix string

const (
	CONSTANTS_MSGPREFIX_JOINED  msgPrefix = "001"
	CONSTANTS_MSGPREFIX_REMOVED msgPrefix = "002"
	CONSTANTS_MSGPREFIX_MESSAGE msgPrefix = "003"
)
