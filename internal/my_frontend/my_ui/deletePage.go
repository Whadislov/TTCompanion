package myapp

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	mf "github.com/Whadislov/TTCompanion/internal/my_functions"
	mt "github.com/Whadislov/TTCompanion/internal/my_types"
	"github.com/google/uuid"
)

// DeletePage sets up the page for deleting players, teams, and clubs.
func DeletePage(db *mt.Database, w fyne.Window, a fyne.App) {
	var rebuildUI func()

	// Rebuild UI on modifications
	rebuildUI = func() {

		pageTitle := setTitle(T("delete_element"), 32)

		pLabel := widget.NewLabel(T("players"))
		tLabel := widget.NewLabel(T("teams"))
		cLabel := widget.NewLabel(T("clubs"))

		returnToFonctionalityPageButton := widget.NewButton(T("return_to_functionalities"), func() {
			fonctionalityPage := FunctionalityPage(db, w, a)
			w.SetContent(fonctionalityPage)
		})

		// "Sort" maps
		sortedPlayers := sortMap(db.Players)
		sortedTeams := sortMap(db.Teams)
		sortedClubs := sortMap(db.Clubs)

		// Players
		acp := widget.NewAccordion()

		for _, sortedPlayer := range sortedPlayers {
			// i is PlayerID
			i := sortedPlayer.Key
			p := db.Players[i]
			item := widget.NewAccordionItem(
				fmt.Sprintf("%v %v", p.Firstname, p.Lastname),
				container.NewVBox(
					PlayerInfos(p),
					widget.NewButton(T("delete"), func() {
						showConfirmationDialog(w, fmt.Sprintf(T("delete_player"), p.Firstname, p.Lastname), func() {
							err := mf.DeletePlayer(p, db)
							if err != nil {
								dialog.ShowError(err, w)
							} else {
								successMsg := fmt.Sprintf(T("delete_player_succes"), p.Firstname, p.Lastname)
								fmt.Println(successMsg)
								dialog.ShowInformation(T("success"), successMsg, w)

								// Set the flag to true to indicate that the database has changed
								HasChanged = true

								// Reload UI
								rebuildUI()
							}
						})
					}),
				),
			)
			acp.Append(item)
		}

		// Teams
		act := widget.NewAccordion()

		for _, sortedTeam := range sortedTeams {
			// i is TeamID
			i := sortedTeam.Key
			t := db.Teams[i]
			item := widget.NewAccordionItem(
				t.Name,
				container.NewVBox(
					TeamInfos(t),
					widget.NewButton(T("delete"), func() {
						showConfirmationDialog(w, fmt.Sprintf(T("delete_team"), t.Name), func() {
							err := mf.DeleteTeam(t, db)
							if err != nil {
								dialog.ShowError(err, w)
							} else {
								successMsg := fmt.Sprintf(T("delete_team_succes"), t.Name)
								fmt.Println(successMsg)
								dialog.ShowInformation(T("success"), successMsg, w)

								// Set the flag to true to indicate that the database has changed
								HasChanged = true

								// Reload UI
								rebuildUI()
							}
						})
					}),
				),
			)
			act.Append(item)
		}

		// Clubs
		acc := widget.NewAccordion()

		for _, sortedClub := range sortedClubs {
			// i is ClubID
			i := sortedClub.Key
			c := db.Clubs[i]
			item := widget.NewAccordionItem(
				c.Name,
				container.NewVBox(
					ClubInfos(c),
					widget.NewButton(T("delete"), func() {
						teamNames := ""
						for _, teamName := range c.TeamIDs {
							teamNames += teamName + ", "
						}
						// Remove extra ", "
						if len(teamNames) > 2 {
							teamNames = teamNames[:len(teamNames)-2]
						}

						showConfirmationDialog(w, fmt.Sprintf(T("delete_club"), c.Name, teamNames), func() {
							// Get teamIDs without the link with the club (to avoid slice modification while iterating)
							var teamIDs []uuid.UUID
							for teamID := range c.TeamIDs {
								teamIDs = append(teamIDs, teamID)
							}
							// Delete the inner teams as well
							for _, teamID := range teamIDs {
								err := mf.DeleteTeam(db.Teams[teamID], db)
								if err != nil {
									dialog.ShowError(err, w)
								}
							}
							// Delete the club
							err := mf.DeleteClub(c, db)
							if err != nil {
								dialog.ShowError(err, w)
							} else {
								successMsg := fmt.Sprintf(T("delete_club_succes"), c.Name)
								fmt.Println(successMsg)
								dialog.ShowInformation(T("success"), successMsg, w)

								// Set the flag to true to indicate that the database has changed
								HasChanged = true

								// Reload UI
								rebuildUI()
							}
						})
					}),
				),
			)
			acc.Append(item)
		}

		playerVBox := container.NewVBox(
			pLabel,
			acp)

		teamVBox := container.NewVBox(
			tLabel,
			act)

		clubVBox := container.NewVBox(
			cLabel,
			acc)

		w.SetContent(container.NewVBox(
			pageTitle,
			returnToFonctionalityPageButton,
			container.NewHBox(
				playerVBox,
				teamVBox,
				clubVBox),
		),
		)
	}

	// Build UI for the first time
	rebuildUI()

}
