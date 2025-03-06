package myapp

import (
	"errors"
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	mt "github.com/Whadislov/TTCompanion/internal/my_types"
)

func AddInfoToPlayerPage(db *mt.Database, w fyne.Window, a fyne.App) {
	pageTitle := setTitle(T("edit_player_information"), 32)

	returnToFonctionalityPageButton := widget.NewButton(T("return_to_functionalities"), func() { w.SetContent(FunctionalityPage(db, w, a)) })

	// When the player is not yet selected, display nothing
	rightContent := container.NewVBox()

	sortedPlayers := sortMap(db.Players)
	playerButtons := []fyne.CanvasObject{}

	for _, p := range sortedPlayers {
		player := p.Value
		playerButton := widget.NewButton(fmt.Sprintf("%v %v", player.Firstname, player.Lastname), func() { AddInfoToSelectedPlayerPage(player, db, w, a) })
		playerButtons = append(playerButtons, playerButton)
	}

	content := container.NewVBox(
		pageTitle,
		returnToFonctionalityPageButton,
		container.NewGridWithColumns(
			2,
			container.NewVBox(playerButtons...),
			rightContent,
		),
	)

	w.SetContent(content)
}

func AddInfoToSelectedPlayerPage(p *mt.Player, db *mt.Database, w fyne.Window, a fyne.App) {
	pageTitle := setTitle(T("edit_player_information"), 32)

	returnToFonctionalityPageButton := widget.NewButton(T("return_to_functionalities"), func() { w.SetContent(FunctionalityPage(db, w, a)) })

	cancelButton := widget.NewButton(T("cancel"), func() { AddInfoToPlayerPage(db, w, a) })

	playerLabel := widget.NewLabel(fmt.Sprintf(T("you_have_selected")+" %v %v.", p.Firstname, p.Lastname))

	// Do not display -1
	var playerAgeStr string
	if p.Age == -1 {
		playerAgeStr = T("not_indicated")
	} else {
		playerAgeStr = strconv.Itoa(p.Age)
	}
	var playerRankingStr string
	if p.Ranking == -1 {
		playerRankingStr = T("not_indicated")
	} else {
		playerRankingStr = strconv.Itoa(p.Ranking)
	}

	// Do not display blank for material
	var material1 string
	if p.Material[0] == "" {
		material1 = T("not_indicated")
	} else {
		material1 = p.Material[0]
	}
	var material2 string
	if p.Material[1] == "" {
		material2 = T("not_indicated")
	} else {
		material2 = p.Material[1]
	}
	var material3 string
	if p.Material[2] == "" {
		material3 = T("not_indicated")
	} else {
		material3 = p.Material[2]
	}

	// Here are optional informations that can be added to the player
	ageEntry := widget.NewEntry()
	entryAgeHolder := playerAgeStr
	ageEntry.SetPlaceHolder(entryAgeHolder)

	rankingEntry := widget.NewEntry()
	entryRankingHolder := playerRankingStr
	rankingEntry.SetPlaceHolder(entryRankingHolder)

	forehandEntry := widget.NewEntry()
	entryForehandHolder := material1
	forehandEntry.SetPlaceHolder(entryForehandHolder)

	backhandEntry := widget.NewEntry()
	entryBackhandHolder := material2
	backhandEntry.SetPlaceHolder(entryBackhandHolder)

	bladeEntry := widget.NewEntry()
	entryBladeHolder := material3
	bladeEntry.SetPlaceHolder(entryBladeHolder)

	var isAgeModified bool
	var isRankingModified bool
	var isForehandModified bool
	var isBackhandModified bool
	var isBladeModified bool

	confirmButton := widget.NewButton(T("confirm"), func() {
		// Check player age
		if ageEntry.Text != "" {
			a, errAge := strconv.Atoi(ageEntry.Text)
			if errAge != nil {
				ageEntry.SetPlaceHolder(entryAgeHolder)
				// Check if the age is a number
				if !isNumbersOnly(ageEntry.Text) {
					dialog.ShowError(errors.New(T("err_age_must_be_number")), w)
					return
				} else {
					dialog.ShowError(errAge, w)
					return
				}
			} else {
				if a != -1 {
					isAgeModified = true
					p.SetPlayerAge(a)
				}
			}
		}

		// Check player ranking
		if rankingEntry.Text != "" {
			r, errRanking := strconv.Atoi(rankingEntry.Text)
			if errRanking != nil {
				rankingEntry.SetPlaceHolder(entryRankingHolder)
				// Check if the ranking is a number
				if !isNumbersOnly(rankingEntry.Text) {
					dialog.ShowError(errors.New(T("err_ranking_must_be_number")), w)
					return
				} else {
					dialog.ShowError(errRanking, w)
					return
				}
			} else {
				if r != -1 {
					isRankingModified = true
					p.SetPlayerRanking(r)
				}
			}
		}

		// Check player material
		if forehandEntry.Text != "" {
			forehandEntry.Text = standardizeSpaces(forehandEntry.Text)
			b, err := isValidMaterialName(forehandEntry.Text)
			if !b {
				dialog.ShowError(err, w)
				forehandEntry.SetPlaceHolder(entryForehandHolder)
				return
			} else {
				isForehandModified = true
				p.SetPlayerMaterial(forehandEntry.Text, p.Material[1], p.Material[2])
			}
		}
		if backhandEntry.Text != "" {
			backhandEntry.Text = standardizeSpaces(backhandEntry.Text)
			b, err := isValidMaterialName(backhandEntry.Text)
			if !b {
				dialog.ShowError(err, w)
				backhandEntry.SetPlaceHolder(entryBackhandHolder)
				return
			} else {
				isBackhandModified = true
				p.SetPlayerMaterial(p.Material[0], backhandEntry.Text, p.Material[2])
			}
		}
		if bladeEntry.Text != "" {
			bladeEntry.Text = standardizeSpaces(bladeEntry.Text)
			b, err := isValidMaterialName(bladeEntry.Text)
			if !b {
				dialog.ShowError(err, w)
				bladeEntry.SetPlaceHolder(entryBladeHolder)
				return
			} else {
				isBladeModified = true
				p.SetPlayerMaterial(p.Material[0], p.Material[1], bladeEntry.Text)
			}
		}

		// Has something changed ?
		if isAgeModified || isRankingModified || isForehandModified || isBackhandModified || isBladeModified {
			HasChanged = true
			dialog.ShowInformation(T("success"), fmt.Sprintf("%v %v "+T("has_been_modified"), p.Firstname, p.Lastname), w)
			AddInfoToPlayerPage(db, w, a)
		} else {
			dialog.ShowInformation(T("information"), fmt.Sprintf("%v %v "+T("has_not_been_modified"), p.Firstname, p.Lastname), w)
		}

	})

	sortedPlayers := sortMap(db.Players)
	playerButtons := []fyne.CanvasObject{}

	for _, p := range sortedPlayers {
		player := p.Value
		playerButton := widget.NewButton(fmt.Sprintf("%v %v", player.Firstname, player.Lastname), func() { AddInfoToSelectedPlayerPage(player, db, w, a) })
		playerButtons = append(playerButtons, playerButton)
	}

	leftContent := container.NewVBox(playerButtons...)

	formLayout := container.New(layout.NewFormLayout(),
		widget.NewLabel(T("age")+":"), ageEntry,
		widget.NewLabel(T("ranking")+":"), rankingEntry,
		widget.NewLabel(T("forehand")+":"), forehandEntry,
		widget.NewLabel(T("backhand")+":"), backhandEntry,
		widget.NewLabel(T("blade")+":"), bladeEntry,
	)

	rightContent := container.NewVBox(
		playerLabel,
		formLayout,
		container.NewGridWithColumns(2, confirmButton, cancelButton),
	)

	content := container.NewVBox(
		pageTitle,
		returnToFonctionalityPageButton,
		container.NewGridWithColumns(
			2,
			leftContent,
			rightContent,
		),
	)
	w.SetContent(content)
}
