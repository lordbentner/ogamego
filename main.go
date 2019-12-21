package main

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/alaingilbert/ogame"
	"github.com/fatih/structs"
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
var bot, err = ogame.New("Pasiphae", "nemesism@hotmail.fr", "pencilcho44", "fr")

func satProduction(id ogame.PlanetID) {
	pl, _ := bot.GetPlanet(id)
	//fac, _ := bot.GetResourcesBuildings(ogame.CelestialID(id))
	temp := pl.Temperature
	satprod := ogame.SolarSatellite.Production(temp, 1)
	//cenprice := 20 * math.Pow(1.1, float64(fac.SolarPlant))
	fmt.Print("sattelitte production:")
	fmt.Println(satprod)
}

func setresearch(id ogame.CelestialID) {
	bot.BuildTechnology(id, ogame.AstrophysicsID)
	bot.BuildTechnology(id, ogame.IntergalacticResearchNetworkID)
	bot.BuildTechnology(id, ogame.ComputerTechnologyID)
	bot.BuildTechnology(id, ogame.ShieldingTechnologyID)
	bot.BuildTechnology(id, ogame.PlasmaTechnologyID)
	res := bot.GetResearch()
	bot.BuildTechnology(id, ogame.EspionageTechnologyID)
	bot.BuildTechnology(id, ogame.CombustionDriveID)
	if res.EnergyTechnology < 12 {
		bot.BuildTechnology(id, ogame.EnergyTechnologyID)
	}

	mresearch := structs.Map(res)
	for k, v := range mresearch {
		if v.(int64) != 0 {
			fmt.Println(k, v)
		}
	}

}

func gestionGlobal(id ogame.CelestialID) {

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
		fmt.Println("build plama...")
	}

	if fac.ResearchLab < 12 {
		bot.BuildBuilding(id, ogame.ResearchLabID)
	}

	mres := structs.Map(res)
	mresource := structs.Map(resource)
	mfac := structs.Map(fac)
	for k, v := range mres {
		if v.(int64) != 0 {
			fmt.Println(k, v)
		}
	}

	for k, v := range mresource {
		if v.(int64) != 0 {
			fmt.Println(k, v)
		}
	}

	for k, v := range mfac {
		if v.(int64) != 0 {
			fmt.Println(k, ":", v)
		}
	}
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
			setresearch(id)
		}
	}

}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
	t, _ := template.ParseFiles("index.html")
	p := "test"

	t.Execute(w, p)
}

func main() {

	// Instantiate loader for kubeconfig file.
	/*kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)

	// Determine the Namespace referenced by the current context in the
	// kubeconfig file.
	namespace, _, err := kubeconfig.Namespace()
	if err != nil {
		panic(err)
	}

	// Get a rest.Config from the kubeconfig file.  This will be passed into all
	// the client objects we create.
	restconfig, err := kubeconfig.ClientConfig()
	if err != nil {
		panic(err)
	}

	// Create a Kubernetes core/v1 client.
	coreclient, err := corev1client.NewForConfig(restconfig)
	if err != nil {
		panic(err)
	}

	// Create an OpenShift build/v1 client.
	buildclient, err := buildv1client.NewForConfig(restconfig)
	if err != nil {
		panic(err)
	}

	// List all Pods in our current Namespace.
	pods, err := coreclient.Pods(namespace).List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Pods in namespace %s:\n", namespace)
	for _, pod := range pods.Items {
		fmt.Printf("  %s\n", pod.Name)
	}

	// List all Builds in our current Namespace.
	/*builds, err := buildclient.Builds(namespace).List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}*/

	/*fmt.Printf("Builds in namespace %s:\n", namespace)
	for _, build := range builds.Items {
		fmt.Printf("  %s\n", build.Name)
	}*/

	go launch()
	//http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", handler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
