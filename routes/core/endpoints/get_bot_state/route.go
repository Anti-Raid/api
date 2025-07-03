package get_bot_state

import (
	"net/http"

	"github.com/Anti-Raid/api/rpc"
	"github.com/Anti-Raid/api/state"
	"github.com/Anti-Raid/api/types"

	docs "github.com/anti-raid/eureka/doclib"
	"github.com/anti-raid/eureka/uapi"
)

func Docs() *docs.Doc {
	return &docs.Doc{
		Summary:     "Get Bot State",
		Description: "This endpoint returns all commands on AntiRaid.",
		Resp:        types.BotState{},
		Params:      []docs.Parameter{},
	}
}

var BotStateCache *types.BotState

func Route(d uapi.RouteData, r *http.Request) uapi.HttpResponse {
	if BotStateCache != nil {
		return uapi.HttpResponse{
			Json: BotStateCache,
		}
	}

	bs, err := rpc.BotState(state.Context)

	if err != nil {
		return uapi.HttpResponse{
			Status: http.StatusInternalServerError,
			Json:   types.ApiError{Message: "Error getting modules:" + err.Error()},
		}
	}

	BotStateCache = &types.BotState{
		Commands: bs.Commands,
	}

	return uapi.HttpResponse{
		Json: BotStateCache,
	}
}
