package sailor

import (
	"fmt"
	"net/http"

	chi "github.com/go-chi/chi/v5"
)

type Api struct {
	Address string
	Port    int
	Worker  *Sailor
	Router  *chi.Mux
}

type ErrResponse struct {
	HTTPStatusCode int
	Message        string
}

func (a *Api) initRouter() {
	a.Router = chi.NewRouter()
	a.Router.Route("/tasks", func(r chi.Router) {
		r.Post("/", a.StartTaskHandler)
		r.Get("/", a.GetTasksHandler)
		r.Route("/{taskID}", func(r chi.Router) {
			r.Delete("/", a.StopTaskHandler)
			r.Get("/", a.InspectTaskHandler)
		})
	})
	a.Router.Route("/stats", func(r chi.Router) {
		r.Get("/", a.GetStatsHandler)
	})
}

func (a *Api) Start() {
	fmt.Print("api started")
	a.initRouter()
	http.ListenAndServe(fmt.Sprintf("%s:%d", a.Address, a.Port), a.Router)
}
