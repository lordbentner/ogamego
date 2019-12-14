package main

import "fmt"
import "os"
import "github.com/alaingilbert/ogame"

func main() {
	universe := "Aquarius" // eg: Bellatrix
	username := "nemesism@hotmail.fr" // eg: email@gmail.com
	password := "pencilcho44" // eg: *****
	language := "fr" // eg: en
	bot, err := ogame.New(universe, username, password, language)
	if err != nil {
		panic(err)
	}
	attacked, _ := bot.IsUnderAttack()
	fmt.Println("izi")
	fmt.Println(attacked) // False
}