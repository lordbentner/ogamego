package main

import (
	"fmt"
	"math"
	"strings"

	"github.com/alaingilbert/ogame"
	"github.com/fatih/structs"
)

type PlaneteInfos struct {
	facilities   map[string]interface{}
	resources    map[string]interface{}
	res_build    map[string]interface{}
	ships        ogame.ShipsInfos
	consInBuild  ogame.ID
	countInBuild string
}

type Login struct {
	Universe string `form:"Universe"`
	User     string `form:"User"`
	Password string `form:"Password"`
}

//global informations of all in struct list
type GlobalList struct {
	planetes           []ogame.Planet
	researchs          map[string]interface{}
	fleets             []map[string]interface{}
	planetinfos        map[string]PlaneteInfos
	facilities         []map[string]interface{}
	resources          []map[string]interface{}
	res_build          []map[string]interface{}
	ships              []map[string]interface{}
	consInBuild        []ogame.ID
	countInBuild       []string
	researchInBuild    ogame.ID
	countResearchBuild string
}

func gestionUnderAttack(id ogame.CelestialID) {
	isAttack, _ := bot.IsUnderAttack()
	var i ogame.ID
	if isAttack {
		fmt.Println("ON EST ATTAQUES!!!!")
		for i = 408; i > 400; i-- {
			bot.BuildDefense(id, i, 10000)
		}
	}
}

func setExpedition(id ogame.CelestialID, coord ogame.Coordinate) {
	sh, _ := bot.GetShips(id)
	q := ogame.Quantifiable{ID: ogame.EspionageProbeID, Nbr: 10}
	r := ogame.Quantifiable{ID: ogame.LargeCargoID, Nbr: sh.LargeCargo}
	co := ogame.Coordinate{Galaxy: coord.Galaxy, System: coord.System, Position: 16}
	var quantList []ogame.Quantifiable
	quantList = append(quantList, q)
	quantList = append(quantList, r)
	bot.SendFleet(id, quantList, 100, co, ogame.Expedition, ogame.Resources{}, 0, 0)
}

func attackSpy(id ogame.CelestialID, coord ogame.Coordinate) {
	q := ogame.Quantifiable{ID: ogame.EspionageProbeID, Nbr: 1}
	var quantList []ogame.Quantifiable
	quantList = append(quantList, q)
	bot.SendFleet(id, quantList, 100, coord, ogame.Attack, ogame.Resources{}, 0, 0)
}

func gestionrapport(id ogame.CelestialID) {
	erm, _ := bot.GetEspionageReportMessages()

	if len(RapportEspionnage) > len(erm) {
		RapportEspionnage = make([]map[string]interface{}, len(erm))
	}

	for _, er := range erm {

		if er.Type == ogame.Action {
			fmt.Println(bot.GetEspionageReport(er.ID))
		}

		if er.Type == ogame.Report {
			msgR, _ := bot.GetEspionageReport(er.ID)
			re := structs.Map(msgR)
			totalres := msgR.Resources.Deuterium + msgR.Resources.Metal + msgR.Resources.Crystal
			if msgR.HasDefensesInformation == false || msgR.HasFleetInformation == false || totalres < 1000000 {
				bot.DeleteMessage(er.ID)
				return
			}

			di := structs.Map(msgR.ShipsInfos())
			for k, nbfl := range di {
				if nbfl.(int64) > 0 && !strings.Contains(k, "Solar") && !strings.Contains(k, "Probe") {
					fmt.Println("Vaisseaux detectes!!")
					bot.DeleteMessage(er.ID)
					return
				}
			}

			df := structs.Map(msgR.DefensesInfos())
			for k, nbdef := range df {
				if nbdef.(int64) > 0 && !strings.Contains(k, "Missiles") {
					fmt.Println("defense detectes!!")
					bot.DeleteMessage(er.ID)
					return
				}
			}

			hasAttacked := gestionAttack(id, totalres, msgR.Coordinate)
			fmt.Println("coordonnées:", msgR.Coordinate)
			RapportEspionnage = append(RapportEspionnage, re)
			if hasAttacked {
				bot.DeleteMessage(er.ID)
			}
		}
	}
}

func gestionEspionnage(id ogame.CelestialID, gal int64, sys int64) {
	galInfo, _ := bot.GalaxyInfos(gal, sys)
	var i int64
	var quantList []ogame.Quantifiable
	for i = 1; i <= 15; i++ {
		pos := galInfo.Position(i)
		if pos != nil {
			if pos.Inactive == true {
				q := ogame.Quantifiable{ID: ogame.EspionageProbeID, Nbr: 30}
				quantList = append(quantList, q)
				bot.SendFleet(id, quantList, 100, pos.Coordinate, ogame.Spy, ogame.Resources{}, 0, 0)
			}
		}
	}
}

func setShips(id ogame.CelestialID) map[string]interface{} {
	ships, _ := bot.GetShips(id)
	if ships.EspionageProbe < 30 {
		bot.BuildShips(id, ogame.EspionageProbeID, 1)
	}

	if ships.EspionageProbe < 100 {
		bot.BuildShips(id, ogame.LargeCargoID, 1)
	}

	sh := structs.Map(ships)
	return sh
}

func gestionAttack(id ogame.CelestialID, resource int64, where ogame.Coordinate) bool {
	ship, _ := bot.GetShips(id)
	if ship.SmallCargo == 0 || ship.LargeCargo == 0 {
		return false
	}

	fmt.Println("Lancement d'attaque sur ", where)
	var quantList []ogame.Quantifiable
	q := ogame.Quantifiable{ID: ogame.LargeCargoID, Nbr: resource / 25000}
	quantList = append(quantList, q)
	bot.SendFleet(id, quantList, 100, where, ogame.Attack, ogame.Resources{}, 0, 0)
	return true
}

func gestionGlobal(id ogame.CelestialID) PlaneteInfos {
	res, _ := bot.GetResourcesBuildings(id)
	resource, _ := bot.GetResources(id)
	fac, _ := bot.GetFacilities(id)
	if fac.RoboticsFactory < 10 {
		bot.BuildBuilding(id, ogame.RoboticsFactoryID)
	}

	if fac.Shipyard < 8 {
		bot.BuildBuilding(id, ogame.ShipyardID)
	}

	bot.BuildBuilding(id, ogame.NaniteFactoryID)
	bot.BuildBuilding(id, ogame.TerraformerID)
	if resource.Energy < 0 {
		bot.BuildBuilding(id, ogame.SolarPlantID)
	} else if res.MetalMine < res.CrystalMine+4 {
		bot.BuildBuilding(id, ogame.MetalStorageID)
		bot.BuildBuilding(id, ogame.MetalMineID)
	} else if res.CrystalMine < res.DeuteriumSynthesizer+4 {
		bot.BuildBuilding(id, ogame.CrystalStorageID)
		bot.BuildBuilding(id, ogame.CrystalMineID)
	} else {
		bot.BuildBuilding(id, ogame.DeuteriumTankID)
		bot.BuildBuilding(id, ogame.DeuteriumSynthesizerID)
	}

	coutTerraform := 100000 * math.Pow(2.0, float64(fac.Terraformer))
	if resource.Deuterium > int64(coutTerraform) && fac.ResearchLab < 1 {
		bot.BuildDefense(id, ogame.PlasmaTurretID, 1)
		fmt.Println("building plasma...")
	}

	if fac.SpaceDock < 7 {
		bot.BuildBuilding(id, ogame.SpaceDockID)
	}

	consInBuild, ctInBld, resinbuild, countres := bot.ConstructionsBeingBuilt(id)
	time := fmt.Sprintf("%dh %dmn %ds", ctInBld/3600, (ctInBld%3600)/60, ctInBld%60)
	var planetinfo PlaneteInfos
	planetinfo.res_build = structs.Map(res)
	planetinfo.resources = structs.Map(resource)
	planetinfo.facilities = structs.Map(fac)
	planetinfo.consInBuild = consInBuild
	planetinfo.countInBuild = time
	items.researchInBuild = resinbuild
	items.countResearchBuild = fmt.Sprintf("%dh %dmn %ds", countres/3600, (countres%3600)/60, countres%60)
	return planetinfo
}

func setresearch(id ogame.CelestialID) map[string]interface{} {
	bot.BuildTechnology(id, ogame.AstrophysicsID)
	bot.BuildTechnology(id, ogame.IntergalacticResearchNetworkID)
	res := bot.GetResearch()
	fac, _ := bot.GetFacilities(id)
	if res.EnergyTechnology < 8 {
		bot.BuildTechnology(id, ogame.EnergyTechnologyID)
	}

	if res.ImpulseDrive < 4 {
		bot.BuildTechnology(id, ogame.ImpulseDriveID)
	}

	if res.LaserTechnology < 10 {
		bot.BuildTechnology(id, ogame.LaserTechnologyID)
	}

	if fac.ResearchLab < 12 {
		bot.BuildBuilding(id, ogame.ResearchLabID)
	}

	if res.IonTechnology < 5 {
		bot.BuildTechnology(id, ogame.IonTechnologyID)
	}

	if res.HyperspaceTechnology < 8 {
		bot.BuildTechnology(id, ogame.HyperspaceTechnologyID)
	}

	bot.BuildTechnology(id, ogame.EspionageTechnologyID)
	bot.BuildTechnology(id, ogame.CombustionDriveID)
	bot.BuildTechnology(id, ogame.ArmourTechnologyID)
	bot.BuildTechnology(id, ogame.PlasmaTechnologyID)
	bot.BuildTechnology(id, ogame.PlasmaTechnologyID)
	bot.BuildTechnology(id, ogame.WeaponsTechnologyID)
	bot.BuildTechnology(id, ogame.ShieldingTechnology.ID)
	mresearch := structs.Map(res)
	mr := make(map[string]interface{})
	for k, v := range mresearch {
		if v.(int64) != 0 {
			mr[strings.Replace(k, "Technology", "", -1)] = v
		}
	}

	return mr
}

func transporter(id ogame.CelestialID, idwhere ogame.Coordinate) {
	ships, _ := bot.GetShips(id)
	if ships.LargeCargo < 5 {
		return
	}
	res, _ := bot.GetResources(id)
	fac, _ := bot.GetFacilities(id)
	if res.Deuterium > (100000*fac.Terraformer + 100000) {
		q := ogame.Quantifiable{ID: ogame.LargeCargoID, Nbr: 5}
		var quantList []ogame.Quantifiable
		quantList = append(quantList, q)
		co := ogame.Coordinate{Galaxy: idwhere.Galaxy, System: idwhere.System,
			Position: idwhere.Position}
		bot.SendFleet(id, quantList, 100, co, ogame.Transport,
			ogame.Resources{Deuterium: 100000}, 0, 0)
		fmt.Println("transport de déuterium vers:", idwhere)
	}
}

func add(a int, b int) int {
	return (a + b)
}

func sub(a int, b int) int {
	return (a - b)
}

func mul(a int, b int) int {
	return (a * b)
}

func div(a int, b int) int {
	return (a / b)
}
