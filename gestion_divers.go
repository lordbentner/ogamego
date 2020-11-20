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

type coordSpy struct {
	Galaxy int64
	System int64
}

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

func getJSONDataboard() coordSpy {
	jsonfile, err := os.Open("databoard.json")
	bytevalue, _ := ioutil.ReadAll(jsonfile)
	var co coordSpy
	json.Unmarshal(bytevalue, &co)
	if err != nil {
		panic(err)
	}

	return co
}

func getTimeInGame() (string, string) {
	login := getJSONlogin()
	if bot == nil {
		return "Non Connecté!", ""
	} else if !bot.IsLoggedIn() {
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
	lastspy := getJSONDataboard()
	for {
		if Logout {
			fmt.Println("Déconnecté!!")
			break
		}

		if bot == nil {
			fmt.Println("bot vide!!!!")
			login := getJSONlogin()
			startLog = time.Now()
			bot, err = ogame.New(login.Universe, login.User, login.Password, "fr")
			if err != nil {
				panic(err)
			}

			for {
				if bot != nil {
					break
				}
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
		i := 0
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

		if len(items.planetes) > len(items.planetinfos) {
			items.planetinfos = make([]PlaneteInfos, len(items.planetes))
		}

		for _, planete := range items.planetes {
			id := ogame.CelestialID(planete.ID)
			inter := bot.GetResearch().IntergalacticResearchNetwork
			if i <= int(inter) {
				items.researchs = setresearch(id)
			}

			gestionUnderAttack(id)
			plinfo := gestionGlobal(id)
			plinfo.Planetes = planete
			items.planetinfos[i] = plinfo
			satProduction(planete.ID)

			if lastspy.System >= 500 {
				lastspy.System = 1
				lastspy.Galaxy++
			}
			if lastspy.Galaxy >= 5 {
				lastspy.Galaxy = 1
			}

			fmt.Println("Vaisseaux:", plinfo.Ships)
			if len(fl) <= int(bot.GetResearch().ComputerTechnology) {
				gestionEspionnage(id, lastspy.Galaxy, lastspy.System)
				gestionrapport(id)
				lastspy.System++
			}

			setExpedition(id, planete.Coordinate)
			i++
		}

		file, _ := json.Marshal(lastspy)
		_ = ioutil.WriteFile("databoard.json", file, 0777)
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
	ctx.Data["point"] = items.points
}
