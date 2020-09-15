package main

import (
	"fmt"
	"math"

	"github.com/alaingilbert/ogame"
	"github.com/fatih/structs"
)

type PlaneteInfos struct {
	facilities   map[string]interface{}
	resources    map[string]interface{}
	res_build    map[string]interface{}
	ships        ogame.ShipsInfos
	consInBuild  string
	countInBuild int64
}


//global informations of all in struct list
type GlobalList struct {
	planetes    []ogame.Planet
	researchs   map[string]interface{}
	fleets      []ogame.Fleet
	planetinfos map[string]PlaneteInfos
	facilities  []map[string]interface{}
	resources   []map[string]interface{}
	res_build   []map[string]interface{}
	ships       []ogame.ShipsInfos
	consInBuild  []string
	countInBuild []int64
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
	for _, er := range erm {
		if er.Type == ogame.Report {
			msgR, _ := bot.GetEspionageReport(er.ID)
			if msgR.HasFleetInformation == false && msgR.HasDefensesInformation == false && msgR.CounterEspionage == 0 {
				totalres := msgR.Resources.Metal + msgR.Resources.Crystal + msgR.Resources.Deuterium
				if totalres > 100000 {
					hasAttacked := gestionAttack(id, totalres, msgR.Coordinate)
					if hasAttacked {
						bot.DeleteMessage(er.ID)
					}
				}
			} else {
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
				fmt.Println("Coordinate:", pos.Coordinate)
				bot.SendFleet(id, quantList, 100, pos.Coordinate, ogame.Spy, ogame.Resources{}, 0, 0)
			}
		}
	}
}

func setShips(id ogame.CelestialID) ogame.ShipsInfos {
	ships, _ := bot.GetShips(id)
	/*if ships.EspionageProbe < 10 {
		bot.BuildShips(id, ogame.EspionageProbe.GetID(), 1)
	}

	if ships.LargeCargo < 15 {
		bot.BuildShips(id, ogame.LargeCargoID, 1)
	}*/

	return ships
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

	if fac.Shipyard < 4 {
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
		fmt.Println("build plasma...")
	}

	if fac.SpaceDock < 7 {
		bot.BuildBuilding(id, ogame.SpaceDockID)
	}

	consInBuild, countInBuild, _, _ := bot.ConstructionsBeingBuilt(id)
	var planetinfo PlaneteInfos
	planetinfo.res_build = structs.Map(res)
	planetinfo.resources = structs.Map(resource)
	planetinfo.facilities = structs.Map(fac)
	planetinfo.consInBuild = string(consInBuild)
	planetinfo.countInBuild = countInBuild
	for k, v := range planetinfo.res_build {
		if v.(int64) != 0 {
			fmt.Print(k, ":", v, ", ")
		}
	}

	fmt.Println(" ")

	for k, v := range planetinfo.resources {
		if v.(int64) != 0 {
			fmt.Print(k, ":", v, ", ")
		}
	}

	fmt.Println(" ")

	for k, v := range planetinfo.facilities {
		if v.(int64) != 0 {
			fmt.Print(k, ":", v, ", ")
		}
	}

	fmt.Println(" ")
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

	bot.BuildTechnology(id, ogame.EspionageTechnologyID)
	bot.BuildTechnology(id, ogame.CombustionDriveID)
	bot.BuildTechnology(id, ogame.ArmourTechnologyID)
	bot.BuildTechnology(id, ogame.PlasmaTechnologyID)

	mresearch := structs.Map(res)
	for k, v := range items.researchs {
		if v.(int64) != 0 {
			fmt.Print(k, ":", v, ",")
		}
	}

	fmt.Println(" ")

	return mresearch
}

