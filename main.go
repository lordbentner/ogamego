package main

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/alaingilbert/ogame"
	"github.com/fatih/structs"
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

func setresearch(id ogame.CelestialID) map[string]interface{} {
	bot.BuildTechnology(id, ogame.AstrophysicsID)
	bot.BuildTechnology(id, ogame.IntergalacticResearchNetworkID)
	res := bot.GetResearch()
	bot.BuildTechnology(id, ogame.EspionageTechnologyID)
	bot.BuildTechnology(id, ogame.CombustionDriveID)
	if res.EnergyTechnology < 3 {
		bot.BuildTechnology(id, ogame.EnergyTechnologyID)
	}

	if res.ImpulseDrive < 3 {
		bot.BuildTechnology(id, ogame.ImpulseDriveID)
	} else {
		bot.BuildTechnology(id, ogame.ComputerTechnologyID)
		bot.BuildTechnology(id, ogame.ShieldingTechnologyID)
		bot.BuildTechnology(id, ogame.PlasmaTechnologyID)
	}

	mresearch := structs.Map(res)
	for k, v := range items.researchs {
		if v.(int64) != 0 {
			fmt.Println(k, ":", v)
		}
	}

	return mresearch
}

func launch() {

	for {
		items.planetes = bot.GetPlanets()
		items.fleets, _ = bot.GetFleets()
		items.planetinfos = nil
		for _, planete := range items.planetes {
			id := ogame.CelestialID(planete.ID)
			plinfo := gestionGlobal(id)
			items.planetinfos = append(items.planetinfos, plinfo)
			items.planeteName = planete.Name
			fmt.Println(planete.Name)
			//satProduction(planete.ID)
			items.researchs = setresearch(id)
		}
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("index.html")
	http.StripPrefix("/static/", http.FileServer(http.Dir("static")))
	/*
			c, _ := ioutil.ReadFile("index.html")
		    s := string(c)

		    t := template.New("")
		    t, _ = t.Parse(s)
	*/
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
