package constants

const (
	GET = "GET"
	POST = "POST"
	PUT = "PUT"
	PATCH = "PATCH"
	DELETE = "DELETE"
	COPY = "COPY"
	HEAD = "HEAD"
	OPTIONS = "OPTIONS"
	LINK = "LINK"
	UNLINK = "UNLINK"
	PURGE = "PURGE"
	LOCK = "LOCK"
	UNLOCK = "UNLOCK"
	PROPFIND = "PROPFIND"
	VIEW = "VIEW"

	ActionNone    int32 = 0
	ActionRead    int32 = 1
	ActionWrite   int32 = 2
	ActionDelete  int32 = 4
	ActionApprove int32 = 8
	ActionAll     int32 = 2147483647

	ActionReadWrite              int32 = ActionRead | ActionWrite
	ActionReadWriteApprove       int32 = ActionRead | ActionWrite | ActionApprove
	ActionReadWriteApproveDelete int32 = ActionRead | ActionWrite | ActionApprove | ActionDelete
	ActionReadApprove            int32 = ActionRead | ActionApprove
	ActionReadApproveDelete      int32 = ActionRead | ActionApprove | ActionDelete
)
