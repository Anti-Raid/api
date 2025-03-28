package ioauth_download_job

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"time"

	api "github.com/Anti-Raid/api/auth"
	"github.com/Anti-Raid/api/state"
	"github.com/Anti-Raid/api/types"
	"github.com/Anti-Raid/corelib_go/objectstorage"

	"github.com/go-chi/chi/v5"
	docs "github.com/infinitybotlist/eureka/doclib"
	"github.com/infinitybotlist/eureka/uapi"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

var downloadTemplate = template.Must(template.New("download").Parse(`<!DOCTYPE html>
<html>
	Your download should start in a moment. If not, <a href="{{.URL}}">click here</a>
	<script>
		if(window.opener) {
			window.opener.postMessage("dl:{{.URL}}", {{.Domain}});
		} else if(window.parent) {
			window.parent.postMessage("dl::{{.URL}}", {{.Domain}});
		}
		window.location.href = "{{.URL}}";
	</script>
</html>
`))

func Docs() *docs.Doc {
	return &docs.Doc{
		Summary:     "Get IOAuth Download Link",
		Description: "Gets the download link to a jobs output.",
		Params: []docs.Parameter{
			{
				Name:        "id",
				Description: "Task ID",
				Required:    true,
				In:          "path",
				Schema:      docs.IdSchema,
			},
			{
				Name:        "no_redirect",
				Description: "Whether or not to avoid the redirect/text response and merely return the link",
				Required:    true,
				In:          "query",
				Schema:      docs.IdSchema,
			},
		},
		Resp: "URL",
	}
}

var permList = map[string]string{
	"message_prune":       "moderation.prune",
	"guild_create_backup": "server_backups.download",
}

func Route(d uapi.RouteData, r *http.Request) uapi.HttpResponse {
	iotok := r.URL.Query().Get("ioauth")

	if iotok == "" {
		return uapi.HttpResponse{
			Status: http.StatusBadRequest,
			Json: types.ApiError{
				Message: "Missing IOAuth token",
			},
		}
	}

	id := chi.URLParam(r, "id")

	if id == "" {
		return uapi.HttpResponse{
			Status: http.StatusBadRequest,
			Json:   types.ApiError{Message: "`id` is required"},
		}
	}

	// Get the ioauth token
	resp, err := state.Rueidis.Do(d.Context, state.Rueidis.B().Get().Key("ioauth:{"+iotok+"}").Build()).AsBytes()

	if err != nil {
		return uapi.HttpResponse{
			Status: http.StatusInternalServerError,
			Json: types.ApiError{
				Message: "Internal Server Error [while checking ioauth token]: " + err.Error(),
			},
		}
	}

	if resp == nil {
		return uapi.HttpResponse{
			Status: http.StatusBadRequest,
			Json: types.ApiError{
				Message: "Invalid IOAuth token",
			},
		}
	}

	var iot types.IOAuthOutput

	err = json.Unmarshal(resp, &iot)

	if err != nil {
		return uapi.HttpResponse{
			Status: http.StatusInternalServerError,
			Json: types.ApiError{
				Message: "Internal Server Error [while parsing ioauth token]: " + err.Error(),
			},
		}
	}

	// Delete expired jobs first
	_, err = state.Pool.Exec(d.Context, "DELETE FROM jobs WHERE created_at + expiry < NOW()")

	if err != nil {
		state.Logger.Error("Failed to delete expired jobs [db delete]", zap.Error(err))
		return uapi.DefaultResponse(http.StatusInternalServerError)
	}

	var name string
	var guildId string
	var output *map[string]any
	err = state.Pool.QueryRow(d.Context, "SELECT name, guild_id, output FROM jobs WHERE id = $1", id).Scan(&name, &guildId, &output)

	if errors.Is(err, pgx.ErrNoRows) {
		return uapi.HttpResponse{
			Status: http.StatusNotFound,
			Json:   types.ApiError{Message: "Job not found"},
		}
	}

	if err != nil {
		state.Logger.Error("Failed to fetch jobs [db fetch]", zap.Error(err))
		return uapi.DefaultResponse(http.StatusInternalServerError)
	}

	// Check user permissions based on permMap
	if perm, ok := permList[name]; ok {
		hresp, ok := api.HandlePermissionCheck(iot.DiscordUser.ID, guildId, perm)

		if !ok {
			return hresp
		}
	} else {
		return uapi.HttpResponse{
			Status: http.StatusNotFound,
			Json:   types.ApiError{Message: "This task cannot be downloaded using the API at this time: " + name},
		}
	}

	if output == nil {
		return uapi.HttpResponse{
			Status: http.StatusNotFound,
			Json:   types.ApiError{Message: "Task output not found"},
		}
	}

	var outputMap = *output

	filename, ok := outputMap["filename"].(string)

	if !ok {
		return uapi.HttpResponse{
			Status: http.StatusNotFound,
			Json:   types.ApiError{Message: "Task output filename not found"},
		}
	}

	dir := fmt.Sprintf("jobs/%s", id)

	// Now get URL
	url, err := state.ObjectStorage.GetUrl(d.Context, objectstorage.GuildBucket(guildId), dir, filename, 10*time.Minute)

	if err != nil {
		state.Logger.Error("Failed to get url for job", zap.Error(err))
		return uapi.HttpResponse{
			Status: http.StatusInternalServerError,
			Json:   types.ApiError{Message: "Failed to get url for job: " + err.Error()},
		}
	}

	if r.URL.Query().Get("no_redirect") == "true" {
		return uapi.HttpResponse{
			Status: http.StatusOK,
			Json: types.ApiError{
				Message: url.String(),
			},
		}
	} else {
		var buf bytes.Buffer
		err := downloadTemplate.Execute(&buf, map[string]any{
			"URL":    url.String(),
			"Domain": state.Config.Sites.Frontend,
		})

		if err != nil {
			state.Logger.Error("Failed to execute download template", zap.Error(err))
			return uapi.HttpResponse{
				Status: http.StatusInternalServerError,
				Json:   types.ApiError{Message: "Failed to execute download template: " + err.Error()},
			}
		}

		return uapi.HttpResponse{
			Status: http.StatusFound,
			Bytes:  buf.Bytes(),
			Headers: map[string]string{
				"Content-Type": "text/html, charset=utf-8",
			},
		}
	}
}
