package main

import (
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
	"text/template"

	"github.com/alaingilbert/ogame"
	"github.com/fatih/structs"
	"github.com/go-macaron/binding"
	"gopkg.in/macaron.v1"
)

var (
	bot *ogame.OGame
)
var (
	err error
)
var items GlobalList
var RapportEspionnage []map[string]interface{}

func satProduction(id ogame.PlanetID) {
	pl, _ := bot.GetPlanet(id)
	fac, _ := bot.GetResourcesBuildings(ogame.CelestialID(id))
	temp := pl.Temperature
	satprod := ogame.SolarSatellite.Production(temp, 1, true)
	cenprice := 20 * math.Pow(1.1, float64(fac.SolarPlant))
	if cenprice > float64(satprod*2000) {
		pid := ogame.CelestialID(id)
		bot.BuildShips(pid, ogame.SolarSatelliteID, 1)
	}
}

func launch() {
	var gal int64 = 3
	var sys int64 = 19
	for {
		items.planetes = bot.GetPlanets()
		fl, _ := bot.GetFleets()
		stres := bot.GetResearch()
		items.planetinfos = nil
		i := 0
		if len(items.planetes) > len(items.facilities) {
			items.facilities = make([]map[string]interface{}, len(items.planetes))
			items.resources = make([]map[string]interface{}, len(items.planetes))
			items.ships = make([]map[string]interface{}, len(items.planetes))
			items.res_build = make([]map[string]interface{}, len(items.planetes))
			items.consInBuild = make([]ogame.ID, len(items.planetes))
			items.countInBuild = make([]string, len(items.planetes))
		}

		if len(items.fleets) < len(fl) {
			items.fleets = make([]map[string]interface{}, len(fl))
		}

		for j, fle := range fl {
			items.fleets[j] = structs.Map(fle)
		}

		for _, planete := range items.planetes {
			id := ogame.CelestialID(planete.ID)
			gestionUnderAttack(id)
			plinfo := gestionGlobal(id)
			items.facilities[i] = plinfo.facilities
			items.resources[i] = plinfo.resources
			items.res_build[i] = plinfo.res_build
			items.consInBuild[i] = plinfo.consInBuild
			items.countInBuild[i] = plinfo.countInBuild
			satProduction(planete.ID)
			inter := stres.IntergalacticResearchNetwork
			if i < int(inter) {
				items.researchs = setresearch(id)
			} else {
				transporter(id, items.planetes[0].Coordinate)
			}

			items.ships[i] = setShips(id)
			if sys >= 500 {
				sys = 1
				gal++
			}
			if gal >= 5 {
				gal = 1
			}

			comput := items.researchs["Computer"].(int64)
			if len(fl) < int(comput) {
				gestionEspionnage(id, gal, sys)
				gestionrapport(id)
				sys++
			}

			setExpedition(id, planete.Coordinate)

			i++
		}
	}
}

func tem(letter string) {
	t := template.Must(template.New("letter").Parse(letter))
	fmt.Println(t)
}

func buildPage(ctx *macaron.Context, req *http.Request, elInPage int, nbItem int) {
	page, ok := req.URL.Query()["page"]
	ctx.Data["fElem"] = 0
	ctx.Data["lElem"] = elInPage
	ctx.Data["Page"] = 1
	if ok && len(page[0]) > 0 {
		el, _ := strconv.Atoi(page[0])
		ctx.Data["fElem"] = (el - 1) * elInPage
		ctx.Data["lElem"] = el * elInPage
		ctx.Data["Page"] = el
	}
	var nbPages []int
	for i := 0; i < nbItem/elInPage+1; i++ {
		nbPages = append(nbPages, i+1)
	}

	ctx.Data["nbPages"] = nbPages
}

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
			go launch()
			ctx.Redirect("/databoard")
		}
	})
	m.Get("/databoard", func(ctx *macaron.Context) {
		MapItems := structs.Map(items)
		ctx.Data["items"] = MapItems
		ctx.Data["planetes"] = items.planetes
		ctx.Data["researchs"] = items.researchs
		ctx.Data["facilities"] = items.facilities
		ctx.Data["resources"] = items.resources
		ctx.Data["res_build"] = items.res_build
		ctx.Data["ships"] = items.ships
		ctx.Data["consInBuild"] = items.consInBuild
		ctx.Data["countInBuild"] = items.countInBuild
		ctx.Data["resInBuild"] = items.researchInBuild
		ctx.Data["countResInBuild"] = items.countResearchBuild
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
