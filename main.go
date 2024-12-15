package main

import (
	"fmt"
	"github.com/guneyin/printhub/cmd/api"
	"github.com/guneyin/printhub/market"
	"github.com/guneyin/printhub/utils"
	"log"
	"time"
)

// @title PrintHub API Doc
// @version 1.0
// @description Photo printing platform

// @contact.name Hüseyin Güney
// @contact.url https://github.com/guneyin
// @contact.email guneyin@gmail.com

// @BasePath /api
func main() {
	utils.SetLastRun(time.Now())

	app, err := api.NewApplication()
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(app.Server.Listen(fmt.Sprintf(":%s", market.Get().Config.ApiPort)))
}
