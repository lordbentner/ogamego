package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/alaingilbert/ogame"
	"github.com/fatih/structs"
	"github.com/go-macaron/binding"
	"gopkg.in/macaron.v1"
)

func main() {
	// app.Main(func(a app.App) {
	m := macaron.Classic()
	m.Use(macaron.Renderer())
	m.Get("/", func(ctx *macaron.Context) {
		ctx.HTML(200, "index")
	})
	m.Post("/ogame", binding.Form(Login{}), func(login Login, ctx *macaron.Context) {
		if login.Universe == "" || login.User == "" || login.Password == "" {
			ctx.Redirect("/")
		} else {
			bot, err = ogame.New(login.Universe, login.User, login.Password, "fr")
			if err != nil {
				panic(err)
			} else {
				startLog = time.Now()
				go launch()
				ctx.Redirect("/databoard")
			}
		}
	})
	m.Get("/databoard", func(ctx *macaron.Context) {
		MapItems := structs.Map(items)
		ctx.Data["items"] = MapItems
		ctx.Data["planetes"] = items.planetes
		ctx.Data["researchs"] = items.researchs
		ctx.Data["facilities"] = items.facilities
		ctx.Data["resources"] = items.resources
		ctx.Data["resdetails"] = items.detailsRessources
		ctx.Data["res_build"] = items.res_build
		ctx.Data["ships"] = items.ships
		ctx.Data["consInBuild"] = items.consInBuild
		ctx.Data["countInBuild"] = items.countInBuild
		ctx.Data["resInBuild"] = items.researchInBuild
		ctx.Data["countResInBuild"] = items.countResearchBuild
		ctx.Data["time_con"] = getTimeInGame()
		ctx.HTML(200, "ogame")
	})

	m.Get("/flottes", func(ctx *macaron.Context, req *http.Request) {
		buildPage(ctx, req, 4, len(items.fleets))
		ctx.Data["flottes"] = items.fleets
		ctx.HTML(200, "flottes")
	})

	m.Get("/rapports", func(ctx *macaron.Context, req *http.Request) {
		buildPage(ctx, req, 15, len(RapportEspionnage))
		ctx.Data["spy"] = RapportEspionnage
		ctx.HTML(200, "rapports")
	})

	host := os.Getenv("IP")
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	fmt.Println("host:", host, "PORT:", port)
	if len(host) < 1 {
		host = "127.0.0.1"
		port = 8000
	}

	m.Run(host, port)
	// })
}
