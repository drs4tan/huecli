package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/amimof/huego"
)

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
var optOff bool = false
var optList bool = false
var optFind string = ""
var optAlert bool = false
var optColorName string = ""
var optColorRGB string = ""
var optColorHEX string = ""
var optBrightness uint = 255

func init() {
	flag.BoolVar(&optOff, "s", optOff, "Shutoff lights")
	flag.BoolVar(&optList, "ls", optList, "List all Hue lights with ID and name")
	flag.StringVar(&optFind, "f", optFind, "Find Hue lights with the name value")
	flag.BoolVar(&optAlert, "alert", optAlert, "Blink lights")
	flag.StringVar(&optColorRGB, "rgb", optColorRGB, "Specify a color you want the light in format R-G-B (16-16-16)")
	flag.StringVar(&optColorHEX, "hex", optColorHEX, "Specify a color you want the light in hex format (0F0F0F)")
	flag.StringVar(&optColorName, "color", optColorName, "Specify a color you want the light (red, green, blue, white)")
	flag.UintVar(&optBrightness, "brightness", optBrightness, "Set light brightness (0-254)")

	flag.Parse()
}

func main() {
	namedColors := map[string]RGBColor{
		"red":   RGBColor{255, 0, 0},
		"blue":  RGBColor{0, 0, 255},
		"green": RGBColor{0, 255, 0},
		"white": RGBColor{255, 255, 255},
	}

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

	case optOff == true:
		matchLights := findLights(optFind, bridge)
		for i := range matchLights {
			matchLights[i].Off()
		}
	case optAlert == true:
		matchLights := findLights(optFind, bridge)
		for i := range matchLights {
			matchLights[i].Alert("select")
		}

	case optList == true:
		listLights(bridge)

	case optColorRGB != "":
		rgb := parseColorFlag(optColorRGB)
		rgbc := RGBColor{rgb[0], rgb[1], rgb[2]}
		xy := rgbc.ConvToXY()

		fmt.Println(xy.X, xy.Y)
		changeLightColor(bridge, xy)

	case optColorHEX != "":
		b, err := hex.DecodeString(optColorHEX)
		if err != nil {
			fmt.Println("Woops:", err.Error())
			os.Exit(3)
		}
		RGBC := RGBColor{float32(b[0]), float32(b[1]), float32(b[2])}
		xy := RGBC.ConvToXY()

		changeLightColor(bridge, xy)

	case optColorName != "":
		RGBC := namedColor(optColorName, namedColors)
		xy := RGBC.ConvToXY()

		changeLightColor(bridge, xy)

	case optBrightness != 255:
		changeLightBrightness(bridge, optBrightness)
	}

}

func createUsername(uname *string, ip *string) {
	var optUname []byte
	println("No username/ip found! Would you like to create/find them or enter them? ")
	println("Valid options: >create >enter")
	fmt.Scanln(&optUname)

	if string(optUname) == "create" {
		println("Press the link button on the Philips Hue hub, then press enter")
		fmt.Scanln(&optUname)
		bridge, _ := huego.Discover()
		user, err := bridge.CreateUser("huecli") // Link button needs to be pressed
		if err != nil {
			fmt.Printf("Woops: %s", err.Error())
			os.Exit(2)
		}

		bridge = bridge.Login(user)
		finalstr := user + "/" + bridge.Host
		println(finalstr)
		err = ioutil.WriteFile("username", []byte(finalstr), 0777)
		if err != nil {
			fmt.Println("Woops:", err.Error())
			os.Exit(2)
		}

		*uname = string(user)
		*ip = string(bridge.Host)

	} else {

		var unameb, ipb []byte

		println("please enter the Philips Hue username:")
		fmt.Scanln(&unameb)

		println("please enter the Philips Hue hub IP:")
		fmt.Scanln(&ipb)

		finalstr := string(unameb) + "/" + string(ipb)

		err := ioutil.WriteFile("username", []byte(finalstr), 0777)
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

func parseColorFlag(flg string) []float32 {
	strS := strings.Split(flg, "-")
	var fltS []float32

	for i := range strS {
		flt, _ := strconv.ParseFloat(strS[i], 32)
		fltS = append(fltS, float32(flt))
	}
	return fltS

}

func changeLightBrightness(bridge *huego.Bridge, level uint) {
	matchLights := findLights(optFind, bridge)
	for i := range matchLights {
		err := matchLights[i].Bri(uint8(level))
		if err != nil {
			fmt.Println("Woops: ", err.Error())
		}
	}

}

func changeLightColor(bridge *huego.Bridge, xy XYColor) {
	matchLights := findLights(optFind, bridge)
	for i := range matchLights {
		err := matchLights[i].Xy(xy.ConvToArray())
		if err != nil {
			fmt.Println("Woops: ", err.Error())
		}
	}
}

func namedColor(colorstr string, named map[string]RGBColor) (color RGBColor) {
	_, exists := named[colorstr]
	if exists == true {
		return named[colorstr]
	}

	fmt.Printf("Woops: No named color of %s defined", colorstr)
	os.Exit(4)
	return
}

func findLights(nameOfLight string, bridge *huego.Bridge) (light []huego.Light) {
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

func listLights(bridge *huego.Bridge) {
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
