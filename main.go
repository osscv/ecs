package main

import (
	"fyne.io/fyne/v2/app"
)

func main() {
	myApp := app.NewWithID("com.oneclickvirt.goecs")
	myApp.Settings().SetTheme(&customTheme{})

	ui := NewTestUI(myApp)
	ui.window.ShowAndRun()
}
