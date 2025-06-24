package guilds

import (
	api "github.com/Anti-Raid/api/auth"
	"github.com/Anti-Raid/api/routes/guilds/endpoints/get_staff_team"
	"github.com/Anti-Raid/api/routes/guilds/endpoints/settings_execute"
	"github.com/Anti-Raid/corelib_go/splashcore"
	"github.com/go-chi/chi/v5"
	"github.com/anti-raid/eureka/uapi"
)

const tagName = "Guilds"

type Router struct{}

func (b Router) Tag() (string, string) {
	return tagName, "These API endpoints are related to AntiRaid guilds"
}

func (b Router) Routes(r *chi.Mux) {
	uapi.Route{
		Pattern: "/guilds/{guild_id}/staff-team",
		OpId:    "get_staff_team",
		Method:  uapi.GET,
		Docs:    get_staff_team.Docs,
		Handler: get_staff_team.Route,
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
		ExtData: map[string]any{
			api.PERMISSION_CHECK_KEY: nil, // Authz is performed in the handler itself
		},
	}.Route(r)
}
