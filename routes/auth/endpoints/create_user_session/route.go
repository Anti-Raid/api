package create_user_session

import (
	"net/http"
	"time"

	"github.com/Anti-Raid/api/state"
	"github.com/Anti-Raid/api/types"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"github.com/anti-raid/eureka/crypto"
	docs "github.com/anti-raid/eureka/doclib"
	"github.com/anti-raid/eureka/uapi"
)

var (
	compiledMessages = uapi.CompileValidationErrors(types.CreateUserSession{})
)

func Docs() *docs.Doc {
	return &docs.Doc{
		Summary:     "Create User Session",
		Description: "Creates a user session returning the session token. The session token cannot be read after creation.",
		Req:         types.CreateUserSession{},
		Resp:        types.CreateUserSessionResponse{},
		Params:      []docs.Parameter{},
	}
}

func Route(d uapi.RouteData, r *http.Request) uapi.HttpResponse {
	var createData types.CreateUserSession

	hresp, ok := uapi.MarshalReq(r, &createData)

	if !ok {
		return hresp
	}

	err := state.Validator.Struct(createData)

	if err != nil {
		return uapi.ValidatorErrorResponse(compiledMessages, err.(validator.ValidationErrors))
	}

	if createData.Name == "" {
		return uapi.HttpResponse{
			Status: http.StatusBadRequest,
			Json:   types.ApiError{Message: "Name is required"},
		}
	}

	if createData.Type == "" {
		return uapi.HttpResponse{
			Status: http.StatusBadRequest,
			Json:   types.ApiError{Message: "Type is required"},
		}
	}

	if createData.Expiry <= 0 {
		return uapi.HttpResponse{
			Status: http.StatusBadRequest,
			Json:   types.ApiError{Message: "Expiry must be greater than 0"},
		}
	}

	// Create session
	sessionToken := crypto.RandString(128)
	var sessionId string

	expiry := time.Now().Add(time.Duration(createData.Expiry) * time.Second)

	err = state.Pool.QueryRow(d.Context, "INSERT INTO web_api_tokens (token, user_id, name, type, expiry) VALUES ($1, $2, $3, $4, $5) RETURNING id", sessionToken, d.Auth.ID, createData.Name, createData.Type, expiry).Scan(&sessionId)

	if err != nil {
		state.Logger.Error("Error while creating user session", zap.Error(err))
		return uapi.DefaultResponse(http.StatusInternalServerError)
	}

	return uapi.HttpResponse{
		Status: http.StatusCreated,
		Json: types.CreateUserSessionResponse{
			Token:     sessionToken,
			UserID:    d.Auth.ID,
			SessionID: sessionId,
			Expiry:    expiry,
		},
	}

}
