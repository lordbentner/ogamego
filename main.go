package main

import (
	"fmt"

	"github.com/alaingilbert/ogame"
)

type PlaneteInfos struct {
	facilities   ogame.Facilities
	resources    ogame.Resources
	res_build    ogame.ResourcesBuildings
	ships        ogame.ShipsInfos
	consInBuild  string
	countInBuild int
}

type GlobalList struct {
	planetes    []ogame.Planet
	researchs   ogame.Researches
	fleets      []ogame.Fleet
	planetinfos []PlaneteInfos
}

var isInit bool = false
var bot, err = ogame.New("Aquarius", "nemesism@hotmail.fr", "pencilcho44", "fr")

func satProduction(id ogame.PlanetID) {
	pl, _ := bot.GetPlanet(id)
	//fac, _ := bot.GetResourcesBuildings(ogame.CelestialID(id))
	temp := pl.Temperature
	satprod := ogame.SolarSatellite.Production(temp, 1)
	//cenprice := 20 * math.Pow(1.1, float64(fac.SolarPlant))
	fmt.Print("sattelitte production:")
	fmt.Println(satprod)

}

func gestionGlobal(id ogame.CelestialID) {
	res, _ := bot.GetResourcesBuildings(id)
	resource, _ := bot.GetResources(id)
	fac, _ := bot.GetFacilities(id)
	bot.BuildBuilding(id, ogame.NaniteFactoryID)
	bot.BuildBuilding(id, ogame.TerraformerID)
	if resource.Energy < 0 {
		bot.BuildBuilding(id, ogame.SolarPlantID)
	} else if res.MetalMine < res.CrystalMine+4 {
		bot.BuildBuilding(id, ogame.MetalMineID)
	} else if res.CrystalMine < res.DeuteriumSynthesizer+4 {
		bot.BuildBuilding(id, ogame.CrystalMineID)
	} else {
		bot.BuildBuilding(id, ogame.DeuteriumSynthesizerID)
	}

	if fac.ResearchLab < 12 {
		bot.BuildBuilding(id, ogame.ResearchLabID)
	}

	fmt.Println(res)
	fmt.Println(resource)
}

func launch() GlobalList {

	flee, _ := bot.GetFleets()
	planets := bot.GetPlanets()

	list := GlobalList{
		planetes:  bot.GetPlanets(),
		researchs: bot.GetResearch(),
		fleets:    flee,
		//planetinfos: planeteinfos,
	}

	var plinfo PlaneteInfos
	for {
		for _, planete := range planets {
			id := ogame.CelestialID(planete.ID)
			plid := planete.ID
			plinfo.facilities, err = bot.GetFacilities(id)
			plinfo.resources, err = bot.GetResources(id)
			plinfo.res_build, err = bot.GetResourcesBuildings(id)
			list.planetinfos = append(list.planetinfos, plinfo)
			fmt.Println(planete.Name)
			gestionGlobal(id)
			satProduction(plid)
		}
	}

}

func main() {

	launch()
	/*http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", TestHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}*/
}

/*s := reflect.ValueOf(&list.planetinfos[0].facilities).Elem()
typeOfs := s.Type()
for i := 0; i < s.NumField(); i++ {
	f := s.Field(i).Interface()
	fmt.Println("Niveau " + typeOfs.Field(i).Name + ":" + f)
}*/
