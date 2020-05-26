# huecli
A cross-platform command line Philips Hue client built with golang and huego.

## TODO

* refactor as package **NEXT**
* add groups **NOT STARTED**
* add device registration **REVERTED**
* add brightness control **DONE**
* add hex color option **DONE**
* add named color option **DONE**
* clean code **ON GOING**
* comment code **ON GOING**




## GET STARTED
This package currently relies on github.com/amimof/huego to work.

To install huego run:
```shell
go get -u github.com/amimof/huego
```


A simple example to flash the lights at set times:
```go
package main

import (
	"huewrapper"
	"time"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
)
var hub huewrapper.HueHub
func main() {
	k := koanf.New("/")
	if err := k.Load(file.Provider("settings.yml"), yaml.Parser()); err != nil {
		print("error mate")
	}
	hub = huewrapper.HueHub{IP: k.String("hueip"), Username: k.String("hueusername")}
	stp := make(chan uint)
	go timeChange(30, stp, 4)
	<- stp
	print("done")
}


func timeChange(mins int, st chan uint, times int) {
	alert := mins / times
	opt := huewrapper.Options{OptFind: "off", OptAlert: true}
	huewrapper.Run(opt, hub)
	
	for times >= 0 {
		time.Sleep(time.Duration(alert) * time.Minute)
		huewrapper.Run(opt, hub)
		times--
	}
	st <- 1
	

}
```