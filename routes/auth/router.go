package auth

import (
	"github.com/Anti-Raid/api/routes/auth/endpoints/create_oauth2_login"
	"github.com/Anti-Raid/api/routes/auth/endpoints/create_user_session"
	"github.com/Anti-Raid/api/routes/auth/endpoints/get_user_sessions"
	"github.com/Anti-Raid/api/routes/auth/endpoints/revoke_user_session"
	"github.com/Anti-Raid/api/routes/auth/endpoints/test_auth"

	"github.com/anti-raid/eureka/uapi"
	"github.com/go-chi/chi/v5"
)

const tagName = "Auth"

type Router struct{}

func (r Router) Tag() (string, string) {
	return tagName, "Authentication APIs"
}

func (m Router) Routes(r *chi.Mux) {
	uapi.Route{
		Pattern: "/auth/test",
		OpId:    "test_auth",
		Method:  uapi.POST,
		Docs:    test_auth.Docs,
		Handler: test_auth.Route,
	}.Route(r)

	uapi.Route{
		Pattern: "/oauth2",
		OpId:    "create_oauth2_login",
		Method:  uapi.POST,
		Docs:    create_oauth2_login.Docs,
		Handler: create_oauth2_login.Route,
	}.Route(r)

	uapi.Route{
		Pattern: "/sessions",
		OpId:    "get_user_sessions",
		Method:  uapi.GET,
		Docs:    get_user_sessions.Docs,
		Handler: get_user_sessions.Route,
		Auth: []uapi.AuthType{
			{
				Type:         "User",
				AllowedScope: "ban_exempt",
			},
		},
	}.Route(r)

	uapi.Route{
		Pattern: "/sessions",
		OpId:    "create_user_session",
		Method:  uapi.POST,
		Docs:    create_user_session.Docs,
		Handler: create_user_session.Route,
		Auth: []uapi.AuthType{
			{
				Type: "User",
			},
		},
	}.Route(r)

	uapi.Route{
		Pattern: "/sessions/{session_id}",
		OpId:    "revoke_user_session",
		Method:  uapi.DELETE,
		Docs:    revoke_user_session.Docs,
		Handler: revoke_user_session.Route,
		Auth: []uapi.AuthType{
			{
				Type:         "User",
				AllowedScope: "ban_exempt",
			},
		},
	}.Route(r)
}
