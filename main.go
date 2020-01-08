package main

//Hue Username:
//VtGw9pfjWX1V6AYgpWwY2M4I0iyiRp82DXKLOWva

/*
TODO:
	!finish findlights
	!alerts work with both findlight/s
	COLORS!



*/

import (
	"flag"
	"fmt"
	"strings"

	"github.com/amimof/huego"
)

type RGBColor struct {
	R float32
	G float32
	B float32
}

type XYColor struct {
	X float32
	Y float32
}

func (xy XYColor) ConvToArray() []float32 {
	var xyarray []float32
	xyarray = append(xyarray, xy.X, xy.Y)
	return xyarray
}

func (RGB *RGBColor) ConvToXY() XYColor {
	var xy XYColor
	var X, Y, Z, cx, cy float32

	X = 0.4124*RGB.R + 0.3576*RGB.G + 0.1805*RGB.B
	Y = 0.2126*RGB.R + 0.7152*RGB.G + 0.0722*RGB.B
	Z = 0.0192*RGB.R + 0.1192*RGB.G + 0.9505*RGB.B
	cx = X / (X + Y + Z)
	cy = Y / (X + Y + Z)

	xy.X = cx
	xy.Y = cy

	return xy

}

//Vars for defualt flags
//Built version will need changes!

var optList bool = false
var optFind string = "off"
var optAlert bool = false

func init() {
	flag.BoolVar(&optList, "list", optList, "List all Hue lights with ID and name")
	flag.StringVar(&optFind, "f", optFind, "Find Hue lights with the name value")
	flag.BoolVar(&optAlert, "alert", optAlert, "Blink lights")

	flag.Parse()
}

func changelightcolor(bridge *huego.Bridge, xy XYColor) {
	test1 := findlights(optFind, bridge)
	for i := range test1 {
		err := test1[i].Xy(xy.ConvToArray())
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

	RGB := RGBColor{R: 0, G: 0, B: 0}
	XY := RGB.ConvToXY()

	fmt.Println(XY.X, XY.Y)

	changelightcolor(bridge, XY)

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
}
