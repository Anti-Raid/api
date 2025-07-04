package get_template_shop

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Anti-Raid/api/state"
	"github.com/Anti-Raid/api/types"
	"github.com/Anti-Raid/api/utils"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"

	docs "github.com/anti-raid/eureka/doclib"
	"github.com/anti-raid/eureka/uapi"
)

var (
	templateShopRows    = utils.GetCols(types.TemplateShopTemplate{})
	templateShopRowsStr = strings.Join(templateShopRows, ", ")
)

func Docs() *docs.Doc {
	return &docs.Doc{
		Summary:     "Get Template Shop",
		Description: "This endpoint will get the templates in the template shop with said name",
		Resp:        types.TemplateShopPartialTemplate{},
		Params: []docs.Parameter{
			{
				Name:        "name",
				Description: "The name of the template to get",
				In:          "path",
				Required:    true,
				Schema:      docs.IdSchema,
			},
		},
	}
}

func Route(d uapi.RouteData, r *http.Request) uapi.HttpResponse {
	name := chi.URLParam(r, "name")
	if name == "" {
		state.Logger.Error("Template name is empty", zap.String("name", name))
		return uapi.HttpResponse{
			Json: types.ApiError{
				Message: "Template name cannot be empty",
			},
			Status: http.StatusBadRequest,
		}
	}

	var row pgx.Rows
	var err error
	if strings.Contains(name, "@") {
		nameAndVersion := strings.SplitN(name, "@", 2)
		if len(nameAndVersion) != 2 {
			state.Logger.Error("Invalid template name format", zap.String("name", name))
			return uapi.HttpResponse{
				Json: types.ApiError{
					Message: "Invalid template name format. Expected 'name@version'",
				},
				Status: http.StatusBadRequest,
			}
		}

		name = nameAndVersion[0]
		version := nameAndVersion[1]

		if version == "latest" {
			row, err = state.Pool.Query(d.Context, "SELECT "+templateShopRowsStr+" FROM template_shop WHERE type = 'public' AND name = $1 ORDER BY version DESC LIMIT 1", name)
		} else {
			row, err = state.Pool.Query(d.Context, "SELECT "+templateShopRowsStr+" FROM template_shop WHERE type = 'public' AND name = $1 AND version = $2", name, version)
		}
	} else {
		row, err = state.Pool.Query(d.Context, "SELECT "+templateShopRowsStr+" FROM template_shop WHERE type = 'public' AND (id::text = $1 OR name = $1)", name)
	}

	if err != nil {
		state.Logger.Error("Error querying database", zap.Error(err))
		return uapi.DefaultResponse(http.StatusInternalServerError)
	}

	rows, err := pgx.CollectOneRow[types.TemplateShopTemplate](row, pgx.RowToStructByName)

	if errors.Is(err, pgx.ErrNoRows) {
		state.Logger.Warn("Template not found in shop", zap.String("name", name))
		return uapi.HttpResponse{
			Json: types.ApiError{
				Message: "Template not found in shop",
			},
			Status: http.StatusNotFound,
		}
	}

	if err != nil {
		state.Logger.Error("Error collecting rows from database", zap.Error(err))
		return uapi.DefaultResponse(http.StatusInternalServerError)
	}
	return uapi.HttpResponse{
		Json: rows,
	}
}
