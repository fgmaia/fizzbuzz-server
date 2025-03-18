package apps

import (
	"fizzbuzz-server/internal/apps/contracts"
	"fizzbuzz-server/internal/services"
	"sync"

	"github.com/gofiber/fiber/v2"
)

var (
	once sync.Once
	app  *FizzbuzzApp
)

var App = func() *FizzbuzzApp {
	once.Do(func() {
		app = &FizzbuzzApp{}
		app.init()
	})
	return app
}

type FizzbuzzApp struct {
	FiberApp        *fiber.App
	FizzBuzzService contracts.FizzBuzzServiceIface
	StatsService    contracts.StatsServiceIface
}

func (f *FizzbuzzApp) init() {
	//ulog.LogInit(logFileName)
	//ulog.Info("TupanWegGo - Conexão:", "connection string", config.OracleConnectionString())

	f.FizzBuzzService = services.NewFizzBuzzService()
	f.StatsService = services.NewStatsService()
}
