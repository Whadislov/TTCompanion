package myapp

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"

	mdb "github.com/Whadislov/TTCompanion/internal/my_db"
	mr "github.com/Whadislov/TTCompanion/internal/my_frontend/my_requests"
	mt "github.com/Whadislov/TTCompanion/internal/my_types"
)

func Quit(db *mt.Database, w fyne.Window, a fyne.App, HasChanged bool) {

	// Check for unsaved changes before quitting

	if HasChanged {
		dialog.ShowConfirm(T("unsaved_changes"), T("save_before_quit"), func(confirm bool) {
			if confirm {
				if appStartOption == "local" {
					err := mdb.SaveDB(db)
					if err != nil {
						dialog.ShowError(err, w)
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
					}
				}
				a.Quit()
				log.Println("Application exited.")

			} else {
				a.Quit()
				log.Println("Application exited.")

			}
		}, w)
	} else {
		a.Quit()
		log.Println("Application exited.")
	}
}
