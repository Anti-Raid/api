package list_template_shop

import (
	"net/http"
	"strings"

	"github.com/Anti-Raid/api/state"
	"github.com/Anti-Raid/api/types"
	"github.com/Anti-Raid/api/utils"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"

	docs "github.com/anti-raid/eureka/doclib"
	"github.com/anti-raid/eureka/uapi"
)

var (
	templateShopPartialRows    = utils.GetCols(types.TemplateShopPartialTemplate{})
	templateShopPartialRowsStr = strings.Join(templateShopPartialRows, ", ")
)

func Docs() *docs.Doc {
	return &docs.Doc{
		Summary:     "List Template Shop",
		Description: "This endpoint will list all templates in the template shop that are public.",
		Resp:        types.TemplateShopPartialTemplate{},
		Params:      []docs.Parameter{},
	}
}

func Route(d uapi.RouteData, r *http.Request) uapi.HttpResponse {
	row, err := state.Pool.Query(d.Context, "SELECT "+templateShopPartialRowsStr+" FROM template_shop WHERE type = 'public'")

	if err != nil {
		state.Logger.Error("Error querying database", zap.Error(err))
		return uapi.DefaultResponse(http.StatusInternalServerError)
	}

	rows, err := pgx.CollectRows[types.TemplateShopPartialTemplate](row, pgx.RowToStructByName)
	if err != nil {
		state.Logger.Error("Error collecting rows from database", zap.Error(err))
		return uapi.DefaultResponse(http.StatusInternalServerError)
	}
	return uapi.HttpResponse{
		Json: rows,
	}
}
