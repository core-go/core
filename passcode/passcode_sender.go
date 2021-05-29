package passcode

import (
	"context"
	"time"
)

type Sender interface {
	Send(ctx context.Context, to string, passcode string, expireAt time.Time, params interface{}) error
}
