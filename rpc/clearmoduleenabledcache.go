package rpc

import (
	"context"
	"fmt"

	"github.com/Anti-Raid/api/rpc_messages"
	"github.com/Anti-Raid/api/state"
)

// Clears the bots internal modules enabled cache
func ClearModulesEnabledCache(ctx context.Context, data *rpc_messages.ClearModulesEnabledCacheRequest) (*rpc_messages.ClearModulesEnabledCacheResponse, error) {
	return RpcQuery[rpc_messages.ClearModulesEnabledCacheResponse](
		ctx,
		state.IpcClient,
		"POST",
		fmt.Sprintf("%s/clear-modules-enabled-cache", CalcBotAddr()),
		data,
		true,
	)
}
