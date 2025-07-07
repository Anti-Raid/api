package users

import (
	"github.com/Anti-Raid/api/routes/users/endpoints/get_user"
	"github.com/Anti-Raid/api/routes/users/endpoints/get_user_guild_base_info"
	"github.com/Anti-Raid/api/routes/users/endpoints/get_user_guilds"
	"github.com/anti-raid/eureka/uapi"
	"github.com/go-chi/chi/v5"
)

const tagName = "Users"

type Router struct{}

func (b Router) Tag() (string, string) {
	return tagName, "These API endpoints are related to AntiRaid users"
}

func (b Router) Routes(r *chi.Mux) {
	uapi.Route{
		Pattern: "/users/{id}",
		OpId:    "get_user",
		Method:  uapi.GET,
		Docs:    get_user.Docs,
		Handler: get_user.Route,
	}.Route(r)

	uapi.Route{
		Pattern: "/users/@me/guilds",
		OpId:    "get_user_guilds",
		Method:  uapi.GET,
		Docs:    get_user_guilds.Docs,
		Handler: get_user_guilds.Route,
		Auth: []uapi.AuthType{
			{
				Type: "User",
			},
		},
	}.Route(r)

	uapi.Route{
		Pattern: "/users/@me/guilds/{guild_id}",
		OpId:    "get_user_guild_base_info",
		Method:  uapi.GET,
		Docs:    get_user_guild_base_info.Docs,
		Handler: get_user_guild_base_info.Route,
		Auth: []uapi.AuthType{
			{
				Type: "User",
			},
		},
	}.Route(r)
}
