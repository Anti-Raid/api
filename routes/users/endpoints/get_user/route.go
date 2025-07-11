package get_user

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Anti-Raid/api/state"
	"github.com/Anti-Raid/api/types"
	"github.com/Anti-Raid/api/utils"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"

	docs "github.com/anti-raid/eureka/doclib"
	"github.com/anti-raid/eureka/dovewing"
	"github.com/anti-raid/eureka/uapi"
)

var (
	userRows    = utils.GetCols(types.User{})
	userRowsStr = strings.Join(userRows, ", ")
)

func Docs() *docs.Doc {
	return &docs.Doc{
		Summary:     "Get User",
		Description: "This endpoint will return user information given their ID",
		Resp:        types.User{},
		Params: []docs.Parameter{
			{
				Name:        "id",
				Description: "The ID of the user to get information about",
				In:          "path",
				Required:    true,
				Schema:      docs.IdSchema,
			},
		},
	}
}

func Route(d uapi.RouteData, r *http.Request) uapi.HttpResponse {
	id := chi.URLParam(r, "id")

	row, err := state.Pool.Query(d.Context, "SELECT "+userRowsStr+" FROM users WHERE user_id = $1", id)

	if err != nil {
		state.Logger.Error("Error querying database", zap.Error(err))
		return uapi.DefaultResponse(http.StatusInternalServerError)
	}

	user, err := pgx.CollectOneRow[types.User](row, pgx.RowToStructByName)

	if errors.Is(err, pgx.ErrNoRows) {
		return uapi.DefaultResponse(http.StatusNotFound)
	}

	user.User, err = dovewing.GetUser(d.Context, id, state.DovewingPlatformDiscord)

	if err != nil {
		state.Logger.Error("Error getting user from dovewing", zap.Error(err))
		return uapi.DefaultResponse(http.StatusInternalServerError)
	}

	return uapi.HttpResponse{
		Json: user,
	}
}
