package permission

const (
	ActionNone    int32 = 0
	ActionRead    int32 = 1
	ActionWrite   int32 = 2
	ActionDelete  int32 = 4
	ActionApprove int32 = 8
	ActionAll     int32 = 2147483647
)
