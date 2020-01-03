package main

import (
	"fmt"
	"math"
	"net/http"
	"os"
	"text/template"
	//"mime"
	"github.com/alaingilbert/ogame"
)

var isInit bool = false
var bot, err = ogame.New("Janice", "nemesism@hotmail.fr", os.Args[1], "fr")
var items GlobalList

func satProduction(id ogame.PlanetID) {
	pl, _ := bot.GetPlanet(id)
	fac, _ := bot.GetResourcesBuildings(ogame.CelestialID(id))
	temp := pl.Temperature
	satprod := ogame.SolarSatellite.Production(temp, 1)
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
		for _, planete := range items.planetes {
			id := ogame.CelestialID(planete.ID)
			gestionUnderAttack(id)
			plinfo := gestionGlobal(id)
			plinfo.coord = planete.Coordinate
			items.planetinfos = append(items.planetinfos, plinfo)
			items.planeteName = planete.Name
			fmt.Println(planete.Name)
			satProduction(planete.ID)
			if i == 0 {
				items.researchs = setresearch(id)
			}

			setShips(id)
			if sys >= 500 {
				sys = 1
				gal++
			}
			if gal >= 5 {
				gal = 1
			}

			setExpedition(id, plinfo.coord)
			gestionEspionnage(id, gal, sys)
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
	go launch()
	http.HandleFunc("/", handler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
