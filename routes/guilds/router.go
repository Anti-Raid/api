package guilds

import (
	"github.com/Anti-Raid/api/routes/guilds/endpoints/get_settings"
	"github.com/Anti-Raid/api/routes/guilds/endpoints/settings_execute"
	"github.com/Anti-Raid/corelib_go/splashcore"
	"github.com/anti-raid/eureka/uapi"
	"github.com/go-chi/chi/v5"
)

const tagName = "Guilds"

type Router struct{}

func (b Router) Tag() (string, string) {
	return tagName, "These API endpoints are related to AntiRaid guilds"
}

func (b Router) Routes(r *chi.Mux) {
	uapi.Route{
		Pattern: "/guilds/{guild_id}/settings",
		OpId:    "get_settings",
		Method:  uapi.GET,
		Docs:    get_settings.Docs,
		Handler: get_settings.Route,
		Auth: []uapi.AuthType{
			{
				Type: splashcore.TargetTypeUser,
			},
		},
	}.Route(r)

	uapi.Route{
		Pattern: "/guilds/{guild_id}/settings",
		OpId:    "settings_execute",
		Method:  uapi.POST,
		Docs:    settings_execute.Docs,
		Handler: settings_execute.Route,
		Auth: []uapi.AuthType{
			{
				Type: splashcore.TargetTypeUser,
			},
		},
	}.Route(r)
}
