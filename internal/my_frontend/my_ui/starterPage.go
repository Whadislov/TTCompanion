package myapp

import (
	"image/color"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	mr "github.com/Whadislov/TTCompanion/internal/my_frontend/my_requests"
)

// StarterPage creates the introduction page to the UI and the starter page
func StarterPage() fyne.App {
	a := app.NewWithID("com.onrender.TTCompanion")

	icon, err := fyne.LoadResourceFromPath("Icon.png")
	if err != nil {
		log.Printf("Failed to load icon: %v", err)
	}
	a.SetIcon(icon)

	// Set language
	AddTranslationsFS(translations, "translation")

	mainWindow := a.NewWindow("TT Companion")
	mainWindow.Resize(fyne.NewSize(600, 400))
	mainWindow.CenterOnScreen() // Center the window on the monitor

	if appStartOption == "local" {
		// Know if light mode is activated or not
		loadTheme(a)
	} else if appStartOption == "browser" {
		loadThemeWeb(a)
	}

	// Check persistence
	if appStartOption == "browser" {
		hasPersistence, db, id, err := mr.CheckPersistence()
		if err != nil {
			log.Printf("Failed to check persistence: %v", err)
		} else {
			if hasPersistence {
				userOfSession = db.Users[id]
				MainPage(db, mainWindow, a)
				mainWindow.SetMainMenu(MainMenu(db, mainWindow, a))
				return a
			}
		}
	}

	// Starter page
	pageTitle := setTitle(T("welcome_to_tt_companion"), 32)
	starterPage := container.NewCenter(pageTitle)

	// Fade
	go func() {
		time.Sleep(1 * time.Second)
		if appStartOption == "local" {
			themeColor := a.Settings().Theme().Color("foreground", a.Settings().ThemeVariant())
			fadeText(pageTitle, themeColor)
			// go to main page with delay so that the menu is not directly shown
			log.Println("Transitioning to identification page")
			mainWindow.SetContent(AuthentificationPage(mainWindow, a))

		} else if appStartOption == "browser" {
			// No fade because it blinks on the browser and the problem is not yet solved
			log.Println("Transitioning to the authentification page web")
			mainWindow.SetContent(AuthentificationPageWeb(mainWindow, a))
		}

	}()
	log.Println("Displaying welcome page")
	mainWindow.SetContent(starterPage)
	mainWindow.SetMainMenu(nil)
	mainWindow.ShowAndRun()
	return a
}

func fadeText(text *canvas.Text, textColor color.Color) {
	r, g, b, alp := textColor.RGBA()
	var fadeStep uint8 = 5
	var threshold uint8 = 120

	// >> 8 because color.RGBA can only use values of 8 bits (textColor is 16 bits)
	updateUI := func(alpha uint8) {
		text.Color = color.RGBA{
			R: uint8(r >> 8),
			G: uint8(g >> 8),
			B: uint8(b >> 8),
			A: alpha,
		}
		text.Refresh()
	}

	for alpha := uint8(alp >> 8); alpha >= threshold; alpha -= fadeStep {
		updateUI(alpha)
		text.Refresh()
		time.Sleep(20 * time.Millisecond) // Pause to simulate fade
	}

}
