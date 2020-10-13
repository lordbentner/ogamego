package main

import (
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/alaingilbert/ogame"
	"github.com/fatih/structs"
	"gopkg.in/macaron.v1"
)

var (
	bot *ogame.OGame
)
var (
	err error
)
var (
	startLog time.Time
)

var TimeDeconnecte []time.Time
var items GlobalList
var RapportEspionnage []map[string]interface{}

func getTimeInGame() string {
	if !bot.IsConnected() || !bot.IsLoggedIn() {
		startLog = time.Now()
		return "Non Connecté!"
	}

	currentTime := time.Now()
	return currentTime.Sub(startLog).String()
}

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
	var sys int64 = 64
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
			items.detailsRessources = make([]map[string]interface{}, len(items.planetes))
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
			items.detailsRessources[i] = plinfo.detailsRessources
			items.consInBuild[i] = plinfo.consInBuild
			items.countInBuild[i] = plinfo.countInBuild
			satProduction(planete.ID)
			inter := stres.IntergalacticResearchNetwork
			if i <= int(inter) {
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

			//setExpedition(id, planete.Coordinate)

			i++
		}
	}
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

	ctx.Data["time_con"] = getTimeInGame()
	ctx.Data["nbPages"] = nbPages
}
