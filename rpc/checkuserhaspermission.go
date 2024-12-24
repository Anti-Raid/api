package rpc

import (
	"context"
	"fmt"

	"github.com/Anti-Raid/api/rpc_messages"
	"github.com/Anti-Raid/api/state"
)

// check_user_has_permission
// /check-user-has-permission/:guild_id/:user_id
func CheckUserHasPermission(
	ctx context.Context,
	guildID string,
	userID string,
	perm string,
) (res *rpc_messages.CheckCommandPermission, err error) {
	return RpcQuery[rpc_messages.CheckCommandPermission](
		ctx,
		state.IpcClient,
		"POST",
		fmt.Sprintf("%s/check-user-has-permission/%s/%s", CalcBotAddr(), guildID, userID),
		rpc_messages.CheckUserHasKittycatPermissionsRequest{
			Perm: perm,
		},
		true,
	)
}
