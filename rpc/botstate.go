package rpc

import (
	"context"
	"fmt"

	"github.com/Anti-Raid/api/rpc_messages"
	"github.com/Anti-Raid/api/state"
)

func BotState(ctx context.Context) (*rpc_messages.TWState, error) {
	return RpcQuery[rpc_messages.TWState](
		ctx,
		state.IpcClient,
		"GET",
		fmt.Sprintf("%s/state", CalcTWAddr()),
		nil,
		true,
	)
}
