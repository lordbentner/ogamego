package main

import (
	"fmt"

	"github.com/alaingilbert/ogame"
	"github.com/fatih/structs"
)

type PlaneteInfos struct {
	facilities   map[string]interface{}
	resources    map[string]interface{}
	res_build    map[string]interface{}
	ships        map[string]interface{}
	consInBuild  string
	countInBuild int
	coord        ogame.Coordinate
}

//global informations of all in strcut list
type GlobalList struct {
	planetes    []ogame.Planet
	researchs   map[string]interface{}
	fleets      []ogame.Fleet
	planetinfos []PlaneteInfos
	planeteName string
}

func gestionrapport() {
	erm, _ := bot.GetEspionageReportMessages()
	fmt.Print("Rapport d'espionnage:")
	for _, er := range erm {
		fmt.Println(bot.GetEspionageReport(er.ID))
	}
}

func gestionGlobal(id ogame.CelestialID) PlaneteInfos {

	att, _ := bot.IsUnderAttack()
	if att {
		fmt.Println("Nous sommes attaqués!!!!")
	}

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

	if resource.Deuterium > 830.000 && fac.ResearchLab < 1 {
		bot.BuildDefense(id, ogame.PlasmaTurretID, 1)
		fmt.Println("build plasma...")
	}

	var planetinfo PlaneteInfos
	planetinfo.res_build = structs.Map(res)
	planetinfo.resources = structs.Map(resource)
	planetinfo.facilities = structs.Map(fac)
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
	//bot.BuildTechnology(id, ogame.EspionageTechnologyID)
	bot.BuildTechnology(id, ogame.CombustionDriveID)
	bot.BuildTechnology(id, ogame.ArmourTechnologyID)
	if res.EnergyTechnology < 3 {
		bot.BuildTechnology(id, ogame.EnergyTechnologyID)
	}

	/*if res.ImpulseDrive < 3 {
		bot.BuildTechnology(id, ogame.ImpulseDriveID)
	} else {
		bot.BuildTechnology(id, ogame.ComputerTechnologyID)
		bot.BuildTechnology(id, ogame.ShieldingTechnologyID)
		bot.BuildTechnology(id, ogame.PlasmaTechnologyID)
	}*/

	if fac.ResearchLab < 6 {
		bot.BuildBuilding(id, ogame.ResearchLabID)
	}

	mresearch := structs.Map(res)
	for k, v := range items.researchs {
		if v.(int64) != 0 {
			fmt.Print(k, ":", v, ",")
		}
	}

	fmt.Println(" ")

	return mresearch
}
