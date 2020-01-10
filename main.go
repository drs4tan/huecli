package main

/*
TODO:
	add groups
	add device registration
	add brightness control
	add hex color option
	add named color option
	*add file-based saving of username/ip
	clean code
	comment code
*/

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/amimof/huego"
)

var red RGBColor = RGBColor{255, 0, 0}
var blue RGBColor = RGBColor{0, 0, 255}
var green RGBColor = RGBColor{0, 255, 0}

//Structs and type funcs

//RGBColor holds RGB color infomation
//Red, Green, Blue respecively
type RGBColor struct {
	R, G, B float32
}

//XYColor holds XY color infomation
//X, Y respecively
type XYColor struct {
	X, Y float32
}

//ConvToArray converts XY color to a float32 slice
//Might be depracated in future
func (xy XYColor) ConvToArray() []float32 {
	var xyarray []float32
	xyarray = append(xyarray, xy.X, xy.Y)
	return xyarray
}

//ConvToXY converts RGB color to XY color
//Returns XY Color
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
var optAlert bool = true
var optColorName string = ""
var optColorRGB string = ""
var optBrightness uint8 = 255

func init() {
	flag.BoolVar(&optList, "list", optList, "List all Hue lights with ID and name")
	flag.StringVar(&optFind, "f", optFind, "Find Hue lights with the name value")
	flag.BoolVar(&optAlert, "alrt", optAlert, "Blink lights")
	flag.StringVar(&optColorRGB, "color", optColorRGB, "Specify a color you want the light in format R-G-B")

	flag.Parse()
}

func main() {
	var uname, ip string
	if fileExists("username") {
		b, err := ioutil.ReadFile("username")
		if err != nil {
			fmt.Println("Woops:", err.Error())
			os.Exit(1)
		}
		bs := strings.Split(string(b), "/")

		fmt.Println("username:", bs[0])
		fmt.Println("ip:", bs[1])

		uname = bs[0]
		ip = bs[1]

	} else {

		createUsername(&uname, &ip)

	}

	bridge := huego.New(ip, uname)

	switch {
	case optAlert == true:
		test1 := findlights(optFind, bridge)
		for i := range test1 {
			test1[i].Alert("select")
		}
	case optList == true:
		listlights(bridge)

	case optColorRGB != "":
		rgb := parsecolorflag(optColorRGB)
		rgbc := RGBColor{rgb[0], rgb[1], rgb[2]}
		xy := rgbc.ConvToXY()

		fmt.Println(xy.X, xy.Y)
		changelightcolor(bridge, xy)

	}

}

func createUsername(uname *string, ip *string) {
	var optUname []byte
	println("No username/ip found! Would you like to create/find them or enter them? ")
	println("Valid options: >create >enter")
	fmt.Scanln(&optUname)

	if string(optUname) == "create" {

	} else {

		var unameb, ipb []byte

		println("please enter the Philips Hue username:")
		fmt.Scanln(&unameb)

		println("please enter the Philips Hue hub IP:")
		fmt.Scanln(&ipb)

		finalbs := string(unameb) + "/" + string(ipb)

		err := ioutil.WriteFile("username", []byte(finalbs), 0777)
		if err != nil {
			fmt.Println("Woops:", err.Error())
			os.Exit(2)
		}

		*uname = string(unameb)
		*ip = string(ipb)
	}

}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func parsecolorflag(flg string) []float32 {
	strS := strings.Split(flg, "-")
	var fltS []float32

	for i := range strS {
		flt, _ := strconv.ParseFloat(strS[i], 32)
		fltS = append(fltS, float32(flt))
	}
	return fltS

}

func changelightcolor(bridge *huego.Bridge, xy XYColor) {
	test1 := findlights(optFind, bridge)
	for i := range test1 {
		err := test1[i].Xy(xy.ConvToArray())
		if err != nil {
			fmt.Println("Woops: ", err.Error())
		}
	}
}

func findlights(nameOfLight string, bridge *huego.Bridge) (light []huego.Light) {
	allTheLights, err := bridge.GetLights()
	if err != nil {
		fmt.Println("Woops:", err.Error())
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

func listlights(bridge *huego.Bridge) {
	allTheLights, err := bridge.GetLights()
	if err != nil {
		fmt.Println("Woops:", err.Error())
	}
	println("Listing all lights:")
	for i := range allTheLights {
		fmt.Println(allTheLights[i].ID, allTheLights[i].Name)
	}
}

/*
This func is depracated
func findlight(nameOfLight string, bridge *huego.Bridge) (light huego.Light) {
	allTheLights, err := bridge.GetLights()
	if err != nil {
		fmt.Println("Woops:", err.Error())
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
*/
