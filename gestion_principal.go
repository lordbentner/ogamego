package main

import (
	"fmt"
	"strings"
	"sync"

	"github.com/alaingilbert/ogame"
	"github.com/fatih/structs"
	"github.com/stretchr/stew/slice"
)

type PlaneteInfos struct {
	Planetes          ogame.Planet
	Facilities        map[string]interface{}
	Resources         map[string]interface{}
	Res_build         map[string]interface{}
	DetailsRessources map[string]interface{}
	Ships             map[string]interface{}
	ConsInBuild       ogame.ID
	CountInBuild      string
	CountProductions  string
	Productions       map[ogame.ID]int64
	Defenses          map[string]interface{}
}

type Login struct {
	Universe string `form:"Universe"`
	User     string `form:"User"`
	Password string `form:"Password"`
}

//global informations of all in struct list
type GlobalList struct {
	planetes           []ogame.Planet
	lunes              []ogame.Moon
	researchs          map[string]interface{}
	fleets             []map[string]interface{}
	planetinfos        []PlaneteInfos
	researchInBuild    ogame.ID
	countResearchBuild string
	countProductions   []string
	points             int64
	lastEspionnage     [2]int64
}

var mu sync.Mutex
var planetinfo PlaneteInfos

func gestionUnderAttack(id ogame.CelestialID) {
	listattack, _ := bot.GetAttacks()
	var i ogame.ID
	for _, attack := range listattack {
		vlistAttack = append(vlistAttack, structs.Map(attack))
		for _, planet := range items.planetes {
			if attack.Destination.System == planet.Coordinate.System && attack.MissionType == ogame.Attack {
				for i = 408; i > 400; i-- {
					bot.BuildDefense(planet.GetID(), i, 10000)
				}
			}
		}
	}
}

func setExpedition(id ogame.CelestialID, coord ogame.Coordinate) {
	sh, _ := bot.GetShips(id)
	if sh.EspionageProbe > 50 {
		return
	}

	q := ogame.Quantifiable{ID: ogame.EspionageProbeID, Nbr: 10}
	r := ogame.Quantifiable{ID: ogame.LargeCargoID, Nbr: sh.LargeCargo}
	s := ogame.Quantifiable{ID: ogame.SmallCargoID, Nbr: sh.SmallCargo}
	co := ogame.Coordinate{Galaxy: coord.Galaxy, System: coord.System, Position: 16}
	var quantList []ogame.Quantifiable
	quantList = append(quantList, q, r, s)
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
	var Rapport []map[string]interface{}
	for _, er := range erm {
		if er.Type == ogame.Report {
			msgR, _ := bot.GetEspionageReport(er.ID)
			totalres := msgR.Resources.Deuterium + msgR.Resources.Metal + msgR.Resources.Crystal
			if msgR.HasDefensesInformation == false || msgR.HasFleetInformation == false || totalres < 2000000 {
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
			//RapportEspionnage = append(RapportEspionnage, re)
			if hasAttacked {
				bot.DeleteMessage(er.ID)
			}
		} else {
			msgR, _ := bot.GetEspionageReport(er.ID)
			fmt.Println("Rapport autre:", msgR)
			Rapport = append(Rapport, structs.Map(msgR))
		}
	}

	RapportEspionnage = Rapport
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
	prlc := planetinfo.Productions[ogame.LargeCargoID]
	prspy := planetinfo.Productions[ogame.EspionageProbeID]
	if ships.EspionageProbe+prspy < 100 {
		bot.BuildShips(id, ogame.EspionageProbeID, 100-(ships.LargeCargo+prlc))
	}

	if ships.LargeCargo+prlc < 100 {
		bot.BuildShips(id, ogame.LargeCargoID, 100-(ships.LargeCargo+prlc))
	}

	ships, _ = bot.GetShips(id)
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
	sh := resource / 25000
	if sh > ship.LargeCargo/25000 {
		sh = resource / 5000
	}
	q := ogame.Quantifiable{ID: ogame.LargeCargoID, Nbr: sh}
	quantList = append(quantList, q)
	bot.SendFleet(id, quantList, 100, where, ogame.Attack, ogame.Resources{}, 0, 0)
	return true
}

func gestionGlobal(id ogame.CelestialID) PlaneteInfos {
	mu.Lock()
	defer mu.Unlock()
	res, _ := bot.GetResourcesBuildings(id)
	resource, _ := bot.GetResources(id)
	fac, _ := bot.GetFacilities(id)
	detres, _ := bot.GetResourcesDetails(id)
	bot.BuildBuilding(id, ogame.NaniteFactoryID)
	if fac.RoboticsFactory < 10 {
		bot.BuildBuilding(id, ogame.RoboticsFactoryID)
	}

	if fac.Shipyard < 12 {
		bot.BuildBuilding(id, ogame.ShipyardID)
	}

	if fac.MissileSilo < 8 {
		bot.BuildBuilding(id, ogame.MissileSiloID)
	}

	bot.BuildDefense(id, ogame.AntiBallisticMissilesID, 10)
	bot.BuildDefense(id, ogame.SmallShieldDomeID, 1)
	bot.BuildDefense(id, ogame.LargeShieldDomeID, 1)
	bot.BuildBuilding(id, ogame.TerraformerID)
	if resource.Energy < 12 {
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

	if fac.SpaceDock < 7 {
		bot.BuildBuilding(id, ogame.SpaceDockID)
	}

	diffmetal := detres.Metal.Available - detres.Metal.StorageCapacity
	diffcrystal := detres.Crystal.Available - detres.Crystal.StorageCapacity
	diffdeut := detres.Deuterium.Available - detres.Deuterium.StorageCapacity
	resdiff := ogame.Resources{Metal: 0, Crystal: 0, Deuterium: 0}
	if diffmetal >= 0 {
		resdiff.Metal = 500000 + diffmetal
	}
	if diffcrystal >= 0 {
		resdiff.Crystal = 500000 + diffcrystal
	}
	if diffdeut >= 0 {
		resdiff.Deuterium = 500000 + diffdeut
	}

	fmt.Println("ressource:", resdiff)
	if len(bot.GetMoons()) > 0 {
		transporter(id, bot.GetMoons()[0].Coordinate, resdiff)
	}
	consInBuild, ctInBld, resinbuild, countres := bot.ConstructionsBeingBuilt(id)
	time := fmt.Sprintf("%dh %dmn %ds", ctInBld/3600, (ctInBld%3600)/60, ctInBld%60)
	prod, nb, _ := bot.GetProduction(id)
	def, _ := bot.GetDefense(id)
	listprod := make(map[ogame.ID]int64)
	var listTypeprod []ogame.ID
	for _, pr := range prod {
		if !slice.Contains(listTypeprod, pr.ID) {
			listTypeprod = append(listTypeprod, pr.ID)
		}
	}
	for _, ltp := range listTypeprod {
		listprod[ltp] = 0
		for _, pr := range prod {
			if ltp == pr.ID {
				listprod[ltp] += pr.Nbr
			}
		}
	}

	planetinfo.Res_build = structs.Map(res)
	planetinfo.Resources = structs.Map(resource)
	planetinfo.Facilities = structs.Map(fac)
	planetinfo.DetailsRessources = structs.Map(detres)
	planetinfo.Ships = setShips(id)
	planetinfo.Defenses = structs.Map(def)
	planetinfo.ConsInBuild = consInBuild
	planetinfo.CountInBuild = time
	planetinfo.Productions = listprod
	planetinfo.CountProductions = secondsToHuman(int(nb))
	items.researchInBuild = resinbuild
	items.countResearchBuild = secondsToHuman(int(countres))
	return planetinfo
}

func setresearch(id ogame.CelestialID) map[string]interface{} {
	res := bot.GetResearch()
	fac, _ := bot.GetFacilities(id)
	mresearch := structs.Map(res)
	mr := make(map[string]interface{})
	for k, v := range mresearch {
		if v.(int64) != 0 {
			mr[strings.Replace(k, "Technology", "", -1)] = v
		}
	}

	if fac.ResearchLab < 12 {
		bot.BuildBuilding(id, ogame.ResearchLabID)
	}

	bot.BuildTechnology(id, ogame.AstrophysicsID)
	bot.BuildTechnology(id, ogame.ComputerTechnologyID)
	bot.BuildTechnology(id, ogame.IntergalacticResearchNetworkID)
	fmt.Println("Recherche...")
	if res.EnergyTechnology < 12 {
		bot.BuildTechnology(id, ogame.EnergyTechnologyID)
	}

	if res.ImpulseDrive < 4 {
		bot.BuildTechnology(id, ogame.ImpulseDriveID)
	}

	if res.LaserTechnology < 10 {
		bot.BuildTechnology(id, ogame.LaserTechnologyID)
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
	return mr
}

func transporter(id ogame.CelestialID, idwhere ogame.Coordinate, resource ogame.Resources) {
	ships, _ := bot.GetShips(id)
	if ships.LargeCargo < 60 {
		return
	}

	if resource.Metal == 0 && resource.Crystal == 0 && resource.Deuterium == 0 {
		return
	}

	q := ogame.Quantifiable{ID: ogame.LargeCargoID, Nbr: 60}
	var quantList []ogame.Quantifiable
	quantList = append(quantList, q)
	co := ogame.Coordinate{Galaxy: idwhere.Galaxy, System: idwhere.System,
		Position: idwhere.Position, Type: ogame.MoonType}
	bot.SendFleet(id, quantList, 100, co, ogame.Transport,
		resource, 0, 0)
	fmt.Println("transport de ressources vers:", idwhere, "avec:", resource)
}
