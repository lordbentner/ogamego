package main

import (
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"

	//"strconv"
	//"strings"

	"text/template"

	"github.com/alaingilbert/ogame"
	"github.com/go-macaron/binding"
	"gopkg.in/macaron.v1"
	//"golang.org/x/mobile/app"
)

//var bot, err = ogame.New("Aquarius", os.Args[1], os.Args[2], "fr")
var (
	bot *ogame.OGame
)
var (
	err error
)
var items GlobalList

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
	//P:1:360:6
	fmt.Print("sattelitte production:")
	fmt.Println(satprod)
}

func launch() {
	var gal int64 = 1
	var sys int64 = 1
	for {
		items.planetes = bot.GetPlanets()
		items.fleets, _ = bot.GetFleets()
		items.planetinfos = nil
		i := 0
		if len(items.planetes) > len(items.facilities) {
			items.facilities = make([]map[string]interface{}, len(items.planetes))
			items.resources = make([]map[string]interface{}, len(items.planetes))
			items.ships = make([]map[string]interface{}, len(items.planetes))
			items.res_build = make([]map[string]interface{}, len(items.planetes))
			items.consInBuild = make([]string, len(items.planetes))
			items.countInBuild = make([]int64, len(items.planetes))
		}

		for _, planete := range items.planetes {
			fmt.Println("Nom de la planéte:", planete.Name)
			id := ogame.CelestialID(planete.ID)
			gestionUnderAttack(id)
			plinfo := gestionGlobal(id)
			items.facilities[i] = plinfo.facilities
			items.resources[i] = plinfo.resources
			items.res_build[i] = plinfo.res_build
			items.consInBuild[i] = plinfo.consInBuild
			items.countInBuild[i] = plinfo.countInBuild
			satProduction(planete.ID)
			if i == 0 {
				items.researchs = setresearch(id)
			}

			items.ships[i] = setShips(id)
			if sys >= 500 {
				sys = 1
				gal++
			}
			if gal >= 5 {
				gal = 1
			}

			setExpedition(id, planete.Coordinate)
			//	gestionEspionnage(id, gal, sys)
			gestionrapport(id)
			//gestionAttack(id)
			sys++
			i++
		}
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("index.html")
	http.StripPrefix("/static/", http.FileServer(http.Dir("static")))
	t.Execute(w, items)
}

func main() {
	//app.Main(func(a app.App) {
	m := macaron.Classic()
	m.Use(macaron.Renderer())
	m.Get("/", func(ctx *macaron.Context) {
		ctx.HTML(200, "index")
	})
	m.Post("/ogame", binding.Form(Login{}), func(login Login, ctx *macaron.Context) {
		fmt.Println("User:", login.User, "Password:", login.Password)
		bot, err = ogame.New("Aquarius", login.User, login.Password, "fr")
		go launch()
		ctx.Redirect("/databoard")
	})
	m.Get("/databoard", func(ctx *macaron.Context) {
		ctx.Data["items"] = items
		ctx.Data["planetes"] = items.planetes
		ctx.Data["researchs"] = items.researchs
		ctx.Data["facilities"] = items.facilities
		ctx.Data["resources"] = items.resources
		ctx.Data["res_build"] = items.res_build
		ctx.Data["ships"] = items.ships
		ctx.Data["consInBuild"] = items.consInBuild
		ctx.Data["countInBuild"] = items.countInBuild
		ctx.HTML(200, "ogame")
	})

	host := os.Getenv("IP")
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	fmt.Println("host:", host, "PORT:", port)
	m.Run("127.0.0.1", "8000")
	//})
}
