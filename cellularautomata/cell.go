// Output cellular automata to PNGs.
// Author: Matt Godshall
// Date  : 08-31-2013
package cell

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"strconv"
)

type CellularAutomaton struct {
	Cells      []uint8
	Rule       []uint8
	Generation int
}

// Initialize a cellular automaton.
func NewAutomaton(rule uint8, size int) *CellularAutomaton {
	cell := new(CellularAutomaton)
	rule_bits := strconv.FormatUint(uint64(rule), 2)
	cell.Cells = make([]uint8, size)
	cell.Cells[size/2] = 1
	cell.Rule = make([]uint8, 8)
	cell.Generation = 0

	bits := len(rule_bits)
	for i, j := 0, bits-1; j >= 0 && i < bits; i, j = i+1, j-1 {
		if rule_bits[j] == '0' {
			cell.Rule[i] = 0
		} else if rule_bits[j] == '1' {
			cell.Rule[i] = 1
		}
	}

	return cell
}

// Calculate the cell value in the next generation.
func (cell *CellularAutomaton) nextCell(left, middle, right uint8) uint8 {
	return cell.Rule[(left<<2)+(middle<<1)+(right<<0)]
}

// Simulate a generation.
func (cell *CellularAutomaton) Generate() {
	size := len(cell.Cells)
	new_cells := make([]uint8, size)
	for i, _ := range cell.Cells {
		left := cell.Cells[((size-1)+i)%size]
		middle := cell.Cells[i]
		right := cell.Cells[(i+1)%size]
		new_cells[i] = cell.nextCell(left, middle, right)
	}
	cell.Cells = new_cells
	cell.Generation++
}

// Create a PNG image of the cellular automaton.
func (cell *CellularAutomaton) CreateImage(width, height int, fg, bg color.RGBA, file string) {
	rect := image.Rect(0, 0, width, height)
	img := image.NewNRGBA(rect)

	for cell.Generation < height {
		for i := 0; i < width; i++ {
			if cell.Cells[i] == 0 {
				img.Set(i, cell.Generation, bg)
			} else {
				img.Set(i, cell.Generation, fg)
			}
		}
		cell.Generate()
	}

	handle, err := os.Create(file)
	if err != nil {
		fmt.Println("could not open file for writing")
		return
	}
	png.Encode(handle, img)
}
