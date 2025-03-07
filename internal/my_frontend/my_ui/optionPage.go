package myapp

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	mt "github.com/Whadislov/TTCompanion/internal/my_types"
)

// OptionPage sets up the option page, in which the user can change the theme color, the language
func OptionPage(db *mt.Database, w fyne.Window, a fyne.App) *fyne.Container {
	pageTitle := setTitle(T("options"), 32)

	themeButton := widget.NewButton(T("change_theme_color"), func() {
		returnToMainMenuButton := widget.NewButton(T("return_to_main_page"), func() {
			MainPage(db, w, a)
		})

		if darkTheme.IsActivated {
			a.Settings().SetTheme(&lightTheme)
			lightTheme.IsActivated = true
			darkTheme.IsActivated = false
			w.SetContent(container.NewVBox(OptionPage(db, w, a), returnToMainMenuButton))
		} else {
			a.Settings().SetTheme(&darkTheme)
			lightTheme.IsActivated = false
			darkTheme.IsActivated = true
			w.SetContent(container.NewVBox(OptionPage(db, w, a), returnToMainMenuButton))
		}
	})

	changeLanguageButton := widget.NewButton(T("change_language"), func() { w.SetContent(ChangeLanguagePage(db, w, a)) })

	optionPage := container.NewVBox(
		pageTitle,
		themeButton,
		changeLanguageButton,
	)

	return optionPage
}

func ChangeLanguagePage(db *mt.Database, w fyne.Window, a fyne.App) *fyne.Container {
	returnToMainMenuButton := widget.NewButton(T("return_to_main_page"), func() {
		MainPage(db, w, a)
	})

	languageSelector := widget.NewSelect([]string{"Deutsch", "English", "Français"}, func(selected string) {
		switch selected {
		case "English":
			setLanguage("en")
			currentSelectedLanguage = "English"
			// Refresh
			returnToMainMenuButton = widget.NewButton(T("return_to_main_page"), func() {
				MainPage(db, w, a)
			})
			w.SetContent(container.NewVBox(OptionPage(db, w, a), returnToMainMenuButton))
			w.SetMainMenu(MainMenu(db, w, a))
		case "Français":
			setLanguage("fr")
			currentSelectedLanguage = "Français"
			// Refresh
			returnToMainMenuButton = widget.NewButton(T("return_to_main_page"), func() {
				MainPage(db, w, a)
			})
			w.SetContent(container.NewVBox(OptionPage(db, w, a), returnToMainMenuButton))
			w.SetMainMenu(MainMenu(db, w, a))
		case "Deutsch":
			setLanguage("de")
			currentSelectedLanguage = "Deutsch"
			// Refresh
			returnToMainMenuButton = widget.NewButton(T("return_to_main_page"), func() {
				MainPage(db, w, a)
			})
			w.SetContent(container.NewVBox(OptionPage(db, w, a), returnToMainMenuButton))
			w.SetMainMenu(MainMenu(db, w, a))
		}
	})
	languageSelector.PlaceHolder = currentSelectedLanguage
	languageSelector.Alignment = fyne.TextAlignCenter
	//languageSelector.SetSelected(currentSelectedLanguage)

	changeLanguagePage := container.NewVBox(
		OptionPage(db, w, a),
		languageSelector,
		returnToMainMenuButton,
	)

	return changeLanguagePage

}
