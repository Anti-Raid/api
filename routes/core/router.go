package core

import (
	"github.com/Anti-Raid/api/routes/core/endpoints/get_api_config"
	"github.com/Anti-Raid/api/routes/core/endpoints/get_bot_state"
	"github.com/Anti-Raid/api/routes/core/endpoints/get_bot_stats"
	"github.com/Anti-Raid/api/routes/core/endpoints/get_template_shop"
	"github.com/Anti-Raid/api/routes/core/endpoints/list_template_shop"
	"github.com/anti-raid/eureka/uapi"
	"github.com/go-chi/chi/v5"
)

const tagName = "Core"

type Router struct{}

func (b Router) Tag() (string, string) {
	return tagName, "These API endpoints are related to core functionality"
}

func (b Router) Routes(r *chi.Mux) {
	uapi.Route{
		Pattern: "/config",
		OpId:    "get_api_config",
		Method:  uapi.GET,
		Docs:    get_api_config.Docs,
		Handler: get_api_config.Route,
	}.Route(r)

	uapi.Route{
		Pattern: "/bot-state",
		OpId:    "get_modules",
		Method:  uapi.GET,
		Docs:    get_bot_state.Docs,
		Handler: get_bot_state.Route,
	}.Route(r)

	uapi.Route{
		Pattern: "/bot-stats",
		OpId:    "get_bot_stats",
		Method:  uapi.GET,
		Docs:    get_bot_stats.Docs,
		Handler: get_bot_stats.Route,
	}.Route(r)

	uapi.Route{
		Pattern: "/template-shop",
		OpId:    "list_template_shop",
		Method:  uapi.GET,
		Docs:    list_template_shop.Docs,
		Handler: list_template_shop.Route,
	}.Route(r)

	uapi.Route{
		Pattern: "/template-shop/{name}",
		OpId:    "get_template_shop",
		Method:  uapi.GET,
		Docs:    get_template_shop.Docs,
		Handler: get_template_shop.Route,
	}.Route(r)
}
