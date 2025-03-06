package myapp

import (
	"errors"
	"fmt"
	"log"

	mr "github.com/Whadislov/TTCompanion/internal/my_frontend/my_requests"
	mf "github.com/Whadislov/TTCompanion/internal/my_functions"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// AuthentificationPageWeb returns a page that contains a log in button and a sign up button
func AuthentificationPageWeb(w fyne.Window, a fyne.App) *fyne.Container {

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
		w.SetContent(loginPageWeb(w, a))
	})

	signUpButton := widget.NewButton(T("sign_up"), func() {
		w.SetContent(signUpPageWeb(w, a))
	})

	returnToAuthentificationPageWebButton := widget.NewButton(T("return_to_authentification_page"), func() {
		w.SetContent(AuthentificationPageWeb(w, a))
	})

	optionPageButton := widget.NewButton(T("options"), func() {
		w.SetContent(container.NewVBox(
			OptionAuthPageWeb(w, a),
			returnToAuthentificationPageWebButton,
		))
	})

	demoButton := widget.NewButton(T("demo"), func() {

		db, token, err := mr.Login("demo", "demo")
		if err != nil {
			dialog.ShowError(errors.New(T("internal_error")), w)
		}

		credToken = token
		log.Println("Demo account logged in")
		for _, user := range db.Users {
			if user.Name == "demo" {
				userOfSession = user
			} else {
				dialog.ShowError(errors.New(T("failed_to_load_profile")), w)
			}
		}
		w.SetContent(MainPage(db, w, a))
		w.SetMainMenu(MainMenu(db, w, a))
	})

	authentificationPageWeb := container.NewVBox(
		pageTitle,
		authLabel,
		logInButton,
		signUpButton,
		optionPageButton,
		demoButton,
	)

	w.SetContent(authentificationPageWeb)
	return authentificationPageWeb

}

// signUpPageWeb returns a page that contains a create username and a create password field. Adds a new user in the database
func signUpPageWeb(w fyne.Window, a fyne.App) *fyne.Container {
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

	validationButton := widget.NewButton(T("confirm"), func() {
		log.Println("Creating new User")
		_, err := mf.IsValidEmail(emailEntry.Text)
		_, err2 := mf.IsValidUsername(usernameEntry.Text)
		_, err3 := mf.IsValidPassword(passwordEntry.Text)

		if err != nil {
			dialog.ShowError(errors.New(T("email.must_be_valid")), w)
			log.Println("e-mail is not valid:", err)
			w.SetContent(signUpPageWeb(w, a))
		} else if err2 != nil {
			dialog.ShowError(errors.New(T("username.must_be_valid")), w)
			log.Println("username is not valid:", err)
			w.SetContent(signUpPageWeb(w, a))
		} else if err3 != nil {
			dialog.ShowError(errors.New(T("password.must_be_valid")), w)
			log.Println("password is not valid:", err)
			w.SetContent(signUpPageWeb(w, a))
		} else if passwordEntry.Text != confirmPasswordEntry.Text {
			dialog.ShowError(errors.New(T("password.not_matching")), w)
			log.Println("passwords do not match:")
			w.SetContent(signUpPageWeb(w, a))
		} else {
			db, token, err := mr.SignUp(usernameEntry.Text, passwordEntry.Text, emailEntry.Text)
			credToken = token
			// Last check if username or email already exist
			if err != nil {
				// Need to recheck this err
				log.Println("Username or e-mail already exists")
				dialog.ShowError(errors.New(T("failed_to_sign_up")), w)
				w.SetContent(signUpPageWeb(w, a))
			} else {
				log.Println("Sign up is successfull")
				// Set the user of the session for the profile page on the menu. There is only one user on the Users map
				for _, user := range db.Users {
					if user.Name == usernameEntry.Text {
						userOfSession = user
					} else {
						dialog.ShowError(errors.New(T("failed_to_load_profile")), w)
					}
				}
				w.SetContent(MainPage(db, w, a))
				w.SetMainMenu(MainMenu(db, w, a))
			}
		}
	})

	cancelButton := widget.NewButton(T("cancel"), func() {
		w.SetContent(AuthentificationPageWeb(w, a))
	})

	signUpPageWeb := container.NewVBox(
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

	return signUpPageWeb
}

// logInPageWeb returns a page that contains a enter username and a enter password field. Checks if the user is in the database. If yes, sets the variable user_id for the rest of the program
func loginPageWeb(w fyne.Window, a fyne.App) *fyne.Container {
	pageTitle := setTitle(T("login"), 32)

	usernameLabel := widget.NewLabel(T("username"))
	usernameEntry := widget.NewEntry()

	passwordLabel := widget.NewLabel(T("password"))
	passwordEntry := widget.NewPasswordEntry()

	validationButton := widget.NewButton(T("login"), func() {
		stopChan := make(chan string)
		go func() {
			loadingWindow(w, stopChan)
		}()
		db, token, err := mr.Login(usernameEntry.Text, passwordEntry.Text)
		stopChan <- "ok"
		credToken = token
		if err != nil {
			err_username_password_missmatch := errors.New("username or password is invalid")
			if err_username_password_missmatch == err {
				dialog.ShowError(fmt.Errorf(T("username_and_password_missmatch"), err), w)
				w.SetContent(loginPageWeb(w, a))
				return
			} else {
				dialog.ShowError(fmt.Errorf(T("failed_to_log_in"), " %v", err), w)
				w.SetContent(loginPageWeb(w, a))
				return
			}
		} else {
			log.Println("Login is successfull")
			// Set the user of the session for the profile page on the menu. There is only one user on the Users map
			for _, user := range db.Users {
				if user.Name == usernameEntry.Text {
					userOfSession = user
				} else {
					dialog.ShowError(errors.New(T("failed_to_load_profile")), w)
				}
			}
			w.SetContent(MainPage(db, w, a))
			w.SetMainMenu(MainMenu(db, w, a))
		}
	})

	forgotPasswordButton := widget.NewButtonWithIcon(T("forgot_password"), nil, func() {
		w.SetContent(reinitPasswordPageWeb(w, a))
	})
	forgotPasswordButton.Importance = widget.LowImportance

	cancelButton := widget.NewButton(T("cancel"), func() {
		w.SetContent(AuthentificationPageWeb(w, a))
	})

	loginPageWeb := container.NewVBox(
		pageTitle,
		usernameLabel,
		usernameEntry,
		container.NewBorder(nil, nil, passwordLabel, forgotPasswordButton),
		passwordEntry,
		validationButton,
		cancelButton,
	)

	return loginPageWeb
}

// reinitPasswordPage offers the user the ability to be sent an Email for a password reset
func reinitPasswordPageWeb(w fyne.Window, a fyne.App) *fyne.Container {
	pageTitle := setTitle(T("forgot_password"), 32)
	pageText := sexText(T("text.forgot_password"), 12)
	emailString := T("enter_email_address")

	emailEntry := widget.NewEntry()
	emailEntry.SetPlaceHolder(emailString)

	validationButton := widget.NewButton(T("reset_your_password"), func() {
		b, err := mf.IsValidEmail(emailEntry.Text)
		if !b {
			dialog.ShowError(err, w)
			w.SetContent(reinitPasswordPageWeb(w, a))
		} else {

			// E-mail found
			var resetEmailDialogSuccess *dialog.CustomDialog
			cautionMessage := T("caution_message")
			resetEmailMessage := widget.NewLabel(fmt.Sprintf(cautionMessage+T("message.email_sent"), emailEntry.Text))
			returnToLoginPageWeb := widget.NewButton(T("return_to_login_screen"), func() {
				resetEmailDialogSuccess.Hide()
				w.SetContent(loginPageWeb(w, a))
			})
			resetEmailContent := container.NewVBox(
				resetEmailMessage,
				returnToLoginPageWeb,
			)

			resetEmailDialogSuccess = dialog.NewCustomWithoutButtons(T("success"), resetEmailContent, w)
			resetEmailDialogSuccess.Show()
		}

	})

	cancelButton := widget.NewButton(T("cancel"), func() {
		w.SetContent(loginPageWeb(w, a))
	})

	reinitPasswordPageWeb := container.NewVBox(
		pageTitle,
		pageText,
		emailEntry,
		validationButton,
		cancelButton,
	)

	return reinitPasswordPageWeb
}
