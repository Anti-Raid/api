package get_user_guilds

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/Anti-Raid/api/rpc"
	"github.com/Anti-Raid/api/state"
	"github.com/Anti-Raid/api/types"
	"github.com/Anti-Raid/api/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"

	docs "github.com/anti-raid/eureka/doclib"
	"github.com/anti-raid/eureka/ratelimit"
	"github.com/anti-raid/eureka/uapi"
)

func Docs() *docs.Doc {
	return &docs.Doc{
		Summary:     "Get User Guilds",
		Description: "This endpoint will return all guilds the user is in along with information regarding whether the bot is in said guild",
		Resp:        types.DashboardGuildData{},
		Params: []docs.Parameter{
			{
				Name:        "refresh",
				Description: "Whether to refresh the user's guilds from discord",
				In:          "query",
				Required:    false,
				Schema:      docs.BoolSchema,
			},
		},
	}
}

func Route(d uapi.RouteData, r *http.Request) uapi.HttpResponse {
	var limit ratelimit.Limit
	var err error
	if r.URL.Query().Get("refresh") == "true" {
		limit, err = ratelimit.Ratelimit{
			Expiry:      5 * time.Minute,
			MaxRequests: 3,
			Bucket:      "get_user_guilds_refresh",
		}.Limit(d.Context, r)

		if err != nil {
			state.Logger.Error("Error while ratelimiting", zap.Error(err), zap.String("bucket", "get_user_guilds_refresh"))
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
	} else {
		limit, err = ratelimit.Ratelimit{
			Expiry:      5 * time.Minute,
			MaxRequests: 5,
			Bucket:      "get_user_guilds_norefresh",
		}.Limit(d.Context, r)

		if err != nil {
			state.Logger.Error("Error while ratelimiting", zap.Error(err), zap.String("bucket", "login"))
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
	}

	refresh := r.URL.Query().Get("refresh") == "true"

	if !refresh {
		// Once case where we must refresh is if guilds_cache is NULL
		var count int64

		err := state.Pool.QueryRow(d.Context, "SELECT COUNT(*) FROM users WHERE user_id = $1 AND guilds_cache IS NULL", d.Auth.ID).Scan(&count)

		if err != nil {
			state.Logger.Error("Failed to query database", zap.Error(err))
			return uapi.DefaultResponse(http.StatusInternalServerError)
		}

		if count != 0 {
			refresh = true
		}
	}

	// Fetch guild list of user
	var dashguilds []*types.DashboardGuild
	if refresh {
		var accesstoken pgtype.Text

		err = state.Pool.QueryRow(d.Context, "SELECT access_token FROM users WHERE user_id = $1", d.Auth.ID).Scan(&accesstoken)

		if err != nil {
			state.Logger.Error("Failed to query database", zap.Error(err))
			return uapi.DefaultResponse(http.StatusInternalServerError)
		}

		if !accesstoken.Valid {
			return uapi.HttpResponse{
				Status: http.StatusBadRequest,
				Json: types.ApiError{
					Message: "User has not authorized via oauth2 yet!",
				},
			}
		}

		// Refresh guilds
		httpReq, err := http.NewRequestWithContext(d.Context, "GET", state.Config.Meta.Proxy+"/api/v10/users/@me/guilds", nil)

		if err != nil {
			state.Logger.Error("Failed to create oauth2 request to discord", zap.Error(err))
			return uapi.HttpResponse{
				Json: types.ApiError{
					Message: "Failed to create request to Discord to fetch guilds",
				},
				Status:  http.StatusInternalServerError,
				Headers: limit.Headers(),
			}
		}

		httpReq.Header.Set("Authorization", "Bearer "+accesstoken.String)

		cli := &http.Client{}

		httpResp, err := cli.Do(httpReq)

		if err != nil {
			state.Logger.Error("Failed to send oauth2 request to discord", zap.Error(err))
			return uapi.HttpResponse{
				Json: types.ApiError{
					Message: "Failed to send oauth2 request to Discord [user guilds]",
				},
				Status:  http.StatusInternalServerError,
				Headers: limit.Headers(),
			}
		}

		defer httpResp.Body.Close()

		body, err := io.ReadAll(httpResp.Body)

		if err != nil {
			state.Logger.Error("Failed to read oauth2 response from discord", zap.Error(err))
			return uapi.HttpResponse{
				Json: types.ApiError{
					Message: "Failed to read oauth2 response from Discord [user guilds]",
				},
				Status:  http.StatusInternalServerError,
				Headers: limit.Headers(),
			}
		}

		var guilds []*discordgo.UserGuild

		err = json.Unmarshal(body, &guilds)

		if err != nil {
			state.Logger.Error("Failed to parse oauth2 response from discord", zap.Error(err))
			return uapi.HttpResponse{
				Json: types.ApiError{
					Message: "Failed to parse oauth2 response from Discord [user guilds]",
				},
				Status:  http.StatusInternalServerError,
				Headers: limit.Headers(),
			}
		}

		for _, guild := range guilds {
			dashguilds = append(dashguilds, &types.DashboardGuild{
				ID:          guild.ID,
				Name:        guild.Name,
				Permissions: guild.Permissions,
				Avatar: func() string {
					return utils.IconURL(guild.Icon, discordgo.EndpointGuildIcon(guild.ID, guild.Icon), discordgo.EndpointGuildIconAnimated(guild.ID, guild.Icon), "64")
				}(),
			})
		}

		// Now update the database
		_, err = state.Pool.Exec(d.Context, "UPDATE users SET guilds_cache = $1 WHERE user_id = $2", dashguilds, d.Auth.ID)

		if err != nil {
			state.Logger.Error("Failed to update database", zap.Error(err))
			return uapi.DefaultResponse(http.StatusInternalServerError)
		}
	} else {
		err := state.Pool.QueryRow(d.Context, "SELECT guilds_cache FROM users WHERE user_id = $1", d.Auth.ID).Scan(&dashguilds)

		if err != nil {
			state.Logger.Error("Failed to query database", zap.Error(err))
			return uapi.DefaultResponse(http.StatusInternalServerError)
		}
	}

	// Get list of guild ids
	var guilds = []string{}

	for _, guild := range dashguilds {
		guilds = append(guilds, guild.ID)
	}

	// Now send the requests
	var botInGuild []string
	var unknownGuilds []string

	guildsExistResp, err := rpc.GuildsExist(d.Context, guilds)

	if err != nil {
		state.Logger.Error("Failed to check if bot is in guilds", zap.Error(err))
		return uapi.HttpResponse{
			Status: http.StatusInternalServerError,
			Json: types.ApiError{
				Message: "Failed to check if bot is in guilds: " + err.Error(),
			},
		}
	}

	guildsExist := *guildsExistResp

	if len(guildsExist) != len(guilds) {
		state.Logger.Error("Mismatch in guildsExist response", zap.Any("guildsExist", guildsExist), zap.Any("guilds", guilds))
		return uapi.HttpResponse{
			Status: http.StatusInternalServerError,
			Json: types.ApiError{
				Message: "Mismatch in guildsExist response",
			},
		}
	}

	for i, v := range guildsExist {
		if v == 1 {
			botInGuild = append(botInGuild, guilds[i])
		}
	}

	return uapi.HttpResponse{
		Json: &types.DashboardGuildData{
			Guilds:        dashguilds,
			BotInGuilds:   botInGuild,
			UnknownGuilds: unknownGuilds,
		},
		Headers: limit.Headers(),
	}
}
