// Driver for cellular automata image output.
// Author: Matt Godshall
// Date  : 08-31-2013
package cell

import (
	"encoding/hex"
	"flag"
	"fmt"
	"image/color"
)

type colorError struct {
	message string
}

func (c *colorError) Error() string {
	return c.message
}

// Convert a color hex string (ex: "#9FEF00" or "9FEF00") to color.RGBA.
func parseColor(hexColor string) (color.RGBA, error) {
	if hexColor[0] == '#' {
		hexColor = hexColor[1:]
	}

	bytes, err := hex.DecodeString(hexColor)
	if err != nil {
		return color.RGBA{0, 0, 0, 0}, err
	}

	if len(bytes) != 3 {
		return color.RGBA{0, 0, 0, 0}, &colorError{"hex color string must be 3 bytes"}
	}

	rgbaColor := color.RGBA{bytes[0], bytes[1], bytes[2], 255}
	return rgbaColor, nil
}

func main() {
	var height, width int
	var rule uint
	var file, fg, bg string

	flag.IntVar(&height, "height", 768, "image height")
	flag.IntVar(&width, "width", 1024, "image width")
	flag.UintVar(&rule, "rule", 110, "cellular automata rule to generate")
	flag.StringVar(&file, "file", "ca.png", "image output file")
	flag.StringVar(&fg, "fg", "9FEF00", "foreground color")
	flag.StringVar(&bg, "bg", "000000", "background color")
	flag.Parse()

	if rule > 255 || rule < 0 {
		fmt.Println("rule must be between 0 and 255")
		return
	}

	fgRGBA, err := parseColor(fg)
	if err != nil {
		fmt.Println(err)
		return
	}

	bgRGBA, err := parseColor(bg)
	if err != nil {
		fmt.Println(err)
		return
	}

	cell := NewAutomaton(uint8(rule), width)
	cell.CreateImage(width, height, fgRGBA, bgRGBA, file)
}
