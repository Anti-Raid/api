package rpc

import (
	"context"
	"fmt"

	"github.com/Anti-Raid/api/state"
	"github.com/Anti-Raid/corelib_go/silverpelt"
)

func Modules(ctx context.Context) (*[]silverpelt.CanonicalModule, error) {
	return RpcQuery[[]silverpelt.CanonicalModule](
		ctx,
		state.IpcClient,
		"GET",
		fmt.Sprintf("%s/modules", CalcBotAddr()),
		nil,
		true,
	)
}
