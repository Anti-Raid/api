package revoke_user_session

import (
	"errors"
	"net/http"

	"github.com/Anti-Raid/api/state"
	"github.com/Anti-Raid/api/types"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"

	docs "github.com/anti-raid/eureka/doclib"
	"github.com/anti-raid/eureka/uapi"
)

func Docs() *docs.Doc {
	return &docs.Doc{
		Summary:     "Revoke User Session",
		Description: "Revokes a session of a user based on session ID",
		Resp:        types.ApiError{},
		Params: []docs.Parameter{
			{
				Name:        "session_id",
				Description: "The ID of the session to revoke",
				In:          "path",
				Required:    true,
				Schema:      docs.IdSchema,
			},
		},
	}
}

func Route(d uapi.RouteData, r *http.Request) uapi.HttpResponse {
	sessionId := chi.URLParam(r, "session_id")

	if sessionId == "" {
		return uapi.DefaultResponse(http.StatusBadRequest)
	}

	var count int64

	err := state.Pool.QueryRow(d.Context, "SELECT COUNT(*) FROM web_api_tokens WHERE user_id = $1 AND id = $2", d.Auth.ID, sessionId).Scan(&count)

	if errors.Is(err, pgx.ErrNoRows) {
		return uapi.HttpResponse{
			Status: http.StatusNotFound,
			Json:   types.ApiError{Message: "No sessions found"},
		}
	}

	if err != nil {
		state.Logger.Error("Error while getting user session", zap.Error(err))
		return uapi.DefaultResponse(http.StatusInternalServerError)
	}

	if count == 0 {
		return uapi.HttpResponse{
			Status: http.StatusNotFound,
			Json:   types.ApiError{Message: "No sessions found"},
		}
	}

	_, err = state.Pool.Exec(d.Context, "DELETE FROM web_api_tokens WHERE id = $1", sessionId)

	if err != nil {
		state.Logger.Error("Error while revoking user session", zap.Error(err))
		return uapi.DefaultResponse(http.StatusInternalServerError)
	}

	return uapi.DefaultResponse(http.StatusNoContent)
}
