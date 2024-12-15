package jobs

import (
	"github.com/Anti-Raid/api/routes/jobs/endpoints/ioauth_download_job"

	"github.com/go-chi/chi/v5"
	"github.com/infinitybotlist/eureka/uapi"
)

const tagName = "Jobs"

type Router struct{}

func (b Router) Tag() (string, string) {
	return tagName, "These API endpoints are related to jobs"
}

func (b Router) Routes(r *chi.Mux) {
	uapi.Route{
		Pattern: "/jobs/{id}/ioauth/download-link",
		OpId:    "ioauth_download_job",
		Method:  uapi.GET,
		Docs:    ioauth_download_job.Docs,
		Handler: ioauth_download_job.Route,
		ExtData: map[string]any{
			"ioauth": []string{"identify"},
		},
	}.Route(r)
}
