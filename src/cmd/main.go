package main

import (
	"context"
	"flag"
	"io"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/SAMBA-Research/microservice-template/internal/config"
	"github.com/SAMBA-Research/microservice-template/internal/db"
	"github.com/SAMBA-Research/microservice-template/internal/utils"
	"github.com/SAMBA-Research/microservice-template/service"
	"github.com/SAMBA-Research/microservice-template/version"
)

const (
	eventQuit = iota
)

type sysEventMessage struct {
	event int
	idata int
}

var sysEventChannel = make(chan sysEventMessage, 5)
var logOutput io.Writer
var startTime time.Time

var logFileName = flag.String("log", "-", "Log file ('-' for only stderr)")

func main() {
	os.Setenv("TZ", "UTC")
	startTime = time.Now()
	rand.Seed(startTime.UnixNano())

	defaultCtx := context.Background()

	cfg, err := config.InitConfig()
	if err != nil {
		panic(err)
	}
	if cfg.LogFileName != "" {
		*logFileName = cfg.LogFileName
	}

	flag.Parse()

	if *logFileName != "-" {
		f, err := os.OpenFile(*logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
		if err != nil {
			log.Fatal().Msg("Cannot open log file " + *logFileName)
		}
		defer f.Close()
		logOutput = io.MultiWriter(os.Stderr, f)
	} else {
		logOutput = os.Stderr
	}
	log.Logger = zerolog.New(logOutput).With().Timestamp().Logger()

	log.Info().Msg("Starting up...")

	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, syscall.SIGINT)

	otelShutdown, err := setupOTelSDK(defaultCtx, version.ServiceName, version.ServiceVersion, cfg)
	if err != nil {
		panic(err)
	}
	defer otelShutdown(defaultCtx)

	db := db.NewDbConnection(cfg)

	msServer, err := service.NewMicroservice(cfg, db)
	if err != nil {
		panic(err)
	}
	go msServer.Run()

	//go webServer()
	//go infraWebServer()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	oldAlloc := int64(m.Alloc)
	printMemStats(&m)

	for {
		select {
		case msg := <-sysEventChannel:
			switch msg.event {
			case eventQuit:
				log.Warn().Msg("Exiting")
				os.Exit(msg.idata)
			}
		case sig := <-sigChannel:
			switch sig {
			case syscall.SIGINT:
				sysEventChannel <- sysEventMessage{event: eventQuit, idata: 0}
				log.Warn().Msg("^C detected")
			}
		case <-time.After(60 * time.Second):

			runtime.ReadMemStats(&m)
			if utils.Abs(int64(m.Alloc)-oldAlloc) > 1024*1024 {
				printMemStats(&m)
				oldAlloc = int64(m.Alloc)
			}
		case <-time.After(15 * time.Minute):
			//cleanupDb()
		}
	}
}

func printMemStats(m *runtime.MemStats) {
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	log.Info().Msgf("Alloc: %v MiB\tTotalAlloc: %v MiB\tSys: %v MiB\tNumGC: %v\tUptime: %0.1fh\n",
		utils.BToMB(m.Alloc), utils.BToMB(m.TotalAlloc), utils.BToMB(m.Sys), m.NumGC, time.Since(startTime).Hours())
}
