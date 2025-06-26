package rpc

import (
	"context"
	"fmt"

	"github.com/Anti-Raid/api/rpc_messages"
	"github.com/Anti-Raid/api/state"
)

// Calls the BaseGuildUserInfo method to get basic user + guild info (base-guild-user-info/:guild_id/:user_id)
func BaseGuildUserInfo(
	ctx context.Context,
	guildID string,
	userID string,
) (res *rpc_messages.BaseGuildUserInfo, err error) {
	return RpcQuery[rpc_messages.BaseGuildUserInfo](
		ctx,
		state.IpcClient,
		"GET",
		fmt.Sprintf("%s/base-guild-user-info/%s/%s", CalcTWAddr(), guildID, userID),
		nil,
		true,
	)
}
