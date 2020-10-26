package main

//https://ogamebot.uc.r.appspot.com
import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alaingilbert/ogame"
	"github.com/fatih/structs"
	"github.com/go-macaron/binding"
	"gopkg.in/macaron.v1"
)

func main() {
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
		MapItems := structs.Map(items)
		ctx.Data["items"] = MapItems
		ctx.Data["lunes"] = items.lunes
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
		ctx.Data["productions"] = items.productions
		time, user := getTimeInGame()
		ctx.Data["time_con"] = time
		ctx.Data["user"] = user
		ctx.Data["point"] = items.points
		ctx.Data["BuildLune"] = BuildLune
		ctx.HTML(200, "ogame")
	})

	m.Get("/flottes", func(ctx *macaron.Context, req *http.Request) {
		buildPage(ctx, req, 4, len(items.fleets))
		ctx.Data["flottes"] = items.fleets
		ctx.Data["point"] = items.points
		ctx.HTML(200, "flottes")
	})

	m.Get("/rapports", func(ctx *macaron.Context, req *http.Request) {
		buildPage(ctx, req, 15, len(RapportEspionnage))
		ctx.Data["spy"] = vlistAttack
		ctx.Data["point"] = items.points
		ctx.HTML(200, "rapports")
	})

	m.Get("/quit", func(ctx *macaron.Context) {
		bot.Logout()
		Logout = true
		ctx.Redirect("/")
	})

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
