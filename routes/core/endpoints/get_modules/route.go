package get_modules

import (
	"net/http"

	"github.com/Anti-Raid/api/rpc"
	"github.com/Anti-Raid/api/state"
	"github.com/Anti-Raid/api/types"
	"github.com/Anti-Raid/corelib_go/silverpelt"

	docs "github.com/infinitybotlist/eureka/doclib"
	"github.com/infinitybotlist/eureka/uapi"
)

func Docs() *docs.Doc {
	return &docs.Doc{
		Summary:     "Get Modules",
		Description: "This endpoint returns the modules on AntiRaid.",
		Resp:        []silverpelt.CanonicalModule{},
		Params:      []docs.Parameter{},
	}
}

var ModulesCache *[]silverpelt.CanonicalModule

func Route(d uapi.RouteData, r *http.Request) uapi.HttpResponse {
	if ModulesCache != nil {
		return uapi.HttpResponse{
			Json: ModulesCache,
		}
	}

	modules, err := rpc.Modules(state.Context)

	if err != nil {
		return uapi.HttpResponse{
			Status: http.StatusInternalServerError,
			Json:   types.ApiError{Message: "Error getting modules:" + err.Error()},
		}
	}

	ModulesCache = modules

	return uapi.HttpResponse{
		Json: modules,
	}
}
