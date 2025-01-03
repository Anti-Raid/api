package rpc

import (
	"context"
	"fmt"

	"github.com/Anti-Raid/api/rpc_messages"
	"github.com/Anti-Raid/api/state"
)

// Calls the CheckCommandPermission method to check whether or not a user has permission to run a command
func CheckCommandPermission(
	ctx context.Context,
	guildID string,
	userID string,
	command string,
) (res *rpc_messages.CheckCommandPermission, err error) {
	return RpcQuery[rpc_messages.CheckCommandPermission](
		ctx,
		state.IpcClient,
		"POST",
		fmt.Sprintf("%s/check-command-permission/%s/%s", CalcBotAddr(), guildID, userID),
		rpc_messages.CheckCommandPermissionRequest{
			Command: command,
		},
		true,
	)
}
