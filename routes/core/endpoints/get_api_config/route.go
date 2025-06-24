package get_api_config

import (
	"net/http"

	"github.com/Anti-Raid/api/state"
	"github.com/Anti-Raid/api/types"

	docs "github.com/anti-raid/eureka/doclib"
	"github.com/anti-raid/eureka/uapi"
)

func Docs() *docs.Doc {
	return &docs.Doc{
		Summary:     "Get API Config",
		Description: "This endpoint returns the API configuration.",
		Resp:        []types.ApiConfig{},
		Params:      []docs.Parameter{},
	}
}

func Route(d uapi.RouteData, r *http.Request) uapi.HttpResponse {
	return uapi.HttpResponse{
		Json: types.ApiConfig{
			MainServer:          state.Config.Servers.Main,
			SupportServerInvite: state.Config.Meta.SupportServerInvite,
			ClientID:            state.Config.DiscordAuth.ClientID,
		},
	}
}
