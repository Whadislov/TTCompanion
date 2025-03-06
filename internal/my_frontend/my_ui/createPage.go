package myapp

import (
	"errors"
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	mf "github.com/Whadislov/TTCompanion/internal/my_functions"
	mt "github.com/Whadislov/TTCompanion/internal/my_types"
)

// CreatePage sets up the page for creating players, teams, and clubs.
func CreatePage(db *mt.Database, w fyne.Window, a fyne.App) {

	pageTitle := setTitle(T("create_new_element"), 32)

	ReturnToFonctionalityPageButton := widget.NewButton(T("return_to_functionalities"), func() {
		fonctionalityPage := FunctionalityPage(db, w, a)
		w.SetContent(fonctionalityPage)
	})

	ReturnToCreatePageButton := widget.NewButton(T("return_to_the_create_menu"), func() {
		CreatePage(db, w, a)
	})

	// Player
	playerButton := widget.NewButton(T("create_a_new_player"), func() {

		// Club Selection
		pageTitle := setTitle(T("create_a_new_player_select_a_club"), 32)

		// clubSelectionPage
		clubSelectionPageButton := widget.NewButton(T("select_a_club"), func() {
			pageTitle := setTitle(T("create_a_new_player_select_a_club"), 32)
			listOfClubs := []fyne.CanvasObject{}

			// Sort clubs for an alphabetical order button display
			sortedClubs := sortMap(db.Clubs)

			for _, c := range sortedClubs {
				club := c.Value
				clubButton := widget.NewButton(club.Name, func() {
					// After club selection
					clubLabel := widget.NewLabel(fmt.Sprintf(T("you_are_going_to_create_a_player_for"), club.Name))

					// We can now create the player
					firstnameEntry := widget.NewEntry()
					entryFirstnameHolder := T("firstname") + " ..."
					firstnameEntry.SetPlaceHolder(entryFirstnameHolder)

					lastnameEntry := widget.NewEntry()
					entryLastnameHolder := T("lastname") + " ..."
					lastnameEntry.SetPlaceHolder(entryLastnameHolder)

					// Here are optional informations that can be added to the player
					ageEntry := widget.NewEntry()
					entryAgeHolder := T("age") + " ..."
					ageEntry.SetPlaceHolder(entryAgeHolder)

					rankingEntry := widget.NewEntry()
					entryRankingHolder := T("ranking") + " ..."
					rankingEntry.SetPlaceHolder(entryRankingHolder)

					forehandEntry := widget.NewEntry()
					entryForehandHolder := T("forehand") + " ..."
					forehandEntry.SetPlaceHolder(entryForehandHolder)

					backhandEntry := widget.NewEntry()
					entryBackhandHolder := T("backhand") + " ..."
					backhandEntry.SetPlaceHolder(entryBackhandHolder)

					bladeEntry := widget.NewEntry()
					entryBladeHolder := T("blade") + " ..."
					bladeEntry.SetPlaceHolder(entryBladeHolder)

					validatationButton := widget.NewButton(T("confirm"), func() {
						age := -1
						ranking := -1

						bf, errf := mf.IsValidName(firstnameEntry.Text)
						if !bf {
							dialog.ShowError(errf, w)
							return
						}

						bl, errl := mf.IsValidName(lastnameEntry.Text)
						if !bl {
							dialog.ShowError(errl, w)
							return
						}

						// Set player age
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
								age = a
							}
						}

						// Set player ranking
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
								ranking = r
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
							}
						} else {
							// Set default player material
							forehandEntry.SetText(T("not_indicated"))
						}
						if backhandEntry.Text != "" {
							backhandEntry.Text = standardizeSpaces(backhandEntry.Text)
							b, err := isValidMaterialName(backhandEntry.Text)
							if !b {
								dialog.ShowError(err, w)
								backhandEntry.SetPlaceHolder(entryBackhandHolder)
								return
							}
						} else {
							// Set default player material
							backhandEntry.SetText(T("not_indicated"))
						}
						if bladeEntry.Text != "" {
							bladeEntry.Text = standardizeSpaces(bladeEntry.Text)
							b, err := isValidMaterialName(bladeEntry.Text)
							if !b {
								dialog.ShowError(err, w)
								bladeEntry.SetPlaceHolder(entryBladeHolder)
								return
							}
						} else {
							// Set default player material
							bladeEntry.SetText(T("not_indicated"))
						}

						// Create the player
						firstname := firstnameEntry.Text
						lastname := lastnameEntry.Text

						p, errName := mf.NewPlayer(firstname, lastname, db)
						if errName != nil {
							dialog.ShowError(errName, w)
							return
						} else {
							p.SetPlayerAge(age)
							p.SetPlayerRanking(ranking)
							p.SetPlayerMaterial(forehandEntry.Text, backhandEntry.Text, bladeEntry.Text)
						}

						// Link the player to the club
						err := mf.AddPlayerToClub(p, club)
						if err != nil {
							dialog.ShowError(err, w)
							return
						} else {
							// Player creation + link to club success
							successMsg := fmt.Sprintf(T("player_has_been_successfully_created"), firstname, lastname)
							fmt.Println(successMsg)
							dialog.ShowInformation(T("success"), successMsg, w)

							// Set the flag to true to indicate that the database has changed
							HasChanged = true

							// Reinit the entry texts
							reinitWidgetEntryText(firstnameEntry, entryFirstnameHolder)
							reinitWidgetEntryText(lastnameEntry, entryLastnameHolder)
							reinitWidgetEntryText(ageEntry, entryAgeHolder)
							reinitWidgetEntryText(rankingEntry, entryRankingHolder)
							reinitWidgetEntryText(forehandEntry, entryForehandHolder)
							reinitWidgetEntryText(backhandEntry, entryBackhandHolder)
							reinitWidgetEntryText(bladeEntry, entryBladeHolder)
						}

					})
					// Create a player in this club page
					pageTitle := setTitle(T("create_a_new_player_add_info"), 32)

					w.SetContent(container.NewVBox(
						pageTitle,
						clubLabel,
						firstnameEntry,
						lastnameEntry,
						ageEntry,
						rankingEntry,
						container.NewGridWithColumns(3, forehandEntry, backhandEntry, bladeEntry),
						validatationButton,
						ReturnToCreatePageButton,
					))
				})
				listOfClubs = append(listOfClubs, clubButton)
			}
			// Choose a club page
			w.SetContent(container.NewVBox(
				pageTitle,
				container.NewVBox(listOfClubs...),
			))
		})
		// Club selection page
		w.SetContent(container.NewVBox(
			pageTitle,
			clubSelectionPageButton,
			ReturnToCreatePageButton,
		))

	})

	// Team
	teamButton := widget.NewButton(T("create_a_new_team"), func() {

		// Club Selection
		pageTitle := setTitle(T("create_a_new_team_select_a_club"), 32)

		// clubSelectionPage
		clubSelectionPageButton := widget.NewButton(T("select_a_club"), func() {
			pageTitle := setTitle(T("create_a_new_team_select_a_club"), 32)
			listOfClubs := []fyne.CanvasObject{}

			for _, club := range db.Clubs {
				clubButton := widget.NewButton(club.Name, func() {
					// After club selection
					clubLabel := widget.NewLabel(fmt.Sprintf(T("you_are_going_to_create_a_team_for"), club.Name))

					// We can now create the team
					nameEntry := widget.NewEntry()
					entryHolder := T("enter_your_team_name_here")
					nameEntry.SetPlaceHolder(entryHolder)

					validatationButton := widget.NewButton(T("confirm"), func() {
						name := nameEntry.Text

						// If team name already exists, do not create the team
						for _, value := range db.Teams {
							if value.Name == name {
								err := fmt.Errorf(T("team_already_exists_in"), name, club.Name)
								dialog.ShowError(err, w)
								// Reinit the text
								nameEntry.SetText("")
								nameEntry.SetPlaceHolder(entryHolder)
								return
							}
						}

						t, err := mf.NewTeam(name, db)

						if err != nil {
							dialog.ShowError(err, w)
							return
						} else {
							// Link the team to the club
							err := mf.AddTeamToClub(t, club)
							if err != nil {
								dialog.ShowError(err, w)
								return
							} else {
								// team creation + link to club success
								successMsg := fmt.Sprintf(T("team_has_been_successfully_created"), name)
								fmt.Println(successMsg)
								dialog.ShowInformation(T("success"), successMsg, w)

								// Set the flag to true to indicate that the database has changed
								HasChanged = true

								// Reinit the text
								nameEntry.SetText("")
								nameEntry.SetPlaceHolder(entryHolder)
							}
						}
					})
					// Create a team in this club page
					pageTitle := setTitle(T("create_new_team"), 32)
					w.SetContent(container.NewVBox(
						pageTitle,
						clubLabel,
						nameEntry,
						validatationButton,
						ReturnToCreatePageButton,
					))
				})
				listOfClubs = append(listOfClubs, clubButton)
			}
			// Choose a club page
			w.SetContent(container.NewVBox(
				pageTitle,
				container.NewVBox(listOfClubs...),
			))
		})
		// Club selection page
		w.SetContent(container.NewVBox(
			pageTitle,
			clubSelectionPageButton,
			ReturnToCreatePageButton,
		))

	})
	// Club
	clubButton := widget.NewButton(T("create_a_new_club"), func() {
		nameEntry := widget.NewEntry()
		entryHolder := T("enter_your_club_name_here")
		nameEntry.SetPlaceHolder(entryHolder)

		validatationButton := widget.NewButton(T("confirm"), func() {
			name := nameEntry.Text
			_, err := mf.NewClub(name, db)

			if err != nil {
				dialog.ShowError(err, w)
				return
			} else {
				successMsg := fmt.Sprintf(T("club_has_been_successfully_created"), name)
				fmt.Println(successMsg)
				dialog.ShowInformation(T("success"), successMsg, w)

				// Set the flag to true to indicate that the database has changed
				HasChanged = true

				// Reinit the text
				nameEntry.SetText("")
				nameEntry.SetPlaceHolder(entryHolder)
			}
		})
		// Create a club page
		pageTitle := setTitle(T("create_new_club"), 32)
		w.SetContent(container.NewVBox(
			pageTitle,
			nameEntry,
			validatationButton,
			ReturnToCreatePageButton,
		))
	})

	// Page definition

	// If there is no club, a club must first be created

	if len(db.Clubs) < 1 {
		label := widget.NewLabel(T("you_currently_have_0_club_available_please_create"))

		createPage := container.NewVBox(
			pageTitle,
			label,
			clubButton,
			ReturnToFonctionalityPageButton,
		)
		// Create menu page
		w.SetContent(createPage)
	} else {
		createPage := container.NewVBox(
			pageTitle,
			playerButton,
			teamButton,
			clubButton,
			ReturnToFonctionalityPageButton,
		)
		// Create menu page
		w.SetContent(createPage)
	}
}
