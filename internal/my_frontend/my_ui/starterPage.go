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
func StarterPage() {
	log.Printf("Client : launching app ")
	a := app.NewWithID("com.onrender.TTCompanion")

	// Set language
	AddTranslationsFS(translations, "translation")

	w := a.NewWindow("TT Companion")
	w.Resize(fyne.NewSize(600, 400))
	w.CenterOnScreen() // Center the window on the monitor

	if appStartOption == "local" {
		icon, err := fyne.LoadResourceFromPath("Icon.png")
		if err != nil {
			log.Printf("Failed to load icon: %v", err)
		}
		a.SetIcon(icon)
		// Know if light mode is activated or not
		loadTheme(a)
	} else if appStartOption == "browser" {
		loadThemeWeb(a)
	}

	// Check persistence
	if appStartOption == "browser" {
		log.Printf("Client : checking persistence ")
		hasPersistence, db, id, err := mr.CheckPersistence()
		//hasPersistence, _, _, err := mr.CheckPersistence()
		if err != nil {
			log.Printf("Failed to check persistence: %v", err)
		} else {
			if hasPersistence {
				log.Printf("Client : persistence is on ")
				userOfSession = db.Users[id]
				// Change the theme here, because it is normally done on the auth page, which is bypassed by persistence
				if a.Settings().ThemeVariant() == 1 {
					lightTheme.IsActivated = true
				} else {
					darkTheme.IsActivated = true
				}
				MainPage(db, w, a)
				w.SetMainMenu(MainMenu(db, w, a))
				w.ShowAndRun()
				return
			}
		}
	}

	// Starter page
	log.Printf("Client : persistence is off ")
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
			w.SetContent(AuthentificationPage(w, a))

		} else if appStartOption == "browser" {
			// No fade because it blinks on the browser and the problem is not yet solved
			log.Println("Transitioning to the authentification page web")
			w.SetContent(AuthentificationPageWeb(w, a))
		}

	}()
	log.Println("Displaying welcome page")
	w.SetContent(starterPage)
	w.SetMainMenu(nil)
	w.ShowAndRun()
	return
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
