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
}

//global informations of all in strcut list
type GlobalList struct {
	planetes    []ogame.Planet
	researchs   map[string]interface{}
	fleets      []ogame.Fleet
	planetinfos []PlaneteInfos
	planeteName string
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

	if fac.ResearchLab < 12 {
		bot.BuildBuilding(id, ogame.ResearchLabID)
	}

	var planetinfo PlaneteInfos
	planetinfo.res_build = structs.Map(res)
	planetinfo.resources = structs.Map(resource)
	planetinfo.facilities = structs.Map(fac)
	for k, v := range planetinfo.res_build {
		if v.(int64) != 0 {
			fmt.Println(k, v)
		}
	}

	for k, v := range planetinfo.resources {
		if v.(int64) != 0 {
			fmt.Println(k, v)
		}
	}

	for k, v := range planetinfo.facilities {
		if v.(int64) != 0 {
			fmt.Println(k, ":", v)
		}
	}

	return planetinfo
}
