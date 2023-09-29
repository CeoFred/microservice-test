package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
	"github.com/uptrace/bunrouter"
	"github.com/uptrace/bunrouter/extra/reqlog"

	"github.com/SAMBA-Research/microservice-template/internal/config"
	"github.com/SAMBA-Research/microservice-template/internal/db"
	"github.com/SAMBA-Research/microservice-template/internal/tracing"
	"github.com/SAMBA-Research/microservice-template/version"
)

type Microservice struct {
	cfg *config.Config
	db  *bun.DB
}

type MessageBody struct {
	Message string `json:"message"`
}

func NewMicroservice(cfg *config.Config, db *bun.DB) (srv *Microservice, err error) {
	srv = &Microservice{
		cfg: cfg,
		db:  db,
	}
	return
}

func (srv *Microservice) Run() {
	router := bunrouter.New(
		bunrouter.Use(reqlog.NewMiddleware()),
	)

	router.POST("/data", srv.handleData)
	router.GET("/data", srv.retrieveMessages)

	log.Info().Msgf("Microservice %s listening on %s:%d", version.ServiceName, srv.cfg.ServiceBind, srv.cfg.ServicePort)
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", srv.cfg.ServiceBind, srv.cfg.ServicePort), router)
	if err != nil {
		panic(err)
	}
}

func (srv *Microservice) retrieveMessages(w http.ResponseWriter, req bunrouter.Request) error {

	_, span := tracing.Tracer().Start(req.Context(), "service.retrieveMessages")
	defer span.End()
	ctx := context.Background()

	var messages []db.Message
	err := srv.db.NewSelect().Model(&messages).Scan(ctx)

	if err != nil {
		w.WriteHeader(400)
		w.Header().Add("Content-Type", "application/json")
		return bunrouter.JSON(w, bunrouter.H{
			"error":   err.Error(),
			"success": false,
		})
	}

	return bunrouter.JSON(w, bunrouter.H{
		"success": true,
		"data":    messages,
	})
}

func (srv *Microservice) handleData(w http.ResponseWriter, req bunrouter.Request) error {

	_, span := tracing.Tracer().Start(req.Context(), "service.handleData")
	defer span.End()
	var body MessageBody

	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&body); err != nil {

		return bunrouter.JSON(w, bunrouter.H{
			"error":   err.Error(),
			"success": false,
		})
	}
	ctx := context.Background()

	m := body.Message

	message := &db.Message{Message: m}
	_, err := srv.db.NewInsert().Model(message).Exec(ctx)

	if err != nil {
		return bunrouter.JSON(w, bunrouter.H{
			"error":   err.Error(),
			"success": false,
		})
	}

	if m == "" {
		return bunrouter.JSON(w, bunrouter.H{
			"error":   "provide a message",
			"success": false,
		})
	}

	if err != nil {
		return bunrouter.JSON(w, bunrouter.H{
			"error":   err.Error(),
			"success": false,
		})
	}

	return bunrouter.JSON(w, bunrouter.H{
		"success": true,
		"data":    nil,
	})
}
