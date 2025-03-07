package myapp

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	mt "github.com/Whadislov/TTCompanion/internal/my_types"
)

// FunctionalityPage creates the functionality page
func FunctionalityPage(db *mt.Database, w fyne.Window, a fyne.App) *fyne.Container {
	pageTitle := setTitle(T("functionalities"), 32)

	returnToMainMenuButton := widget.NewButton(T("return_to_main_page"), func() {
		MainPage(db, w, a)
	})

	createMenuButton := widget.NewButton(T("create_new_element"), func() { CreatePage(db, w, a) })
	createAddMenuButton := widget.NewButton(T("add_relationships"), func() { AddPage(db, w, a) })
	createRemoveMenuButton := widget.NewButton(T("remove_relationships"), func() { RemovePage(db, w, a) })
	createDeleteMenuButton := widget.NewButton(T("delete_element"), func() { DeletePage(db, w, a) })
	createAddInfoToPlayerButton := widget.NewButton(T("edit_player_information"), func() { AddInfoToPlayerPage(db, w, a) })

	functionalityPage := container.NewVBox(
		pageTitle,
		createMenuButton,
		createAddMenuButton,
		createRemoveMenuButton,
		createDeleteMenuButton,
		createAddInfoToPlayerButton,
		returnToMainMenuButton,
	)

	return functionalityPage
}
