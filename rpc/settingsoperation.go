package rpc

import (
	"context"
	"fmt"

	"github.com/Anti-Raid/api/rpc_messages"
	"github.com/Anti-Raid/api/state"
)

// Calls the SettingsOperation method to execute a settings operation (settings-operation/:guild_id/:user_id)
func SettingsOperation(
	ctx context.Context,
	guildID string,
	userID string,
	settingsOpReq *rpc_messages.SettingsOperationRequest,
) (res *rpc_messages.CanonicalSettingsResult, err error) {
	return RpcQuery[rpc_messages.CanonicalSettingsResult](
		ctx,
		state.IpcClient,
		"POST",
		fmt.Sprintf("%s/settings-operation/%s/%s", CalcTWAddr(), guildID, userID),
		settingsOpReq,
		true,
	)
}
