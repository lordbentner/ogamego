package main

//https://ogamebot.uc.r.appspot.com
import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alaingilbert/ogame"
	"github.com/go-macaron/binding"
	"github.com/kardianos/service"
	"gopkg.in/macaron.v1"
)

var logger service.Logger
var m *macaron.Macaron

type program struct{}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}
func (p *program) run() {
	http.Handle("/", m)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	return nil
}

func main() {
	m = macaron.Classic()
	m.Use(macaron.Renderer())
	m.Get("/", func(ctx *macaron.Context) {
		ctx.HTML(200, "index")
	})
	m.Post("/ogame", binding.Form(Login{}), func(login Login, ctx *macaron.Context, r *http.Request) {
		if login.Universe == "" || login.User == "" || login.Password == "" {
			ctx.Redirect("/")
		} else {
			bot, err = ogame.New(login.Universe, login.User, login.Password, "fr")
			file, _ := json.Marshal(login)
			_ = ioutil.WriteFile("data.json", file, 0777)
			if err != nil {
				panic(err)
			} else {
				startLog = time.Now()
				Logout = false
				go launch()
				/*c := appengine.NewContext(r)
				rerr := runtime.RunInBackground(c, func(c appengine.Context) {
					launch()
				})*/
				ctx.Redirect("/databoard")
			}
		}
	})
	m.Get("/databoard", func(ctx *macaron.Context) {
		/*v := reflect.ValueOf(items.planetinfos)
		for i := 0; i < v.NumField(); i++ {
				if v.Type().Field(i).Type.String() == "[]map[string]interface {}" {
			}
			fmt.Println("key:", v.Type().Field(i).Name, "type:", v.Type().Field(i).Type.String())
		}*/
		/*	ctx.Data["lunes"] = items.lunes
			ctx.Data["planetes"] = items.planetes
			ctx.Data["researchs"] = items.researchs
			ctx.Data["facilities"] = items.facilities
			ctx.Data["resources"] = items.resources
			ctx.Data["resdetails"] = items.detailsRessources
			ctx.Data["res_build"] = items.res_build
			ctx.Data["ships"] = items.ships
			ctx.Data["consInBuild"] = items.consInBuild
			ctx.Data["countInBuild"] = items.countInBuild
			ctx.Data["resInBuild"] = items.researchInBuild
			ctx.Data["countResInBuild"] = items.countResearchBuild
			ctx.Data["productions"] = items.productions
			ctx.Data["countProductions"] = items.countProductions*/
		ctx.Data["resInBuild"] = items.researchInBuild
		ctx.Data["point"] = items.points
		ctx.Data["planetinfos"] = items.planetinfos
		ctx.Data["researchs"] = items.researchs
		ctx.Data["countResInBuild"] = items.countResearchBuild
		time, user := getTimeInGame()
		ctx.Data["time_con"] = time
		ctx.Data["user"] = user
		ctx.Data["BuildLune"] = BuildLune
		ctx.HTML(200, "ogame")
	})

	m.Get("/flottes", func(ctx *macaron.Context, req *http.Request) {
		buildPage(ctx, req, 4, len(items.fleets))
		ctx.Data["flottes"] = items.fleets
		ctx.HTML(200, "flottes")
	})

	m.Get("/rapports", func(ctx *macaron.Context, req *http.Request) {
		buildPage(ctx, req, 15, len(RapportEspionnage))
		ctx.Data["spy"] = vlistAttack
		ctx.HTML(200, "rapports")
	})

	m.Get("/quit", func(ctx *macaron.Context) {
		if bot != nil {
			bot.Logout()
		}
		Logout = true
		ctx.Redirect("/")
	})

	http.Handle("/", m)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}

	/*svcConfig := &service.Config{
		Name:        "OgameBot",
		DisplayName: "ogame bot",
		Description: "This is a test Go service.",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	logger, err = s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}
	err = s.Run()
	if err != nil {
		logger.Error(err)
	}*/
}
