# huecli
A cross-platform command line Philips Hue client built with golang and huego.

## TODO

* refactor as package **NEXT**
* add groups
* add loggy based logging **DONE**
* add device registration **REVERTED**
* add brightness control **DONE**
* add hex color option **DONE**
* add named color option **DONE**
* add file-based saving of username/ip **DONE**
* clean code
* comment code

## SETUP

On first run huecli will look for a YAML file "settings.yml"

```
settings.yml format:

hueusername: Your hue username
hueip: Your hue IP
kbteam: Keybase team for output
kbchan: Keybase team channel for output, most likely - general -
prog: Program name for keybae output
 
```

## USAGE

```
Usage of ./huecli:
  -a    
        Blink lights
  -b uint
        Set light brightness (0-254) (default 255)
  -clr string
        Specify a color you want the light (red, green, blue, white)
  -f string
        Find Hue lights with the name value
  -hex string
        Specify a color you want the light in hex format (0F0F0F)
  -ls
        List all Hue lights with ID and name
  -rgb string
        Specify a color you want the light in format R-G-B (16-16-16)
  -s    
        Shutoff lights

```