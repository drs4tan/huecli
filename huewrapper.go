package huewrapper

import (
	"encoding/hex"
	"fmt"
	"github.com/amimof/huego"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
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

//Options to pass to run
type Options struct {
	OptOff        bool
	OptList       bool
	OptFind       string
	OptAlert      bool
	OptColorName  string
	OptColorRGB   string
	OptColorHEX   string
	OptBrightness uint
	OptTemp       uint
	OptDelay      uint
}

//HueHub to pass hue username and IP
type HueHub struct {
	IP 			string
	Username 	string
}


//Run takes the Options and HueHub to execute desired commands
func Run(opt Options, hub HueHub) {
	
	namedColors := map[string]RGBColor{
		"red":   RGBColor{255, 0, 0},
		"blue":  RGBColor{0, 0, 255},
		"green": RGBColor{0, 255, 0},
		"white": RGBColor{255, 255, 255},
	}

	bridge := huego.New(hub.IP, hub.Username)

	if opt.OptDelay > 0 {
		matchLights := findLights(opt.OptFind, bridge)
		for i := range matchLights {
			matchLights[i].TransitionTime(uint16(opt.OptDelay))
		}
	}

	if opt.OptOff == true {
		matchLights := findLights(opt.OptFind, bridge)
		for i := range matchLights {
			matchLights[i].Off()
		}
	}
	if opt.OptAlert == true {
		matchLights := findLights(opt.OptFind, bridge)
		for i := range matchLights {
			matchLights[i].Alert("select")
		}
	}
	if opt.OptColorRGB != "" {
		rgb := parseColorFlag(opt.OptColorRGB)
		rgbc := RGBColor{rgb[0], rgb[1], rgb[2]}
		xy := rgbc.ConvToXY()
		//log.LogInfo(fmt.Sprintf("X: %v Y: %v", xy.X, xy.Y))
		changeLightColor(bridge, xy, opt.OptFind)
	}

	if opt.OptColorHEX != "" {
		b, err := hex.DecodeString(opt.OptColorHEX)
		if err != nil {
			os.Exit(3)
		}
		RGBC := RGBColor{float32(b[0]), float32(b[1]), float32(b[2])}
		xy := RGBC.ConvToXY()

		changeLightColor(bridge, xy, opt.OptFind)
	}

	if opt.OptColorName != "" {
		RGBC := namedColor(opt.OptColorName, namedColors)
		xy := RGBC.ConvToXY()

		changeLightColor(bridge, xy, opt.OptFind)
	}

	if opt.OptBrightness != 0 {
		changeLightBrightness(bridge, opt.OptBrightness, opt.OptFind)
	}

	if opt.OptTemp > 0 {
		changeLightTemp(bridge, uint16(opt.OptTemp), opt.OptFind)
	}

	if opt.OptList == true {
		listLights(bridge)
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

func changeLightBrightness(bridge *huego.Bridge, level uint, optFind string) {
	matchLights := findLights(optFind, bridge)
	for i := range matchLights {
		err := matchLights[i].Bri(uint8(level))
		if err != nil {
			print("error")
		}

	}

}

func changeLightTemp(bridge *huego.Bridge, temp uint16, optFind string) {
	matchLights := findLights(optFind, bridge)
	for i := range matchLights {
		matchLights[i].Ct(temp)
		err := matchLights[i].Ct(temp)
		if err != nil {
			print("error")
		}
	}

}

func changeLightColor(bridge *huego.Bridge, xy XYColor, optFind string) {
	matchLights := findLights(optFind, bridge)
	for i := range matchLights {
		err := matchLights[i].Xy(xy.ConvToArray())
		if err != nil {
			print("error")
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
		print("Error")
	}

	var matchedLights []huego.Light
	for i := range allTheLights {
		if strings.Contains(allTheLights[i].Name, nameOfLight) {
			//log.LogInfo(fmt.Sprintf("Found light: '%v' (%v)", allTheLights[i].Name, allTheLights[i].ID))
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
	for i := range allTheLights {

		fmt.Printf("%v (%v) \t Mode: %v \t %v \n", allTheLights[i].Name, allTheLights[i].ID, allTheLights[i].State.ColorMode, allTheLights[i].State.Reachable)

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
