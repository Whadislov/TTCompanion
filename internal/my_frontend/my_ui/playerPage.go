package myapp

import (
	"fmt"
	"slices"
	"strconv"

	mt "github.com/Whadislov/TTCompanion/internal/my_types"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// PlayerInfos returns a container that displays all the infos of a player.
func PlayerInfos(player *mt.Player) *fyne.Container {

	teams := []string{}
	clubs := []string{}

	for _, team := range player.TeamIDs {
		teams = append(teams, team)
	}
	// Sort teams alphabetically
	slices.Sort(teams)

	for _, club := range player.ClubIDs {
		clubs = append(clubs, club)
	}
	// Sort clubs alphabetically
	slices.Sort(clubs)

	// Do not display -1
	var playerAgeStr string
	if player.Age == -1 {
		playerAgeStr = T("not_indicated")
	} else {
		playerAgeStr = strconv.Itoa(player.Age)
	}
	var playerRankingStr string
	if player.Ranking == -1 {
		playerRankingStr = T("not_indicated")
	} else {
		playerRankingStr = strconv.Itoa(player.Ranking)
	}

	leftText := T("firstname") + ":\n" + T("lastname") + ":\n" + T("age") + ":\n" + T("forehand") + ":\n" + T("backhand") + ":\n" + T("blade") + ":\n" + T("ranking") + ":\n" + T("teams") + ":\n" + T("clubs") + ":"
	rightText := player.Firstname + "\n" + player.Lastname + "\n" + playerAgeStr + "\n" + player.Material[0] + "\n" + player.Material[1] + "\n" + player.Material[2] + "\n" + playerRankingStr + "\n" + strHelper(teams) + "\n" + strHelper(clubs)

	playerInfosContent := container.NewGridWithColumns(2, widget.NewLabel(leftText), widget.NewLabel(rightText))
	return playerInfosContent
}

// PlayerPage sets up the page for displaying player info.
func PlayerPage(db *mt.Database, w fyne.Window, a fyne.App) {

	pageTitle := setTitle(T("players"), 32)
	ac := widget.NewAccordion()

	// "Sort the map"
	sortedPlayers := sortMap(db.Players)

	for _, player := range sortedPlayers {
		item := widget.NewAccordionItem(
			fmt.Sprintf("%v %v", player.Value.Firstname, player.Value.Lastname),
			PlayerInfos(player.Value),
		)

		ac.Append(item)
	}

	returnToDatabasePageButton := widget.NewButton(T("return_to_database"), func() {
		databasePage := DatabasePage(db, w, a)
		w.SetContent(databasePage)
	})

	w.SetContent(container.NewVBox(
		pageTitle,
		returnToDatabasePageButton,
		ac),
	)
}
