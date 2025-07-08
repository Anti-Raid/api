package staff

import (
	"github.com/Anti-Raid/api/routes/staff/endpoints/proxy"
	"github.com/Anti-Raid/corelib_go/splashcore"
	"github.com/anti-raid/eureka/uapi"
	"github.com/go-chi/chi/v5"
)

const tagName = "Staff"

type Router struct{}

func (b Router) Tag() (string, string) {
	return tagName, "These API endpoints are AntiRaid staff only"
}

func (b Router) Routes(r *chi.Mux) {
	uapi.Route{
		Pattern: "/staff/proxy",
		OpId:    "proxy",
		Method:  uapi.POST,
		Docs:    proxy.Docs,
		Handler: proxy.Route,
		Auth: []uapi.AuthType{
			{
				Type: splashcore.TargetTypeUser,
			},
		},
	}.Route(r)
}
