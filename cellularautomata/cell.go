// Output cellular automata to PNGs.
// Author: Matt Godshall
// Date  : 08-31-2013
package main

import (
  "fmt"
  "strconv"
  "os"
  "encoding/hex"
  "flag"
  "image"
  "image/color"
  "image/png"
)

type cellularAutomaton struct {
  Cells      []uint8
  Rule       []uint8
  Generation int
}

// Initialize a cellular automaton.
func newAutomaton(rule uint8, size int) *cellularAutomaton {
  cell := new(cellularAutomaton)
  rule_bits := strconv.FormatUint(uint64(rule), 2)
  cell.Cells = make([]uint8, size)
  cell.Cells[size / 2] = 1
  cell.Rule  = make([]uint8, 8)
  cell.Generation = 0

  bits := len(rule_bits)
  for i , j := 0, bits - 1; j >= 0 && i < bits; i , j = i + 1, j - 1 {
    if rule_bits[j] == '0' {
      cell.Rule[i] = 0
    } else if rule_bits[j] == '1' {
      cell.Rule[i] = 1
    }
  }

  return cell
}

// Calculate the cell value in the next generation.
func (cell *cellularAutomaton) nextCell(left, middle, right uint8) uint8 {
  return cell.Rule[(left << 2) + (middle << 1) + (right << 0)]
}

// Simulate a generation.
func (cell *cellularAutomaton) generate() {
  size := len(cell.Cells)
  new_cells := make([]uint8, size)
  for i, _ := range cell.Cells {
    left   := cell.Cells[((size - 1) + i) % size]
    middle := cell.Cells[i]
    right  := cell.Cells[(i + 1) % size]
    new_cells[i] = cell.nextCell(left, middle, right)
  }
  cell.Cells = new_cells
  cell.Generation++
}

// Create a PNG image from the cellular automaton.
func (cell *cellularAutomaton) createImage(width, height int, fg, bg color.RGBA, file string) {
  rect := image.Rect(0, 0, width, height)
  img := image.NewNRGBA(rect)

  fore := color.RGBA{0x9F, 0xEF, 0x00, 255}
  back := color.RGBA{0, 0, 0, 255}

  for cell.Generation < height {
    for i := 0; i < width; i++ {
      if cell.Cells[i] == 0 {
        img.Set(i, cell.Generation, back)
      } else {
        img.Set(i, cell.Generation, fore)
      }
    }
    cell.generate()
  }

  handle, err := os.Create(file)
  if err != nil {
    fmt.Println("could not open file for writing")
    return
  }
  png.Encode(handle, img)
}

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
    return color.RGBA{0, 0, 0, 0}, &colorError{ "hex color string must be 3 bytes" }
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

  if (rule > 255 || rule < 0) {
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

  cell := newAutomaton(uint8(rule), width)
  cell.createImage(width, height, fgRGBA, bgRGBA, file)
}
