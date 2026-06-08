package ui

import (
	"fmt"
	"strings"

	"cot/internal/currency"
)

// BrailleCanvas represents a terminal sub-pixel drawing grid using Unicode Braille patterns (2x4 dots per character)
type BrailleCanvas struct {
	Width      int // Dot width (CharWidth * 2)
	Height     int // Dot height (CharHeight * 4)
	CharWidth  int // Width in character columns
	CharHeight int // Height in character rows
	Grid       [][]byte
}

// NewBrailleCanvas initializes a new canvas with character dimensions
func NewBrailleCanvas(charWidth, charHeight int) *BrailleCanvas {
	w := charWidth * 2
	h := charHeight * 4
	grid := make([][]byte, charHeight)
	for i := range grid {
		grid[i] = make([]byte, charWidth)
	}
	return &BrailleCanvas{
		Width:      w,
		Height:     h,
		CharWidth:  charWidth,
		CharHeight: charHeight,
		Grid:       grid,
	}
}

// Set illuminates a specific sub-pixel dot (0,0 is bottom-left, Width-1, Height-1 is top-right)
func (c *BrailleCanvas) Set(dotX, dotY int) {
	if dotX < 0 || dotX >= c.Width || dotY < 0 || dotY >= c.Height {
		return
	}

	// Invert Y because grid rendering goes top-to-bottom, but charts plot bottom-to-top
	invertedY := c.Height - 1 - dotY

	charX := dotX / 2
	charY := invertedY / 4

	localX := dotX % 2
	localY := invertedY % 4

	// Mapping coordinates (localX, localY) in a 2x4 cell to Unicode Braille bit patterns:
	// localX=0 (left column):  y=0 -> bit 0 (0x01), y=1 -> bit 1 (0x02), y=2 -> bit 2 (0x04), y=3 -> bit 6 (0x40)
	// localX=1 (right column): y=0 -> bit 3 (0x08), y=1 -> bit 4 (0x10), y=2 -> bit 5 (0x20), y=3 -> bit 7 (0x80)
	var bit byte
	if localX == 0 {
		switch localY {
		case 0:
			bit = 0x01
		case 1:
			bit = 0x02
		case 2:
			bit = 0x04
		case 3:
			bit = 0x40
		}
	} else {
		switch localY {
		case 0:
			bit = 0x08
		case 1:
			bit = 0x10
		case 2:
			bit = 0x20
		case 3:
			bit = 0x80
		}
	}

	c.Grid[charY][charX] |= bit
}

// DrawLine draws a line segment using Bresenham's algorithm
func (c *BrailleCanvas) DrawLine(x1, y1, x2, y2 int) {
	dx := abs(x2 - x1)
	dy := abs(y2 - y1)
	sx := -1
	if x1 < x2 {
		sx = 1
	}
	sy := -1
	if y1 < y2 {
		sy = 1
	}
	err := dx - dy

	for {
		c.Set(x1, y1)
		if x1 == x2 && y1 == y2 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x1 += sx
		}
		if e2 < dx {
			err += dx
			y1 += sy
		}
	}
}

// GetCellRune returns the Unicode Braille pattern rune for a given character cell coordinate
func (c *BrailleCanvas) GetCellRune(charX, charY int) rune {
	return rune(0x2800) + rune(c.Grid[charY][charX])
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// DrawChart draws a beautiful high-resolution line chart using Braille canvas
func DrawChart(prices []float64, symbol, name string, currencySymbol string, currentBRL float64) int {
	charWidth := 50
	charHeight := 10

	canvas := NewBrailleCanvas(charWidth, charHeight)

	// Since each character holds 2 dots horizontally, the chart can support canvas.Width dots
	maxDotsX := canvas.Width

	numPrices := len(prices)
	startIndex := 0
	if numPrices > maxDotsX {
		startIndex = numPrices - maxDotsX
	}
	visiblePrices := prices[startIndex:]

	// Find min and max
	min := visiblePrices[0]
	max := visiblePrices[0]
	for _, p := range visiblePrices {
		if p < min {
			min = p
		}
		if p > max {
			max = p
		}
	}

	// Add padding to prevent division by zero
	if min == max {
		min = min * 0.999
		max = max * 1.001
	}

	priceRange := max - min
	linesPrinted := 0

	// Draw compact header
	var title string
	if currencySymbol == "R$" {
		title = fmt.Sprintf("=== %s (%s) ➔ BRL | Ctrl+C para sair ===", strings.ToUpper(name), strings.ToUpper(symbol))
	} else {
		title = fmt.Sprintf("=== %s (%s) ➔ USD/BRL | Ctrl+C para sair ===", strings.ToUpper(name), strings.ToUpper(symbol))
	}
	fmt.Printf("%s%s%s%s\033[K\n", Cyan, Bold, title, Reset)
	linesPrinted++
	fmt.Printf("\033[K\n") // Empty spacer line
	linesPrinted++

	// Map index to X dot coordinate (right-aligned if history is short)
	getX := func(i int) int {
		if numPrices < maxDotsX {
			return (maxDotsX - numPrices) + i
		}
		return i
	}

	// Map price value to Y dot coordinate
	getY := func(p float64) int {
		return int((p - min) / priceRange * float64(canvas.Height-1))
	}

	// Draw connected segments in the dot space
	for i := 1; i < len(visiblePrices); i++ {
		x1 := getX(i - 1)
		y1 := getY(visiblePrices[i-1])
		x2 := getX(i)
		y2 := getY(visiblePrices[i])
		canvas.DrawLine(x1, y1, x2, y2)
	}

	// Render the character grid row-by-row
	for cy := 0; cy < charHeight; cy++ {
		level := max - priceRange*(float64(cy)/float64(charHeight-1))

		var label string
		if currencySymbol == "R$" {
			label = currency.FormatCurrencySmart(level, "R$", 4)
		} else {
			label = currency.FormatCurrencySmart(level, "US$", 2)
		}

		fmt.Printf("%s%15s │%s ", Gray, label, Reset)

		// Print columns cell-by-cell
		for cx := 0; cx < charWidth; cx++ {
			// Decide color based on price trend in this character cell column
			color := Reset
			dotX := cx * 2

			// Find corresponding price index
			priceIdx := dotX
			if numPrices < maxDotsX {
				priceIdx = dotX - (maxDotsX - numPrices)
			}

			if priceIdx >= 0 && priceIdx < len(visiblePrices) {
				prevIdx := priceIdx - 1
				if prevIdx >= 0 {
					if visiblePrices[priceIdx] > visiblePrices[prevIdx] {
						color = Green
					} else if visiblePrices[priceIdx] < visiblePrices[prevIdx] {
						color = Red
					}
				}
			}

			// Retrieve the unicode braille rune
			r := canvas.GetCellRune(cx, cy)
			if r == 0x2800 {
				// Empty braille cell: render space
				fmt.Print(" ")
			} else {
				fmt.Printf("%s%c%s", color, r, Reset)
			}
		}
		fmt.Print("\033[K\n")
		linesPrinted++
	}

	// Draw X-axis border
	fmt.Printf("%s%15s └%s%s\033[K\n", Gray, "", strings.Repeat("─", charWidth+2), Reset)
	linesPrinted++

	// Print stats in a single compact line below the chart
	currentPrice := visiblePrices[len(visiblePrices)-1]
	var formattedCurrent string
	if currencySymbol == "R$" {
		formattedCurrent = currency.FormatCurrencySmart(currentPrice, "R$", 4)
	} else {
		formattedCurrent = currency.FormatCurrencySmart(currentPrice, "US$", 2)
		if currentBRL > 0 {
			formattedBRL := currency.FormatCurrencySmart(currentBRL, "R$", 2)
			formattedCurrent = fmt.Sprintf("%s (%s)", formattedCurrent, formattedBRL)
		}
	}

	minStr := formatPriceSymbol(min, currencySymbol)
	maxStr := formatPriceSymbol(max, currencySymbol)

	fmt.Printf("Preço Atual: %s%s%s | Mínimo: %s | Máximo: %s\033[K\n",
		Bold, formattedCurrent, Reset, minStr, maxStr)
	linesPrinted++

	return linesPrinted
}

func formatPriceSymbol(val float64, currencySymbol string) string {
	if currencySymbol == "R$" {
		return currency.FormatCurrencySmart(val, "R$", 4)
	}
	return currency.FormatCurrencySmart(val, "US$", 2)
}
