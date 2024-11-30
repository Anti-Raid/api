package rpc

import (
	"context"
	"fmt"

	"github.com/Anti-Raid/api/rpc_messages"
	"github.com/Anti-Raid/api/state"
)

// Calls the CheckCommandPermission method to check whether or not a command is runnable
func JobserverSpawnTask(
	ctx context.Context,
	spawnTask *rpc_messages.JobserverSpawn,
) (res *rpc_messages.JobserverSpawnResponse, err error) {
	return RpcQuery[rpc_messages.JobserverSpawnResponse](
		ctx,
		state.IpcClient,
		"POST",
		fmt.Sprintf("%s/spawn", CalcJobserverAddr()),
		spawnTask,
		true,
	)
}
