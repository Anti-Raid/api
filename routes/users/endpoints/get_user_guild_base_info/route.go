package get_user_guild_base_info

import (
	"net/http"
	"time"

	"github.com/Anti-Raid/api/rpc"
	"github.com/Anti-Raid/api/state"
	"github.com/Anti-Raid/api/types"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	docs "github.com/anti-raid/eureka/doclib"
	"github.com/anti-raid/eureka/ratelimit"
	"github.com/anti-raid/eureka/uapi"
)

func Docs() *docs.Doc {
	return &docs.Doc{
		Summary:     "Get User Guild Base Info",
		Description: "This endpoint will return basic user and guild information given their IDs",
		Resp:        types.UserGuildBaseData{},
		Params: []docs.Parameter{
			{
				Name:        "guild_id",
				Description: "Whether to refresh the user's guilds from discord",
				In:          "path",
				Required:    true,
				Schema:      docs.IdSchema,
			},
		},
	}
}

func Route(d uapi.RouteData, r *http.Request) uapi.HttpResponse {
	limit, err := ratelimit.Ratelimit{
		Expiry:      5 * time.Minute,
		MaxRequests: 5,
		Bucket:      "get_user_guild_base_info",
	}.Limit(d.Context, r)

	if err != nil {
		state.Logger.Error("Error while ratelimiting", zap.Error(err), zap.String("bucket", "get_user_guild_base_info"))
		return uapi.DefaultResponse(http.StatusInternalServerError)
	}

	if limit.Exceeded {
		return uapi.HttpResponse{
			Json: types.ApiError{
				Message: "You are being ratelimited. Please try again in " + limit.TimeToReset.String(),
			},
			Headers: limit.Headers(),
			Status:  http.StatusTooManyRequests,
		}
	}

	guildId := chi.URLParam(r, "guild_id")

	if guildId == "" {
		return uapi.DefaultResponse(http.StatusBadRequest)
	}

	bgui, err := rpc.BaseGuildUserInfo(d.Context, guildId, d.Auth.ID)

	if err != nil {
		return uapi.HttpResponse{
			Status: http.StatusInternalServerError,
			Json: types.ApiError{
				Message: "Error parsing bot response: " + err.Error(),
			},
		}
	}

	return uapi.HttpResponse{
		Json: types.UserGuildBaseData{
			OwnerID:   bgui.OwnerID,
			Name:      bgui.Name,
			Icon:      bgui.Icon,
			Roles:     bgui.Roles,
			UserRoles: bgui.UserRoles,
			BotRoles:  bgui.BotRoles,
			Channels:  bgui.Channels,
		},
	}
}
