package myapp

import (
	"fmt"
	"log"

	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	mdb "github.com/Whadislov/TTCompanion/internal/my_db"
	mf "github.com/Whadislov/TTCompanion/internal/my_functions"
	"golang.org/x/crypto/bcrypt"
)

// AuthentificationPage returns a page that contains a log in button and a sign up button
func AuthentificationPage(w fyne.Window, a fyne.App) *fyne.Container {

	// Know if light mode is activated or not
	if a.Settings().ThemeVariant() == 1 {
		lightTheme.IsActivated = true
	} else {
		darkTheme.IsActivated = true
	}

	pageTitle := setTitle("TT Companion", 32)

	authLabel := widget.NewLabel(T("please_authenticate"))
	authLabel.Alignment = fyne.TextAlignCenter

	logInButton := widget.NewButton(T("log_in"), func() {
		w.SetContent(loginPage(w, a))
	})

	signUpButton := widget.NewButton(T("sign_up"), func() {
		w.SetContent(signUpPage(w, a))
	})

	returnToAuthentificationPageButton := widget.NewButton(T("return_to_authentification_page"), func() {
		w.SetContent(AuthentificationPage(w, a))
	})

	optionPageButton := widget.NewButton(T("options"), func() {
		w.SetContent(container.NewVBox(
			OptionAuthPage(w, a),
			returnToAuthentificationPageButton,
		))
	})

	demoButton := widget.NewButton(T("demo"), func() {

		db, err := mdb.LoadUsersOnly()
		if err != nil {
			panic(err)
		}

		for _, user := range db.Users {
			if user.Name == "demo" {
				err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte("demo"))
				if err != nil {
					dialog.ShowError(fmt.Errorf(T("username_and_password_missmatch"), err), w)
					return
				} else {
					log.Println("Demo account logged in")
					mdb.SetUserIDOfSession(user.ID)
					// Now load the corresponding database of the user
					var err error

					mdb.SetUserIDOfSession(user.ID)
					db, err = mdb.LoadDB()

					if err != nil {
						dialog.ShowError(err, w)
					} else {
						userOfSession = db.Users[user.ID]
						MainPage(db, w, a)
						w.SetMainMenu(MainMenu(db, w, a))
						return
					}
				}
			}
		}
	})

	authentificationPage := container.NewVBox(
		pageTitle,
		authLabel,
		logInButton,
		signUpButton,
		optionPageButton,
		demoButton,
	)
	return authentificationPage

}

// signUpPage returns a page that contains a create username and a create password field. Adds a new user in the database
func signUpPage(w fyne.Window, a fyne.App) *fyne.Container {
	pageTitle := setTitle(T("create_your_account"), 32)

	emailLabel := widget.NewLabel("✉️ E-mail")
	emailEntry := widget.NewEntry()
	emailEntry.SetPlaceHolder("abc@def.com")

	usernameLabel := widget.NewLabel(T("username"))
	usernameEntry := widget.NewEntry()

	passwordLabel := widget.NewLabel(T("password"))
	passwordEntry := widget.NewPasswordEntry()

	confirmPasswordLabel := widget.NewLabel(T("confirm_password"))
	confirmPasswordEntry := widget.NewPasswordEntry()

	db, err := mdb.LoadUsersOnly()
	if err != nil {
		panic(err)
	}

	validationButton := widget.NewButton(T("confirm"), func() {

		log.Println("Creating new User")
		newUser, err := mf.NewUser(usernameEntry.Text, emailEntry.Text, passwordEntry.Text, confirmPasswordEntry.Text, db)
		if err != nil {
			switch err.Error() {
			case "e-mail cannot be empty":
				log.Println("e-mail is empty")
				dialog.ShowError(errors.New(T("email.cannot_be_empty")), w)
				emailEntry.SetPlaceHolder("abc@def.com")
			case "e-mail is already used":
				log.Println("e-mail is already used")
				dialog.ShowError(errors.New(T("email.already_used")), w)
				emailEntry.SetPlaceHolder("abc@def.com")
			case "e-mail must be valid. Example: abc@def.com":
				log.Println("e-mail is not valid")
				dialog.ShowError(errors.New(T("email.must_be_valid")), w)
				emailEntry.SetPlaceHolder("abc@def.com")
			case "username cannot be empty":
				log.Println("username is empty")
				dialog.ShowError(errors.New(T("username.must_be_valid")), w)
				usernameEntry.SetPlaceHolder("")
			case "username must be valid (only letters and figures are allowed, spaces are not allowed)":
				log.Println("username is not valid")
				dialog.ShowError(errors.New(T("dialog_error.invalid_username")), w)
				usernameEntry.SetPlaceHolder("")
			case "username is already taken":
				log.Println("username is already taken")
				dialog.ShowError(errors.New(T("username.already_used")), w)
				usernameEntry.SetPlaceHolder("")
			case "password cannot be empty":
				log.Println("password cannot be empty")
				dialog.ShowError(errors.New(T("password.cannot_be_empty")), w)
				passwordEntry.SetText("")
				confirmPasswordEntry.SetText("")
			case "password must be valid (spaces are not allowed)":
				log.Println("password is not valid")
				dialog.ShowError(errors.New(T("password.must_be_valid")), w)
				passwordEntry.SetText("")
				confirmPasswordEntry.SetText("")
			case "passwords do not match":
				log.Println("passwords do not match")
				dialog.ShowError(errors.New(T("password.not_matching")), w)
				passwordEntry.SetText("")
				confirmPasswordEntry.SetText("")
			default:
				log.Println("unknown error")
				dialog.ShowError(errors.New(T("unknown_error")), w)
				w.SetContent(signUpPage(w, a))
			}
		} else {
			var err error
			dialog.ShowInformation(T("success"), T("message.account_created"), w)
			mdb.SetUserIDOfSession(newUser.ID)

			// Save the new user in the database
			mdb.SaveDB(db)
			log.Println("Sign up is successfull")

			// Now load the whole database
			db, err = mdb.LoadDB()

			if err != nil {
				dialog.ShowError(err, w)
			} else {
				// Only one user is loaded
				for _, user := range db.Users {
					userOfSession = user
					mdb.SetUserIDOfSession(userOfSession.ID)
					MainPage(db, w, a)
					w.SetMainMenu(MainMenu(db, w, a))
					return
				}
			}
		}
	})

	cancelButton := widget.NewButton(T("cancel"), func() {
		w.SetContent(AuthentificationPage(w, a))
	})

	signUpPage := container.NewVBox(
		pageTitle,
		emailLabel,
		emailEntry,
		usernameLabel,
		usernameEntry,
		passwordLabel,
		passwordEntry,
		confirmPasswordLabel,
		confirmPasswordEntry,
		validationButton,
		cancelButton,
	)

	return signUpPage
}

// loginPage returns a page that contains a enter username and a enter password field. Checks if the user is in the database. If yes, sets the variable user_id for the rest of the program
func loginPage(w fyne.Window, a fyne.App) *fyne.Container {
	pageTitle := setTitle(T("login"), 32)

	usernameLabel := widget.NewLabel(T("username"))
	usernameEntry := widget.NewEntry()

	passwordLabel := widget.NewLabel(T("password"))
	passwordEntry := widget.NewPasswordEntry()

	db, err := mdb.LoadUsersOnly()
	if err != nil {
		panic(err)
	}

	validationButton := widget.NewButton(T("login"), func() {

		// Verify if username and password match
		log.Println("Verifying username and password")
		for _, user := range db.Users {
			if usernameEntry.Text == user.Name {
				err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(passwordEntry.Text))
				if err != nil {
					dialog.ShowError(fmt.Errorf(T("username_and_password_missmatch")), w)
					return
				} else {
					log.Println("Login is successfull")
					mdb.SetUserIDOfSession(user.ID)
					// Now load the corresponding database of the user
					var err error

					mdb.SetUserIDOfSession(user.ID)
					db, err = mdb.LoadDB()

					if err != nil {
						dialog.ShowError(err, w)
					} else {
						userOfSession = db.Users[user.ID]
						MainPage(db, w, a)
						w.SetMainMenu(MainMenu(db, w, a))
						return
					}
				}
			}
		}
		dialog.ShowError(errors.New(T("username_and_password_missmatch")), w)
		// Reset entries
		reinitWidgetEntryText(usernameEntry, "")
		passwordEntry.SetText("")

	})

	forgotPasswordButton := widget.NewButtonWithIcon(T("forgot_password"), nil, func() {
		w.SetContent(reinitPasswordPage(w, a))
	})
	forgotPasswordButton.Importance = widget.LowImportance

	cancelButton := widget.NewButton(T("cancel"), func() {
		w.SetContent(AuthentificationPage(w, a))
	})

	loginPage := container.NewVBox(
		pageTitle,
		usernameLabel,
		usernameEntry,
		container.NewBorder(nil, nil, passwordLabel, forgotPasswordButton),
		passwordEntry,
		validationButton,
		cancelButton,
	)

	return loginPage
}

// reinitPasswordPage offers the user the ability to be sent an Email for a password reset
func reinitPasswordPage(w fyne.Window, a fyne.App) *fyne.Container {
	pageTitle := setTitle(T("forgot_password"), 32)
	pageText := sexText(T("text.forgot_password"), 12)
	emailString := T("enter_email_address")

	emailEntry := widget.NewEntry()
	emailEntry.SetPlaceHolder(emailString)

	db, err := mdb.LoadUsersOnly()
	if err != nil {
		panic(err)
	}

	validationButton := widget.NewButton(T("reset_your_password"), func() {
		b, err := mf.IsValidEmail(emailEntry.Text)
		if !b {
			dialog.ShowError(err, w)
			w.SetContent(reinitPasswordPage(w, a))
		} else {
			for _, user := range db.Users {
				if emailEntry.Text == user.Email {
					// E-mail found
					var resetEmailDialogSuccess *dialog.CustomDialog
					resetEmailMessage := widget.NewLabel(fmt.Sprintf(T("message.email_not_found"), emailEntry.Text))
					//resetEmailMessage := widget.NewLabel(fmt.Sprintf("A message has been sent to %v.\nIf you can't find it in your inbox, check your junk mail.\nOtherwise, check that you have correctly entered the e-mail address you used to log in.", emailEntry.Text))
					returnToLoginPage := widget.NewButton(T("return_to_login_screen"), func() {
						resetEmailDialogSuccess.Hide()
						w.SetContent(loginPage(w, a))
					})
					resetEmailContent := container.NewVBox(
						resetEmailMessage,
						returnToLoginPage,
					)

					resetEmailDialogSuccess = dialog.NewCustomWithoutButtons(T("success"), resetEmailContent, w)
					resetEmailDialogSuccess.Show()
				}
			}
			// E-mail valid, but not found
			var resetEmailDialogFail *dialog.CustomDialog
			goToSignupPage := widget.NewButton(T("create_your_account"), func() {
				resetEmailDialogFail.Hide()
				w.SetContent(signUpPage(w, a))
			})
			returnToLoginPage := widget.NewButton(T("try_again"), func() {
				resetEmailDialogFail.Hide()
				w.SetContent(reinitPasswordPage(w, a))
			})
			resetEmailMessageFail := widget.NewLabel(fmt.Sprintf(T("message.email_sent"), emailEntry.Text))
			resetEmailContentFail := container.NewVBox(
				resetEmailMessageFail,
				returnToLoginPage,
				goToSignupPage,
			)

			resetEmailDialogFail = dialog.NewCustomWithoutButtons(T("attention"), resetEmailContentFail, w)
			resetEmailDialogFail.Show()
		}

	})

	cancelButton := widget.NewButton(T("cancel"), func() {
		w.SetContent(loginPage(w, a))
	})

	reinitPasswordPage := container.NewVBox(
		pageTitle,
		pageText,
		emailEntry,
		validationButton,
		cancelButton,
	)

	return reinitPasswordPage
}
