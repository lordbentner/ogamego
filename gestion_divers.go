package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
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
var Logout = false
var BuildLune []map[string]interface{}
var vlistAttack []map[string]interface{}

func getJSONlogin() Login {
	jsonfile, err := os.Open("data.json")
	bytevalue, _ := ioutil.ReadAll(jsonfile)
	var mlogin Login
	json.Unmarshal(bytevalue, &mlogin)
	if err != nil {
		panic(err)
	}

	return mlogin
}

func getTimeInGame() (string, string) {
	login := getJSONlogin()
	if bot == nil {
		return "Non Connecté!", ""
	}

	currentTime := time.Now()
	return currentTime.Sub(startLog).String(), login.User
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
	var gal int64 = 4
	var sys int64 = 2
	for {
		if Logout {
			fmt.Println("Déconnecté!!")
			break
		}

		login := getJSONlogin()
		if bot == nil {
			startLog = time.Now()
			bot, err = ogame.New(login.Universe, login.User, login.Password, "fr")
			if err != nil {
				panic(err)
			}
		}

		items.points = bot.GetUserInfos().Points
		items.planetes = bot.GetPlanets()
		items.lunes = bot.GetMoons()
		BuildLune = nil
		for _, lune := range items.lunes {
			id := ogame.CelestialID(lune.ID)
			botfac, _ := bot.GetFacilities(id)
			BuildLune = append(BuildLune, structs.Map(botfac))
		}
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
			items.productions = make([]map[string]interface{}, len(items.planetes))
		}

		if len(items.fleets) < len(fl) {
			items.fleets = make([]map[string]interface{}, len(fl))
		}

		for j, fle := range fl {
			items.fleets[j] = structs.Map(fle)
		}

		for _, lune := range items.lunes {
			id := ogame.CelestialID(lune.ID)
			bot.BuildBuilding(id, ogame.LunarBaseID)
			bot.BuildBuilding(id, ogame.SensorPhalanxID)
		}

		for _, planete := range items.planetes {
			id := ogame.CelestialID(planete.ID)
			vlistAttack = gestionUnderAttack(id)
			plinfo := gestionGlobal(id)
			items.facilities[i] = plinfo.facilities
			items.resources[i] = plinfo.resources
			items.res_build[i] = plinfo.res_build
			items.detailsRessources[i] = plinfo.detailsRessources
			items.consInBuild[i] = plinfo.consInBuild
			items.countInBuild[i] = plinfo.countInBuild
			items.productions[i] = plinfo.productions
			satProduction(planete.ID)
			items.ships[i] = setShips(id)
			inter := stres.IntergalacticResearchNetwork
			if i <= int(inter) {
				items.researchs = setresearch(id)
			} /*else {
				transporter(id, items.planetes[i].Coordinate)
			}*/

			if sys >= 500 {
				sys = 1
				gal++
			}
			if gal >= 5 {
				gal = 1
			}

			fmt.Println("Vaisseaux:", items.ships[i])
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

	time, user := getTimeInGame()
	ctx.Data["time_con"] = time
	ctx.Data["user"] = user
	ctx.Data["nbPages"] = nbPages
}
