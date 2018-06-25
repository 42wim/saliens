package main

import (
	"strings"

	"github.com/gizak/termui"
)

func printStatus(text string) {
	printText("Status", text, 40, 0, 12)
}

func printScore(text string) {
	printText("Score", text, 10, 41, 6)
}

func printLevel(text string) {
	printText("Level", text, 10, 41, 9)
}

func printPlanet(text string) {
	if text != "?" {
		text = pMap[text].State.Name
		text = strings.Replace(text, "#TerritoryControl_", "", -1)
		//	text = text + strconv.Itoa(pMap[text].State.Difficulty)
	}
	printText("Planet", text, 40, 0, 0)
}

func printZone(text string) {
	printText("Zone", text, 9, 0, 3)
}

func printDifficulty(text string) {
	printText("Difficulty", text, 9, 10, 3)
}

func printZonesLeft(text string) {
	printText("Left", text, 9, 20, 3)
}

func printZoneCapture(text string) {
	printText("Capture", text, 10, 30, 3)
}

func printCapture(text string) {
	printText("Capture", text, 9, 41, 0)
}

func printNextGrind(text string) {
	printText("ETA", text, 9, 52, 6)
}

func printNextLevel(text string) {
	printText("ETA", text, 9, 52, 9)
}

func printText(label string, text string, width int, x, y int) {
	par := termui.NewPar(text)
	par.Height = 3
	par.Width = width
	par.Y = y
	par.X = x
	par.BorderLabel = label
	par.BorderFg = termui.ColorWhite
	termui.Render(par)
}

func updateGauge(percent int) {
	printGauge("Grinding XP", percent, 40, 0, 6)
}

func updateNextLevelGauge(percent int) {
	printGauge("Next Level", percent, 40, 0, 9)
}

func printGauge(label string, percent, width, x, y int) {
	g0 := termui.NewGauge()
	g0.Percent = percent
	g0.Y = y
	g0.X = x
	g0.Width = width
	g0.Height = 3
	g0.BorderLabel = label
	g0.BarColor = termui.ColorRed
	g0.BorderFg = termui.ColorWhite
	g0.BorderLabelFg = termui.ColorCyan
	termui.Render(g0)
}
