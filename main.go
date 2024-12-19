package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"path/filepath"
)

type customTheme struct{}

func (t customTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return color.NRGBA{R: 25, G: 25, B: 25, A: 255}
	case theme.ColorNameForeground:
		return color.White
	case theme.ColorNameButton:
		return color.NRGBA{R: 0, G: 0, B: 0, A: 0}
	default:
		return theme.DefaultTheme().Color(name, variant)
	}
}

func (t customTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (t customTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (t customTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

var activeIcon fyne.CanvasObject

func createIconWrapper(imagePath string, isActive bool, id int) *fyne.Container {
	iconBg := canvas.NewRectangle(color.NRGBA{R: 30, G: 30, B: 30, A: 255})
	iconBg.SetMinSize(fyne.NewSize(80, 80))
	iconBg.CornerRadius = 15

	icon := canvas.NewImageFromFile(imagePath)
	icon.FillMode = canvas.ImageFillContain
	icon.SetMinSize(fyne.NewSize(60, 60))

	rectColor := color.NRGBA{R: 80, G: 80, B: 80, A: 255}
	if isActive {
		rectColor = color.NRGBA{R: 0, G: 255, B: 0, A: 255}
		activeIcon = canvas.NewRectangle(rectColor)
	}

	rect := canvas.NewRectangle(rectColor)
	rect.SetMinSize(fyne.NewSize(8, 8))
	rect.CornerRadius = 4

	tap := widget.NewButton("", func() {
		if activeIcon != nil {
			activeIcon.(*canvas.Rectangle).FillColor = color.NRGBA{R: 80, G: 80, B: 80, A: 255}
			activeIcon.Refresh()
		}
		rect.FillColor = color.NRGBA{R: 0, G: 255, B: 0, A: 255}
		rect.Refresh()
		activeIcon = rect
	})
	tap.Importance = widget.LowImportance

	iconContainer := container.NewMax(
		iconBg,
		container.NewCenter(icon),
		tap,
	)

	return container.NewVBox(
		container.NewCenter(iconContainer),
		container.NewCenter(rect),
	)
}

type customButton struct {
	widget.Button
}

func newCustomButton(text string, tapped func()) *customButton {
	button := &customButton{}
	button.Text = text
	button.OnTapped = tapped
	button.ExtendBaseWidget(button)
	return button
}

func (b *customButton) CreateRenderer() fyne.WidgetRenderer {
	b.Importance = widget.HighImportance
	return b.Button.CreateRenderer()
}

func appendLog(logs *widget.TextGrid, text string) {
	logs.SetText(logs.Text() + text + "\n")
	logs.Refresh()
}

var rightPanel *fyne.Container
var isLogsVisible = true

func main() {
	myApp := app.New()
	myApp.Settings().SetTheme(&customTheme{})
	myWindow := myApp.NewWindow("LockDown AntiAntiVM")
	myWindow.Resize(fyne.NewSize(800, 400))

	title := widget.NewLabelWithStyle("LockDown AntiAntiVM", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	assetsDir := "Assets"
	iconPaths := map[string]string{
		"lockdownbrowser": filepath.Join(assetsDir, "lockdownbrowser.jpg"),
		"safeexam":        filepath.Join(assetsDir, "safeexam.png"),
		"WIP":             filepath.Join(assetsDir, "WIP.png"),
	}

	iconsContainer := container.NewGridWithColumns(3,
		createIconWrapper(iconPaths["lockdownbrowser"], false, 1),
		createIconWrapper(iconPaths["safeexam"], false, 2),
		createIconWrapper(iconPaths["WIP"], false, 3),
	)

	logs := widget.NewTextGrid()
	logs.SetStyleRange(0, 0, 0, 100, &widget.CustomTextGridStyle{
		FGColor: color.White,
	})
	logs.SetText("Welcome to LockDown AntiAntiVM\n")

	scrollableLogs := container.NewScroll(logs)

	injectButton := newCustomButton("INJECT", func() {
		performInjection(logs)
	})

	terminalBg := canvas.NewRectangle(color.NRGBA{R: 30, G: 30, B: 30, A: 255})
	terminalBg.CornerRadius = 8
	terminalBg.SetMinSize(fyne.NewSize(40, 40))

	terminalIconPath := filepath.Join(assetsDir, "Terminal.png")
	terminalIcon := canvas.NewImageFromFile(terminalIconPath)
	terminalIcon.FillMode = canvas.ImageFillContain
	terminalIcon.SetMinSize(fyne.NewSize(40, 40))

	var terminalContainer *fyne.Container

	terminalTap := widget.NewButton("", func() {
		if isLogsVisible {
			rightPanel.Hide()
			isLogsVisible = false
			myWindow.SetFixedSize(false)
			myWindow.Resize(fyne.NewSize(0, 0))
			myWindow.SetFixedSize(true)
		} else {
			rightPanel.Show()
			isLogsVisible = true
			myWindow.SetFixedSize(false)
			myWindow.Resize(fyne.NewSize(800, 400))
			myWindow.SetFixedSize(true)
		}
		myWindow.Content().Refresh()
	})

	terminalContainer = container.NewMax(
		terminalBg,
		container.NewCenter(terminalIcon),
		terminalTap,
	)

	leftPanel := container.NewVBox(
		container.NewPadded(title),
		layout.NewSpacer(),
		iconsContainer,
		layout.NewSpacer(),
		injectButton,
		layout.NewSpacer(),
		container.NewHBox(layout.NewSpacer(), terminalContainer),
	)

	leftBg := canvas.NewRectangle(color.NRGBA{R: 15, G: 15, B: 15, A: 255})
	leftBg.CornerRadius = 12
	leftBg.SetMinSize(fyne.NewSize(200, 400))

	rightExteriorRect := canvas.NewRectangle(color.NRGBA{R: 15, G: 15, B: 15, A: 255})
	rightExteriorRect.CornerRadius = 12
	rightExteriorRect.SetMinSize(fyne.NewSize(530, 430))

	rightInteriorRect := canvas.NewRectangle(color.NRGBA{R: 30, G: 30, B: 30, A: 255})
	rightInteriorRect.CornerRadius = 8
	rightInteriorRect.SetMinSize(fyne.NewSize(500, 400))

	interiorContainer := container.NewMax(
		rightInteriorRect,
		container.NewPadded(scrollableLogs),
	)

	rightPanel = container.NewMax(
		rightExteriorRect,
		container.NewCenter(interiorContainer),
	)
	rightPanel.Resize(fyne.NewSize(600, 280))

	content := container.NewHBox(
		container.NewMax(leftBg, container.NewPadded(leftPanel)),
		rightPanel,
	)

	mainContainer := container.NewPadded(content)
	myWindow.SetContent(mainContainer)
	myWindow.ShowAndRun()
}

func performInjection(logs *widget.TextGrid) {
	appendLog(logs, "Injection started")
	//WIP
	appendLog(logs, "Injection completed")
}
