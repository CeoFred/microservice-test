package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
	"github.com/uptrace/bunrouter"
	"github.com/uptrace/bunrouter/extra/reqlog"

	"github.com/SAMBA-Research/microservice-template/internal/config"
	"github.com/SAMBA-Research/microservice-template/internal/db"
	"github.com/SAMBA-Research/microservice-template/version"
)

type Microservice struct {
	cfg *config.Config
	db  *bun.DB
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

	ctx := context.Background()

	if err := req.ParseForm(); err != nil {
		w.WriteHeader(400)
		return bunrouter.JSON(w, bunrouter.H{
			"error":   err.Error(),
			"success": false,
		})
	}

	m := req.PostForm.Get("message")

	message := &db.Message{Message: m}
	_, err := srv.db.NewInsert().Model(message).Exec(ctx)

	if err != nil {
				w.WriteHeader(400)
		return bunrouter.JSON(w, bunrouter.H{
			"error":   err.Error(),
			"success": false,
		})
	}

	if m == "" {
		w.WriteHeader(400)
		return bunrouter.JSON(w, bunrouter.H{
			"error":   "provide a message",
			"success": false,
		})
	}


	if err != nil {
				w.WriteHeader(400)
		return bunrouter.JSON(w, bunrouter.H{
			"error":   err.Error(),
			"success": false,
		})
	}

	bunrouter.JSON(w, bunrouter.H{
		"success": true,
		"data":    nil,
	})
	return nil
}
