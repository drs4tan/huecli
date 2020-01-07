package main

//VtGw9pfjWX1V6AYgpWwY2M4I0iyiRp82DXKLOWva
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
	flag.BoolVar(&optList, "list", false, "List all Hue lights")
	flag.StringVar(&optFind, "find", "", "Find Hue light")
	flag.BoolVar(&optAlert, "alert", false, "Blink light")
	flag.Parse()
}

func main() {

	println(optAlert)
	bridge := huego.New("192.168.1.101", "VtGw9pfjWX1V6AYgpWwY2M4I0iyiRp82DXKLOWva")

	if optAlert == true {
		testl := findlight(optFind, bridge)
		testl.Alert("select")
	}

	if optList == true {
		listlights(bridge)
	}

	//findlight("he", bridge)
	//findall(bridge)

}

/*
func findlights(nameOfLight string, bridge *huego.Bridge) (light huego.Light) {
	allTheLights, err := bridge.GetLights()
	if err != nil {
		println(err.Error)
	}


	println("ID of lights:")
	for i := range allTheLights {

		if strings.Contains(allTheLights[i].Name, nameOfLight) {

			fmt.Println(allTheLights[i].ID)
			return allTheLights[i]
		}
	}
	return
}
*/

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
