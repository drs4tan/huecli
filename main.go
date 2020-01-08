package main

//Hue Username:
//VtGw9pfjWX1V6AYgpWwY2M4I0iyiRp82DXKLOWva

/*
TODO:
	!finish findlights
	!alerts work with both findlight/s
	COLORS!



*/



type RGBColor struct {
	R uint16
	G uint16
	B uint16
}

type XYColor struct {
	X uint8
	Y uint8
}

func (XY *XYColor) ConvToRGB() {
	var RGB RGBColor

	


}

import (
	"flag"
	"fmt"
	"strings"

	"github.com/amimof/huego"
)

//Vars for defualt flags
//Built version will need changes!

var optList bool = false
var optFind string = "off"
var optAlert bool = false

const 



func init() {
	flag.BoolVar(&optList, "list", optList, "List all Hue lights with ID and name")
	flag.StringVar(&optFind, "f", optFind, "Find Hue lights with the name value")
	flag.BoolVar(&optAlert, "alert", optAlert, "Blink lights")

	flag.Parse()
}

func changelightcolor(bridge *huego.Bridge) {
	test1 := findlights(optFind, bridge)
	for i := range test1 {
		err := test1[i].Hue(30000)
		fmt.Println(err)
	}
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

	changelightcolor(bridge)

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

//This func is depracated
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

