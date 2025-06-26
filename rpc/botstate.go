package rpc

import (
	"context"
	"fmt"

	"github.com/Anti-Raid/api/rpc_messages"
	"github.com/Anti-Raid/api/state"
)

// Temporarily does 2 RPC calls until bot is removed in favor of template worker
func BotState(ctx context.Context) (*rpc_messages.BotState, error) {
	bs, err := RpcQuery[rpc_messages.BSI](
		ctx,
		state.IpcClient,
		"GET",
		fmt.Sprintf("%s/state", CalcBotAddr()),
		nil,
		true,
	)

	if err != nil {
		return nil, err
	}

	ts, err := RpcQuery[rpc_messages.TSI](
		ctx,
		state.IpcClient,
		"GET",
		fmt.Sprintf("%s/state", CalcTWAddr()),
		nil,
		true,
	)

	if err != nil {
		return nil, err
	}

	return &rpc_messages.BotState{
		Commands:           bs.Commands,
		CommandPermissions: bs.CommandPermissions,
		Settings:           ts.Settings,
	}, nil
}
