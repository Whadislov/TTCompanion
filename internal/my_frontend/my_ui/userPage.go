package myapp

import (
	"errors"

	mdb "github.com/Whadislov/TTCompanion/internal/my_db"
	mr "github.com/Whadislov/TTCompanion/internal/my_frontend/my_requests"
	mf "github.com/Whadislov/TTCompanion/internal/my_functions"
	mt "github.com/Whadislov/TTCompanion/internal/my_types"
	"golang.org/x/crypto/bcrypt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// UserPage sets up the page for displaying user info.
func UserPage(user *mt.User, db *mt.Database, w fyne.Window, a fyne.App) {
	//pageTitle := setTitle("Your informations", 32)
	pageTitle := setTitle(T("your_informations"), 32)
	if HasUserProfileChanged {
		pageTitle = setTitle(T("your_informations_in_edition"), 32)
	} else {
		// Display the correct user info after modifications even before saving
		currentUsername = user.Name
		currentEmail = user.Email
	}

	usernameLabel1 := widget.NewLabel(T("username") + ":")
	usernameLabel2 := widget.NewLabel(currentUsername)
	usernameContent := container.NewGridWithColumns(2, usernameLabel1, usernameLabel2)

	emailLabel1 := widget.NewLabel("✉️ E-mail" + ":")
	emailLabel2 := widget.NewLabel(currentEmail)
	emailContent := container.NewGridWithColumns(2, emailLabel1, emailLabel2)

	passwordLabel := widget.NewLabel(T("password") + ":")

	createdAtLabel1 := widget.NewLabel(T("user_since") + ":")
	// Example format: 2025-02-03T09:58:59Z
	// [:10] -> only show the date and not the hours
	readableCreatedAt := user.CreatedAt[:10]
	createdAtLabel2 := widget.NewLabel(readableCreatedAt)
	createdAtContent := container.NewGridWithColumns(2, createdAtLabel1, createdAtLabel2)

	changeUsernameButton := widget.NewButton(T("edit_username"), func() {
		ChangeUsernamePage(user, db, w, a)
	})

	changeEmailButton := widget.NewButton(T("edit_email"), func() {
		ChangeEmailPage(user, db, w, a)
	})

	changePasswordButton := widget.NewButton(T("edit_password"), func() {
		ChangePasswordPage(user, db, w, a)
	})

	passwordContent := container.NewGridWithColumns(2, passwordLabel)

	saveChangesButton := widget.NewButton(T("save_changes"), func() {
		if !HasUserProfileChanged {
			dialog.ShowInformation(T("information"), T("there_is_nothing_new_to_save"), w)
			return
		} else {
			var err error
			if appStartOption == "local" {
				err = mdb.SaveDB(db)
			} else if appStartOption == "browser" {
				stopChan := make(chan string, 1)
				go func() {
					loadingWindow(T("saving"), w, stopChan)
				}()
				err = mr.SaveDB(credToken, db)
				stopChan <- "ok"
			}

			if err != nil {
				dialog.ShowError(err, w)
			} else {
				dialog.ShowInformation(T("success"), T("changes_saved"), w)
				HasUserProfileChanged = false
			}
			UserPage(user, db, w, a)
		}
	})

	returnToMainPageButton := widget.NewButton(T("return_to_main_page"), func() {
		if HasUserProfileChanged {
			dialog.ShowConfirm(T("unsaved_changes"), T("you_have_unsaved_changes_alt"), func(confirm bool) {
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
				}
				//Reset flag in both case
				HasUserProfileChanged = false
				//Set Main page in both case
				w.SetContent(MainPage(db, w, a))
			}, w)
		} else {
			w.SetContent(MainPage(db, w, a))
		}
	})

	content := container.NewVBox(
		pageTitle,
		container.NewGridWithColumns(2, usernameContent, changeUsernameButton),
		container.NewGridWithColumns(2, emailContent, changeEmailButton),
		container.NewGridWithColumns(2, passwordContent, changePasswordButton),
		createdAtContent,
		saveChangesButton,
		returnToMainPageButton,
	)

	w.SetContent(content)
}

// ChangeUsernamePage sets up the page to change user's username.
func ChangeUsernamePage(user *mt.User, db *mt.Database, w fyne.Window, a fyne.App) {

	pageTitle := setTitle(T("editing_username"), 32)

	usernameLabel := widget.NewLabel(T("current_username") + user.Name)
	editUsernameEntry := widget.NewEntry()
	editUsernameEntry.SetPlaceHolder(T("enter_new_username"))

	confirmButton := widget.NewButton(T("confirm"), func() {
		err := mf.ChangeUsername(user.Name, editUsernameEntry.Text, db)
		if err != nil {
			dialog.ShowError(err, w)
			editUsernameEntry.SetPlaceHolder(T("enter_new_username"))
		} else {
			dialog.ShowInformation(T("success"), T("username_successfully_changed"), w)
			HasUserProfileChanged = true
			currentUsername = editUsernameEntry.Text
			UserPage(user, db, w, a)
		}

	})

	cancelButton := widget.NewButton(T("cancel"), func() {
		UserPage(user, db, w, a)
	})

	content := container.NewVBox(
		pageTitle,
		usernameLabel,
		editUsernameEntry,
		confirmButton,
		cancelButton,
	)
	w.SetContent(content)
}

// ChangeEmailPage sets up the page to change user's email.
func ChangeEmailPage(user *mt.User, db *mt.Database, w fyne.Window, a fyne.App) {

	pageTitle := setTitle(T("editing_email"), 32)

	emailLabel := widget.NewLabel(T("current_email") + user.Email)
	editEmailEntry := widget.NewEntry()
	editEmailEntry.SetPlaceHolder(T("enter_new_email"))

	confirmButton := widget.NewButton(T("confirm"), func() {
		err := mf.ChangeEmail(editEmailEntry.Text, user)
		if err != nil {
			dialog.ShowError(err, w)
			editEmailEntry.SetPlaceHolder(T("enter_new_email"))
		} else {
			dialog.ShowInformation(T("success"), T("email_successfully_changed"), w)
			HasUserProfileChanged = true
			currentEmail = editEmailEntry.Text
			UserPage(user, db, w, a)
		}

	})

	cancelButton := widget.NewButton(T("cancel"), func() {
		UserPage(user, db, w, a)
	})

	content := container.NewVBox(
		pageTitle,
		emailLabel,
		editEmailEntry,
		confirmButton,
		cancelButton,
	)
	w.SetContent(content)
}

// ChangePasswordPage sets up the page to change user's email.
func ChangePasswordPage(user *mt.User, db *mt.Database, w fyne.Window, a fyne.App) {

	pageTitle := setTitle(T("editing_password"), 32)

	passwordLabel := widget.NewLabel(T("enter_current_password"))

	currentPasswordEntry := widget.NewPasswordEntry()
	editPasswordEntry := widget.NewPasswordEntry()
	confirmEditPasswordEntry := widget.NewPasswordEntry()
	editPasswordLabel := widget.NewLabel(T("enter_new_password"))
	confirmEditPasswordLabel := widget.NewLabel(T("confirm_your_new_password"))

	confirmButton := widget.NewButton(T("confirm"), func() {
		err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(currentPasswordEntry.Text))
		if err != nil {
			// Password is wrong
			dialog.ShowError(errors.New(T("wrong_password")), w)
			currentPasswordEntry.SetText("")
			editPasswordEntry.SetText("")
			confirmEditPasswordEntry.SetText("")
		} else if editPasswordEntry.Text != confirmEditPasswordEntry.Text {
			dialog.ShowError(errors.New(T("new_passwords_do_not_match")), w)
			currentPasswordEntry.SetText("")
			editPasswordEntry.SetText("")
			confirmEditPasswordEntry.SetText("")
		} else {
			err := mf.ChangePassword(editPasswordEntry.Text, user)
			if err != nil {
				dialog.ShowError(err, w)
				currentPasswordEntry.SetText("")
				editPasswordEntry.SetText("")
				confirmEditPasswordEntry.SetText("")
			} else {
				dialog.ShowInformation(T("success"), T("password_successfully_changed"), w)
				HasUserProfileChanged = true
				UserPage(user, db, w, a)
			}
		}
	})

	cancelButton := widget.NewButton(T("cancel"), func() {
		UserPage(user, db, w, a)
	})

	content := container.NewVBox(
		pageTitle,
		passwordLabel,
		currentPasswordEntry,
		editPasswordLabel,
		editPasswordEntry,
		confirmEditPasswordLabel,
		confirmEditPasswordEntry,
		confirmButton,
		cancelButton,
	)
	w.SetContent(content)

}
