package main

type MsgPrefix string

const (
	MSGPREFIX_JOINED  MsgPrefix = "001"
	MSGPREFIX_REMOVED MsgPrefix = "002"
	MSGPREFIX_MESSAGE MsgPrefix = "003"
	MSGPREFIX_ERROR   MsgPrefix = "999"
)

type ErrorCode int

const (
	ERRORCODE_JOINED ErrorCode = 1
)

func getErrorCodeMessage(errorCode ErrorCode) string {
	message := "Error"
	switch errorCode {
	case ERRORCODE_JOINED:
		message = "接続できませんでした。"
	}

	return message
}
