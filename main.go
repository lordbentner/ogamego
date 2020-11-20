package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Masterminds/sprig"
	"github.com/alaingilbert/ogame"
	"github.com/go-macaron/binding"
	"gopkg.in/macaron.v1"
)

var m *macaron.Macaron

func main() {
	m = macaron.Classic()
	m.Get("/", func(ctx *macaron.Context) {
		ctx.HTML(200, "index")
	})
	m.Post("/ogame", binding.Form(Login{}), func(login Login, ctx *macaron.Context, r *http.Request) {
		if login.Universe == "" || login.User == "" || login.Password == "" {
			ctx.Redirect("/")
		} else {
			bot, err = ogame.New(login.Universe, login.User, login.Password, "fr")
			file, _ := json.Marshal(login)
			_ = ioutil.WriteFile("data.json", file, 0777)
			if err != nil {
				panic(err)
			} else {
				startLog = time.Now()
				Logout = false
				go launch()
				ctx.Redirect("/databoard")
			}
		}
	})

	m.Get("/databoard", func(ctx *macaron.Context) {
		ctx.Data["resInBuild"] = items.researchInBuild
		ctx.Data["point"] = items.points
		ctx.Data["planetinfos"] = items.planetinfos
		ctx.Data["researchs"] = items.researchs
		ctx.Data["countResInBuild"] = items.countResearchBuild
		time, user := getTimeInGame()
		ctx.Data["time_con"] = time
		ctx.Data["user"] = user
		ctx.Data["BuildLune"] = BuildLune
		ctx.HTML(200, "ogame")
	})

	m.Get("/flottes", func(ctx *macaron.Context, req *http.Request) {
		buildPage(ctx, req, 4, len(items.fleets))
		ctx.Data["flottes"] = items.fleets
		ctx.HTML(200, "flottes")
	})

	m.Get("/rapports", func(ctx *macaron.Context, req *http.Request) {
		buildPage(ctx, req, 15, len(RapportEspionnage))
		ctx.Data["spy"] = vlistAttack
		ctx.HTML(200, "rapports")
	})

	m.Get("/quit", func(ctx *macaron.Context) {
		if bot != nil {
			bot.Logout()
		}
		Logout = true
		ctx.Redirect("/")
	})

	m.Use(macaron.Renderer(macaron.RenderOptions{
		Funcs: []template.FuncMap{sprig.FuncMap()},
	}))
	http.Handle("/", m)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
