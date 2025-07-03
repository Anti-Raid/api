package rpc

import (
	"context"
	"fmt"

	"github.com/Anti-Raid/api/state"
)

// Get Settings for a Guild User
func GetSettingsForGuildUser(
	ctx context.Context,
	guildID string,
	userID string,
) (res *any, err error) {
	return RpcQuery[any](
		ctx,
		state.IpcClient,
		"GET",
		fmt.Sprintf("%s/settings/%s/%s", CalcTWAddr(), guildID, userID),
		nil,
		true,
	)
}
