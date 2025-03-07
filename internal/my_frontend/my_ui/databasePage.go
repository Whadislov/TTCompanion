package myapp

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	mt "github.com/Whadislov/TTCompanion/internal/my_types"
)

// DatabasePage sets up the page for showing players, teams, and clubs.
func DatabasePage(db *mt.Database, w fyne.Window, a fyne.App) *fyne.Container {
	pageTitle := setTitle(T("your_database"), 32)

	returnToMainMenuButton := widget.NewButton(T("return_to_main_page"), func() {
		MainPage(db, w, a)
	})

	playerButton := widget.NewButton(T("your_players"), func() { PlayerPage(db, w, a) })
	teamButton := widget.NewButton(T("your_teams"), func() { TeamPage(db, w, a) })
	clubButton := widget.NewButton(T("your_clubs"), func() { ClubPage(db, w, a) })

	databasePage := container.NewVBox(
		pageTitle,
		playerButton,
		teamButton,
		clubButton,
		returnToMainMenuButton,
	)

	return databasePage
}
