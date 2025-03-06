package myapp

import (
	mt "github.com/Whadislov/TTCompanion/internal/my_types"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// RemovePage sets up the main page for removing players from teams and vice versa.
func RemovePage(db *mt.Database, w fyne.Window, a fyne.App) {

	pageTitle := setTitle(T("remove_relationships"), 32)

	ReturnToFonctionalityPageButton := widget.NewButton(T("return_to_functionalities"), func() {
		fonctionalityPage := FunctionalityPage(db, w, a)
		w.SetContent(fonctionalityPage)
	})

	removeTfromPButton := widget.NewButton(T("remove_team_from_a_player"), func() {
		w.SetContent(
			currentSelectionPageTfromP(
				SelectionPageTfromP(db, w, a),
				nil,
				db, w, a,
			),
		)
	})

	removePfromTButton := widget.NewButton(T("remove_player_from_a_team"), func() {
		w.SetContent(
			currentSelectionPagePfromT(
				SelectionPagePfromT(db, w, a),
				nil,
				db, w, a,
			),
		)

	})

	removeCfromPButton := widget.NewButton(T("remove_club_from_a_player"), func() {
		w.SetContent(
			currentSelectionPagePfromT(
				SelectionPageCfromP(db, w, a),
				nil,
				db, w, a,
			),
		)

	})

	removePage := container.NewVBox(
		pageTitle,
		removeTfromPButton,
		removePfromTButton,
		removeCfromPButton,
		ReturnToFonctionalityPageButton,
	)

	w.SetContent(removePage)
}
