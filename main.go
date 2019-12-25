package main

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/alaingilbert/ogame"
)

var isInit bool = false
var bot, err = ogame.New("Janice", "nemesism@hotmail.fr", "pencilcho44", "fr")
var items GlobalList

func satProduction(id ogame.PlanetID) {
	pl, _ := bot.GetPlanet(id)
	//fac, _ := bot.GetResourcesBuildings(ogame.CelestialID(id))
	temp := pl.Temperature
	satprod := ogame.SolarSatellite.Production(temp, 1)
	//cenprice := 20 * math.Pow(1.1, float64(fac.SolarPlant))
	fmt.Print("sattelitte production:")
	fmt.Println(satprod)
}

func launch() {
	gal, _ := bot.GalaxyInfos(1, 337)
	var i int64
	for i = 1; i <= 15; i++ {
		pos := gal.Position(i)
		if pos != nil {
			fmt.Println(pos.Inactive)
		}
	}

	gestionrapport()

	for {
		items.planetes = bot.GetPlanets()
		items.fleets, _ = bot.GetFleets()
		items.planetinfos = nil
		i := 0
		for _, planete := range items.planetes {
			id := ogame.CelestialID(planete.ID)
			plinfo := gestionGlobal(id)
			plinfo.coord = planete.Coordinate
			items.planetinfos = append(items.planetinfos, plinfo)
			items.planeteName = planete.Name
			fmt.Println(planete.Name)
			//satProduction(planete.ID)
			if i == 0 {
				items.researchs = setresearch(id)
			}
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
	//	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
