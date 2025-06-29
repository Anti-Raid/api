package test_auth

import (
	"context"
	"net/http"

	api "github.com/Anti-Raid/api/auth"
	"github.com/Anti-Raid/api/types"

	docs "github.com/anti-raid/eureka/doclib"
	"github.com/anti-raid/eureka/uapi"

	"github.com/go-chi/chi/v5"
)

func Docs() *docs.Doc {
	return &docs.Doc{
		Summary:     "Test Auth",
		Description: "Test your authentication",
		Req:         types.TestAuth{},
		Resp:        uapi.AuthData{},
	}
}

func Route(d uapi.RouteData, r *http.Request) uapi.HttpResponse {
	var payload types.TestAuth

	hresp, ok := uapi.MarshalReq(r, &payload)

	if !ok {
		return hresp
	}

	if payload.TargetID == "" {
		return uapi.HttpResponse{
			Status: http.StatusBadRequest,
			Json:   types.ApiError{Message: "Target ID is required"},
		}
	}

	// Create []AuthType
	rctx := context.Background()
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("test", payload.TargetID)
	rctx = context.WithValue(rctx, chi.RouteCtxKey, ctx)
	authType := []uapi.AuthType{
		{
			URLVar:       "test",
			Type:         payload.AuthType,
			AllowedScope: "ban_exempt",
		},
	}

	reqCtxd := r.WithContext(rctx)

	r.Header.Set("Authorization", payload.Token)

	// Check auth
	authData, hr, ok := api.Authorize(uapi.Route{
		Auth: authType,
	}, reqCtxd)

	if !ok {
		return hr
	}

	return uapi.HttpResponse{
		Json: authData,
	}
}
