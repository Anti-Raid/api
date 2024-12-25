package rpc

import (
	"context"
	"fmt"

	"github.com/Anti-Raid/api/rpc_messages"
	"github.com/Anti-Raid/api/state"
)

func BotState(ctx context.Context) (*rpc_messages.BotState, error) {
	return RpcQuery[rpc_messages.BotState](
		ctx,
		state.IpcClient,
		"GET",
		fmt.Sprintf("%s/state", CalcBotAddr()),
		nil,
		true,
	)
}
