package service

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
	"github.com/uptrace/bunrouter"
	"github.com/uptrace/bunrouter/extra/reqlog"

	"github.com/SAMBA-Research/microservice-template/internal/config"
	"github.com/SAMBA-Research/microservice-template/internal/tracing"
	"github.com/SAMBA-Research/microservice-template/version"
)

type Microservice struct {
	cfg *config.Config
}

func NewMicroservice(cfg *config.Config, db *bun.DB) (srv *Microservice, err error) {
	srv = &Microservice{
		cfg: cfg,
	}
	return
}

func (srv *Microservice) Run() {
	router := bunrouter.New(
		bunrouter.Use(reqlog.NewMiddleware()),
	)

	router.GET("/", srv.indexHandler)
	router.GET("/v1", srv.indexHandler)
	router.GET("/v1/example", srv.exampleHandler)

	log.Info().Msgf("Microservice %s listening on %s:%d", version.ServiceName, srv.cfg.ServiceBind, srv.cfg.ServicePort)
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", srv.cfg.ServiceBind, srv.cfg.ServicePort), router)
	if err != nil {
		panic(err)
	}
}

func (srv *Microservice) indexHandler(w http.ResponseWriter, r bunrouter.Request) (err error) {
	_, span := tracing.Tracer().Start(r.Context(), "service.indexHandler")
	defer span.End()

	w.Write([]byte("This is an API server"))
	return
}

func (srv *Microservice) exampleHandler(w http.ResponseWriter, r bunrouter.Request) (err error) {
	_, span := tracing.Tracer().Start(r.Context(), "service.exampleHandler")
	defer span.End()

	return bunrouter.JSON(w, map[string]any{
		"ok": true,
	})
}
