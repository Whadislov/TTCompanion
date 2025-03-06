package myapp

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	mdb "github.com/Whadislov/TTCompanion/internal/my_db"
	mr "github.com/Whadislov/TTCompanion/internal/my_frontend/my_requests"
	mt "github.com/Whadislov/TTCompanion/internal/my_types"
)

// MainMenu creates the menus
func MainMenu(db *mt.Database, w fyne.Window, a fyne.App) *fyne.MainMenu {

	mainPage := MainPage(db, w, a)

	menu1Item1 := fyne.NewMenuItem(T("main_page"), func() { w.SetContent(mainPage) })
	menu1Item2 := fyne.NewMenuItem(T("my_profile"), func() { UserPage(userOfSession, db, w, a) })
	menu1Item3 := fyne.NewMenuItem(T("options"), func() {
		returnToMainMenuButton := widget.NewButton(T("return_to_main_page"), func() {
			w.SetContent(MainPage(db, w, a))
		})
		w.SetContent(container.NewVBox(OptionPage(db, w, a), returnToMainMenuButton))
	})
	menu1Item4 := fyne.NewMenuItem(T("save_changes"), func() {
		if !HasChanged {
			dialog.ShowInformation(T("information"), T("there_is_nothing_new_to_save"), w)
		} else {
			if appStartOption == "local" {
				err := mdb.SaveDB(db)
				if err != nil {
					dialog.ShowError(err, w)
				} else {
					HasChanged = false
					// Reload the database after saving (refresh the IDs)
					db, err = mdb.LoadDB()
					if err != nil {
						dialog.ShowError(err, w)
					} else {
						dialog.ShowInformation(T("information"), T("changes_saved"), w)
					}
				}

			} else if appStartOption == "browser" {
				stopChan := make(chan string, 1)
				go func() {
					loadingWindow(T("saving"), w, stopChan)
				}()
				err := mr.SaveDB(credToken, db)
				if err != nil {
					stopChan <- "ok"
					dialog.ShowError(err, w)
				} else {
					HasChanged = false
					// Reload the database after saving (refresh the IDs)
					db, err = mr.LoadDB(credToken)
					stopChan <- "ok"
					if err != nil {
						dialog.ShowError(err, w)
					} else {
						dialog.ShowInformation(T("information"), T("changes_saved"), w)
					}
				}
			}
		}
	})
	menu1Item5 := fyne.NewMenuItem(T("log_out"), func() {
		if HasChanged {
			dialog.ShowConfirm(T("unsaved_changes"), T("you_have_unsaved_changes"), func(confirm bool) {
				if confirm {
					// User wants to save the changes
					if appStartOption == "local" {
						err := mdb.SaveDB(db)
						if err != nil {
							dialog.ShowError(err, w)
						} else {
							HasChanged = false
							log.Println("User logged out and saved his changes.")
							w.SetMainMenu(nil)
							w.SetContent(AuthentificationPage(w, a))
						}
					} else if appStartOption == "browser" {
						stopChan := make(chan string, 1)
						go func() {
							loadingWindow(T("saving"), w, stopChan)
						}()
						err := mr.SaveDB(credToken, db)
						stopChan <- "ok"
						if err != nil {
							dialog.ShowError(err, w)
						} else {
							log.Println("User logged out and saved his changes.")
							// Reset the token, the flag, the menu and the database
							credToken = ""
							HasChanged = false
							w.SetMainMenu(nil)
							db = &mt.Database{}
							w.SetContent(AuthentificationPageWeb(w, a))
						}
					}
				} else {
					// User does not want to save the changes
					log.Println("User logged out and saved nothing.")
					if appStartOption == "local" {
						// Reset the flag and the menu
						HasChanged = false
						w.SetMainMenu(nil)
						w.SetContent(AuthentificationPage(w, a))
					} else if appStartOption == "browser" {
						// Reset the token, the flag, the menu and the database
						credToken = ""
						HasChanged = false
						w.SetMainMenu(nil)
						db = &mt.Database{}
						w.SetContent(AuthentificationPageWeb(w, a))
					}
				}
			}, w)
		} else {
			// Nothing has changed
			log.Println("User logged out.")
			if appStartOption == "local" {
				// Reset the menu
				w.SetMainMenu(nil)
				w.SetContent(AuthentificationPage(w, a))
			} else if appStartOption == "browser" {
				// Reset the menu, the token and the database
				credToken = ""
				w.SetMainMenu(nil)
				db = &mt.Database{}
				w.SetContent(AuthentificationPageWeb(w, a))
			}
		}

	})

	// menu item names are not linked with pageTitles, if there is a modification here -> modify also on pages
	newMenu1 := fyne.NewMenu(T("main_menu"), menu1Item1, menu1Item2, menu1Item3, menu1Item4, menu1Item5)

	menu2Item1 := fyne.NewMenuItem(T("players"), func() { PlayerPage(db, w, a) })
	menu2Item2 := fyne.NewMenuItem(T("teams"), func() { TeamPage(db, w, a) })
	menu2Item3 := fyne.NewMenuItem(T("clubs"), func() { ClubPage(db, w, a) })
	newMenu2 := fyne.NewMenu(T("database"), menu2Item1, menu2Item2, menu2Item3)

	menu3Item1 := fyne.NewMenuItem(T("create"), func() { CreatePage(db, w, a) })
	menu3Item2 := fyne.NewMenuItem(T("add"), func() { AddPage(db, w, a) })
	menu3Item3 := fyne.NewMenuItem(T("remove"), func() { RemovePage(db, w, a) })
	menu3Item4 := fyne.NewMenuItem(T("delete"), func() { DeletePage(db, w, a) })
	menu3Item5 := fyne.NewMenuItem(T("edit_player_information"), func() { AddInfoToPlayerPage(db, w, a) })
	newMenu3 := fyne.NewMenu(T("functionalities"), menu3Item1, menu3Item2, menu3Item3, menu3Item4, menu3Item5)

	menu := fyne.NewMainMenu(newMenu1, newMenu2, newMenu3)

	return menu

}
