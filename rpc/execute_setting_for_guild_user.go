package rpc

import (
	"context"
	"fmt"

	"github.com/Anti-Raid/api/rpc_messages"
	"github.com/Anti-Raid/api/state"
)

// Calls the SettingsOperation method to execute a settings operation
func ExecuteSettingForGuildUser(
	ctx context.Context,
	guildID string,
	userID string,
	settingsOpReq *rpc_messages.SettingsOperationRequest,
) (res *map[string]rpc_messages.DispatchResult, err error) {
	return RpcQuery[map[string]rpc_messages.DispatchResult](
		ctx,
		state.IpcClient,
		"POST",
		fmt.Sprintf("%s/settings/%s/%s", CalcTWAddr(), guildID, userID),
		settingsOpReq,
		true,
	)
}
