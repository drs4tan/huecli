# huecli
A cross-platform command line Philips Hue client built with golang and huego.

## TODO

* add groups
* add device registration **DONE**
* add brightness control **DONE**
* add hex color option **DONE**
* add named color option **DONE**
* add file-based saving of username/ip **DONE**
* clean code
* comment code

## SETUP

On first run huecli will look for a file "username". If a file is not found it will ask you to either enter a username/ip or create a user. 

## USAGE

```
  -alert
        Blink lights
  -brightness uint
        Set light brightness (0-254) (default 255)
  -color string
        Specify a color you want the light (red, green, blue, white)
  -f string
        Find Hue lights with the name value
  -hex string
        Specify a color you want the light in hex format (0F0F0F)
  -list
        List all Hue lights with ID and name
  -rgb string
        Specify a color you want the light in format R-G-B (16-16-16)

```