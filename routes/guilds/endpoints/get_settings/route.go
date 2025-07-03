package get_settings

import (
	"net/http"
	"time"

	"github.com/Anti-Raid/api/rpc"
	"github.com/Anti-Raid/api/state"
	"github.com/Anti-Raid/api/types"
	docs "github.com/anti-raid/eureka/doclib"
	"github.com/anti-raid/eureka/ratelimit"
	"github.com/anti-raid/eureka/uapi"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func Docs() *docs.Doc {
	return &docs.Doc{
		Summary:     "Get Settings",
		Description: "Get Settings for a Guild User",
		Req:         types.SettingsExecute{},
		Resp:        map[string]any{},
		Params: []docs.Parameter{
			{
				Name:        "guild_id",
				Description: "The guild ID to execute the operation in",
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
		MaxRequests: 10,
		Bucket:      "settings_execute",
	}.Limit(d.Context, r)

	if err != nil {
		state.Logger.Error("Error while ratelimiting", zap.Error(err), zap.String("bucket", "settings_execute"))
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

	var body types.SettingsExecute

	hresp, ok := uapi.MarshalReqWithHeaders(r, &body, limit.Headers())

	if !ok {
		return hresp
	}

	if body.Setting == "" {
		return uapi.HttpResponse{
			Status: http.StatusBadRequest,
			Json: types.ApiError{
				Message: "`setting` must be provided",
			},
		}
	}

	resp, err := rpc.GetSettingsForGuildUser(
		d.Context,
		guildId,
		d.Auth.ID,
	)

	if err != nil {
		state.Logger.Error("Error executing settings operation", zap.Error(err))
		return uapi.HttpResponse{
			Status: http.StatusInternalServerError,
			Json: types.ApiError{
				Message: "Error executing settings operation: " + err.Error(),
			},
		}
	}

	return uapi.HttpResponse{
		Status: http.StatusOK,
		Json:   resp,
	}
}
