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

// MainPage creates the main page
func MainPage(db *mt.Database, w fyne.Window, a fyne.App) {

	var mainPage *fyne.Container
	pageTitle := setTitle("TT Companion", 32)

	showDBButton := widget.NewButton(T("your_database"), func() {
		w.SetContent(DatabasePage(db, w, a))
	})

	showFuncButton := widget.NewButton(T("functionalities"), func() {
		w.SetContent(FunctionalityPage(db, w, a))
	})

	// Options button
	OptionButton := widget.NewButton(T("options"), func() {
		returnToMainMenuButton := widget.NewButton(T("return_to_main_page"), func() {
			MainPage(db, w, a)
		})
		w.SetContent(container.NewVBox(OptionPage(db, w, a), returnToMainMenuButton))
	})

	SaveButton := widget.NewButton(T("save"), func() {
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

	disconnectButton := widget.NewButton(T("log_out"), func() {
		if HasChanged {
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

	quitButton := widget.NewButton(T("quit"), func() {
		Quit(db, w, a, HasChanged)
	})

	if appStartOption == "local" {
		mainPage = container.NewVBox(
			pageTitle,
			showDBButton,
			showFuncButton,
			OptionButton,
			SaveButton,
			disconnectButton,
			quitButton,
		)

	} else if appStartOption == "browser" {
		// Remove the quit button
		mainPage = container.NewVBox(
			pageTitle,
			showDBButton,
			showFuncButton,
			OptionButton,
			SaveButton,
			disconnectButton,
		)
	}

	// Check for unsaved changes before quitting
	w.SetCloseIntercept(func() {
		Quit(db, w, a, HasChanged)
	})

	w.SetContent(mainPage)

}
