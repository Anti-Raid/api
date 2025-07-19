package proxy

import (
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/Anti-Raid/api/rpc"
	"github.com/Anti-Raid/api/state"
	"github.com/Anti-Raid/api/types"
	docs "github.com/anti-raid/eureka/doclib"
	"github.com/anti-raid/eureka/ratelimit"
	"github.com/anti-raid/eureka/uapi"
	"go.uber.org/zap"
)

func Docs() *docs.Doc {
	return &docs.Doc{
		Summary:     "Staff Proxy",
		Description: "Proxy for staff operations, allowing for moderation and management of AntiRaid services.",
		Req:         types.SettingsExecute{},
		Resp:        map[string]any{},
		RespName:    "AnyStaffProxyResp",
		Params: []docs.Parameter{
			{
				Name:        "__service",
				Description: "The service to proxy the request to (e.g., 'template-worker')",
				Required:    true,
				In:          "query",
				Schema:      docs.IdSchema,
			},
			{
				Name:        "__method",
				Description: "The method to use for the request (e.g., 'GET', 'POST')",
				Required:    true,
				In:          "query",
				Schema:      docs.IdSchema,
			},
			{
				Name:        "__path",
				Description: "The path to use for the reques",
				Required:    true,
				In:          "query",
				Schema:      docs.IdSchema,
			},
		},
	}
}

func Route(d uapi.RouteData, r *http.Request) uapi.HttpResponse {
	if !slices.Contains(state.Config.DiscordAuth.RootUsers, d.Auth.ID) {
		return uapi.HttpResponse{
			Status: http.StatusForbidden,
			Headers: map[string]string{
				"X-Proxy-Error": "true",
			},
			Json: types.ApiError{
				Message: "You are not allowed to access this endpoint",
			},
		}
	}

	limit, err := ratelimit.Ratelimit{
		Expiry:      5 * time.Minute,
		MaxRequests: 10,
		Bucket:      "staffproxy",
	}.Limit(d.Context, r)

	if err != nil {
		state.Logger.Error("Error while ratelimiting", zap.Error(err), zap.String("bucket", "settings_execute"))
		return uapi.DefaultResponse(http.StatusInternalServerError)
	}

	if limit.Exceeded {
		return uapi.HttpResponse{
			Json: types.ApiError{
				Message: "You are being ratelimited. Please try again in " + limit.TimeToReset.String(),
			},
			Headers: limit.Headers(),
			Status:  http.StatusTooManyRequests,
		}
	}

	method := r.URL.Query().Get("__method")

	if method == "" {
		return uapi.HttpResponse{
			Status: http.StatusBadRequest,
			Headers: map[string]string{
				"X-Proxy-Error": "true",
			},
			Json: types.ApiError{
				Message: "Missing __method query parameter",
			},
		}
	}

	service := r.URL.Query().Get("__service")

	if service == "" {
		return uapi.HttpResponse{
			Status: http.StatusBadRequest,
			Headers: map[string]string{
				"X-Proxy-Error": "true",
			},
			Json: types.ApiError{
				Message: "Missing __service query parameter",
			},
		}
	}

	baseurl := ""
	switch service {
	case "template-worker":
		baseurl = rpc.CalcTWAddr()
	default:
		return uapi.HttpResponse{
			Status: http.StatusBadRequest,
			Headers: map[string]string{
				"X-Proxy-Error": "true",
			},
			Json: types.ApiError{
				Message: fmt.Sprintf("Unknown service %s", service),
			},
		}
	}

	if !r.URL.Query().Has("__path") {
		return uapi.HttpResponse{
			Status: http.StatusBadRequest,
			Headers: map[string]string{
				"X-Proxy-Error": "true",
			},
		}
	}

	queryParams := ""
	for k, v := range r.URL.Query() {
		if k != "__method" && k != "__service" && k != "__path" {
			if queryParams == "" {
				queryParams = "?" + k + "=" + strings.Join(v, ",")
			} else {
				queryParams += "&" + k + "=" + strings.Join(v, ",")
			}
		}
	}

	url := fmt.Sprintf("%s%s%s", baseurl, r.URL.Query().Get("__path"), queryParams)
	body := r.Body
	if body == nil {
		body = http.NoBody
	}

	cli := &http.Client{
		Timeout: 2 * time.Minute,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	req, err := http.NewRequest(method, url, body)

	if err != nil {
		return uapi.HttpResponse{
			Status: http.StatusInternalServerError,
			Headers: map[string]string{
				"X-Proxy-Error": "true",
			},
			Json: types.ApiError{
				Message: fmt.Sprintf("Error creating request to backend: %v", err),
			},
		}
	}

	var headers = r.Header.Clone()

	if r.Header != nil {
		// Remove headers that should not be forwarded
		for _, h := range []string{"Host", "Content-Length", "X-Forwarded-Host", "X-Forwarded-For", "X-Proxy-User"} {
			delete(headers, h)
		}

		headers.Set("X-Forwarded-Host", r.Host)
		headers.Set("X-Forwarded-For", r.RemoteAddr)
		headers.Set("X-Proxy-User", d.Auth.ID)
	}

	// Set a default Content-Type if not set
	if headers.Get("Content-Type") == "" {
		headers.Set("Content-Type", "application/json")
	}

	req.Header = headers

	resp, err := cli.Do(req)

	if err != nil {
		state.Logger.Error("Error sending request to backend", zap.Error(err), zap.String("url", url), zap.String("method", method))
		return uapi.HttpResponse{
			Status: http.StatusInternalServerError,
			Headers: map[string]string{
				"X-Proxy-Error": "true",
			},
			Json: types.ApiError{
				Message: fmt.Sprintf("Error sending request to backend: %v", err),
			},
		}
	}

	var respBody []byte
	if resp.Body != nil {
		respBody, err = io.ReadAll(resp.Body)
		if err != nil {
			state.Logger.Error("Error reading response body from backend", zap.Error(err), zap.String("url", url), zap.String("method", method))
			return uapi.HttpResponse{
				Status: http.StatusInternalServerError,
				Headers: map[string]string{
					"X-Proxy-Error": "true",
				},
				Json: types.ApiError{
					Message: fmt.Sprintf("Error reading response body from backend: %v", err),
				},
			}
		}
	}

	return uapi.HttpResponse{
		Status: resp.StatusCode,
		Bytes:  respBody,
	}
}
