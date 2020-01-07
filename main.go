package main

//Hue Username:
//VtGw9pfjWX1V6AYgpWwY2M4I0iyiRp82DXKLOWva

/*
TODO:
	finish findlights
	alerts work with both findlight/s
	COLORS!



*/

import (
	"flag"
	"fmt"
	"strings"

	"github.com/amimof/huego"
)

//Vars for defualt flags
//Built version will need changes!

var optList bool = false
var optFind string = ""
var optAlert bool = false

func init() {
	flag.BoolVar(&optList, "list", optList, "List all Hue lights")
	flag.StringVar(&optFind, "find", optFind, "Find Hue light")
	flag.BoolVar(&optAlert, "alert", optAlert, "Blink light")
	flag.Parse()
}

func main() {

	bridge := huego.New("192.168.1.101", "VtGw9pfjWX1V6AYgpWwY2M4I0iyiRp82DXKLOWva")

	if optAlert == true {
		test1 := findlights(optFind, bridge)
		for i := range test1 {
			test1[i].Alert("select")
		}
	}

	if optList == true {
		listlights(bridge)
	}

	//findlight("he", bridge)
	//findall(bridge)

}

func findlights(nameOfLight string, bridge *huego.Bridge) (light []huego.Light) {
	allTheLights, err := bridge.GetLights()
	if err != nil {
		println(err.Error)
	}

	var matchedLights []huego.Light

	println("ID of lights:")
	for i := range allTheLights {

		if strings.Contains(allTheLights[i].Name, nameOfLight) {

			fmt.Println(allTheLights[i].ID)
			matchedLights = append(matchedLights, allTheLights[i])
		}
	}
	return matchedLights
}

func findlight(nameOfLight string, bridge *huego.Bridge) (light huego.Light) {
	allTheLights, err := bridge.GetLights()
	if err != nil {
		println(err.Error)
	}

	for i := range allTheLights {

		if strings.Contains(allTheLights[i].Name, nameOfLight) {
			println("ID of light:")
			fmt.Println(allTheLights[i].ID)
			return allTheLights[i]
		}
	}
	return
}

func listlights(bridge *huego.Bridge) {
	allTheLights, err := bridge.GetLights()
	if err != nil {
		println(err.Error)
	}
	println("Listing all lights...")
	for i := range allTheLights {
		fmt.Println(allTheLights[i].ID, allTheLights[i].Name)
	}

}
