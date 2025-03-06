package myapp

import (
	"bytes"
	"errors"
	"fmt"
	"image/gif"
	"regexp"
	"sort"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	mr "github.com/Whadislov/TTCompanion/internal/my_frontend/my_requests"
	mt "github.com/Whadislov/TTCompanion/internal/my_types"
	"github.com/google/uuid"
)

// strHelper is a helper fonction that takes from example ["ok1", "ok2" , "ok3"] and returns "ok1, ok2, ok3"
func strHelper(list []string) string {
	str := ""
	for _, word := range list {
		str += word + ", "
	}
	// Remove extra ", "
	if len(str) > 2 {
		str = str[:len(str)-2]
	}
	return str
}

func sortMap[T ~map[uuid.UUID]V, V mt.Entity](m T) []struct {
	Key   uuid.UUID
	Value V
} {
	// This function takes as parameter a variable of type : mt.Player or mt.Team or mt.Club
	// It returns a struct that contains the keys and the name of player/team/club.
	// The slices are alphabeticaly sorted
	// Tip: it does not return a map, because maps can't be sorted

	// Get keys from the map
	keys := make([]uuid.UUID, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}

	// Sort keys alphabeticaly per value using m.Lastname. Keys is the sorted slice
	sort.Slice(keys, func(i, j int) bool {
		return m[keys[i]].GetName() < m[keys[j]].GetName()
	})

	// Build a new slice with the pair key/value
	sorted := make([]struct {
		Key   uuid.UUID
		Value V
	}, len(keys))
	// The new slice is sorted thanks to the slice keys
	for i, k := range keys {
		sorted[i] = struct {
			Key   uuid.UUID
			Value V
		}{
			Key:   k,
			Value: m[k],
		}
	}

	return sorted
}

// Create a confirmation dialog
func showConfirmationDialog(w fyne.Window, message string, onConfirm func()) {
	d := dialog.NewConfirm(T("confirm"), message, func(confirm bool) {
		if confirm {
			onConfirm()
		}
	}, w)
	d.Show()
}

// Reinit the text of a widget entry
func reinitWidgetEntryText(entry *widget.Entry, entryHolder string) {
	entry.SetText("")
	entry.SetPlaceHolder(entryHolder)
}

// isValidName verifies that the name follows some criterias
func IsValidFirstname(name string) (bool, error) {
	if name == "" {
		return false, fmt.Errorf(T("firstname.must_not_be_empty"))
	}

	// Name can be composed (not mandatory with the ()), will start will a -
	// * means that the group can be repeated
	nameRegex := `^[a-zA-ZéèêçàÉÈÊÇÀßöäüÖÜÄ]+(-[a-zA-ZéèêçàÉÈÊÇÀßöäüÖÜÄ]+)*$`

	// Compile the regex
	re := regexp.MustCompile(nameRegex)

	// Verify if the string matches the regex
	if re.MatchString(name) {
		return re.MatchString(name), nil
	} else {
		return re.MatchString(name), fmt.Errorf(T("firstname.must_be_valid"))
	}
}

// isValidName verifies that the name follows some criterias
func IsValidLastname(name string) (bool, error) {
	if name == "" {
		return false, fmt.Errorf(T("lastname.must_not_be_empty"))
	}

	// Name can be composed (not mandatory with the ()), will start will a -
	// * means that the group can be repeated
	nameRegex := `^[a-zA-ZéèêçàÉÈÊÇÀßöäüÖÜÄ]+(-[a-zA-ZéèêçàÉÈÊÇÀßöäüÖÜÄ]+)*$`

	// Compile the regex
	re := regexp.MustCompile(nameRegex)

	// Verify if the string matches the regex
	if re.MatchString(name) {
		return re.MatchString(name), nil
	} else {
		return re.MatchString(name), fmt.Errorf(T("lastname.must_be_valid"))
	}
}

// Verify if the string is numbers only
func isNumbersOnly(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// isValidMaterialName verifies that the name follows some criterias
func isValidMaterialName(s string) (bool, error) {
	sRegex := `^[a-zA-Z0-9éèêçàÉÈÊÇÀßöäüÖÜÄ ]+([- ][a-zA-Z0-9éèêçàÉÈÊÇÀßöäüÖÜÄ ]+)*$`

	// Compile the regex
	re := regexp.MustCompile(sRegex)

	// Verify if the string is a regex
	if re.MatchString(s) {
		return re.MatchString(s), nil
	} else {
		return re.MatchString(s), fmt.Errorf(T("material_name_invalid"))
	}
}

// standardizeSpaces removes spaces at the beginning and end of the string and replaces multiple spaces by one
func standardizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

// setTitle sets the string as a title for the page. The string is centered, respects dark/light mode and has its size
func setTitle(s string, size float32) *canvas.Text {
	themeColor := fyne.CurrentApp().Settings().Theme().Color("foreground",
		func() fyne.ThemeVariant {
			if darkTheme.IsActivated {
				return theme.VariantDark
			} else {
				return theme.VariantLight
			}
		}())

	title := canvas.NewText(s, themeColor)
	title.Alignment = fyne.TextAlignCenter
	title.TextSize = size
	return title
}

// setText sets the string as a text for the page. The string respects dark/light mode and has its size
func sexText(s string, size float32) *canvas.Text {
	a := fyne.CurrentApp()
	themeColor := a.Settings().Theme().Color("foreground",
		func() fyne.ThemeVariant {
			if darkTheme.IsActivated {
				return theme.VariantDark
			} else {
				return theme.VariantLight
			}
		}())
	title := canvas.NewText(s, themeColor)
	title.Alignment = fyne.TextAlignCenter
	title.TextSize = size
	return title
}

// loadTheme sets the flags for the light Theme and the dark Theme
func loadTheme(a fyne.App) {
	if a.Settings().ThemeVariant() == 1 {
		lightTheme.IsActivated = true
	} else {
		// Dark Theme on default
		darkTheme.IsActivated = true
	}
}

// loadTheme sets the flags for the light Theme and the dark Theme, browser
func loadThemeWeb(a fyne.App) {
	darkTheme.IsActivated = true
	a.Settings().SetTheme(&darkTheme)
}

// loadingSpinnerStart starts a spinner gif that indicates that an operation is working
func loadingSpinnerStart(a fyne.App) (fyne.Window, chan struct{}) {
	w := a.NewWindow("Loading")

	gifData, err := loadGIF()
	if err != nil {
		dialog.ShowError(err, w)
		return w, nil
	}

	stopChan := make(chan struct{})

	gifImage := startGIFAnimation(gifData, stopChan)

	w.Resize(fyne.NewSize(300, 100))
	w.SetContent(gifImage)
	w.CenterOnScreen()
	w.SetCloseIntercept(func() {
		// User can't close
	})
	w.Show()

	return w, stopChan
}

// loadDatabaseWithSpinner uses a spinner gif to indicate that the database is loading
func loadDatabaseWithSpinner(a fyne.App, credToken string) (*mt.Database, error) {
	w, stopChan := loadingSpinnerStart(a)

	resultChan := make(chan struct {
		db  *mt.Database // Remplace par ton type
		err error
	})

	go func() {
		db, err := mr.LoadDB(credToken) // Ton chargement de DB
		resultChan <- struct {
			db  *mt.Database
			err error
		}{db, err}
	}()

	timeout := time.After(1 * time.Minute)

	select {
	case result := <-resultChan:
		close(stopChan) // Arrête l'animation
		w.Close()
		return result.db, result.err
	case <-timeout:
		close(stopChan) // Arrête l'animation
		w.Close()
		return nil, errors.New("operation timed out")
	}
}

// loadDatabaseWithSpinner uses a spinner gif to indicate that the database is loading
func loginWithSpinner(a fyne.App, username string, password string) (*mt.Database, string, error) {
	w, stopChan := loadingSpinnerStart(a)

	resultChan := make(chan struct {
		db    *mt.Database
		token string
		err   error
	})

	go func() {
		db, token, err := mr.Login(username, password)
		resultChan <- struct {
			db    *mt.Database
			token string
			err   error
		}{db, token, err}
	}()

	timeout := time.After(1 * time.Minute)

	select {
	case result := <-resultChan:
		close(stopChan) // Arrête l'animation
		w.Close()
		return result.db, result.token, result.err
	case <-timeout:
		close(stopChan) // Arrête l'animation
		w.Close()
		return nil, "", errors.New("operation timed out")
	}
}

// loadGIF loads a gif
func loadGIF() (*gif.GIF, error) {

	// Load gif as a static resource
	gifResource, err := fyne.LoadResourceFromPath("wasm/spinner_dark.gif")
	if err != nil {
		return nil, err
	}

	g, err := gif.DecodeAll(bytes.NewReader(gifResource.Content()))
	if err != nil {
		return nil, err
	}
	return g, nil
}

// startGIFAnimation creates an animation from a gif
func startGIFAnimation(gifData *gif.GIF, stopChan chan struct{}) *canvas.Image {
	// Create an inital image from the first frames
	img := canvas.NewImageFromImage(gifData.Image[0])
	img.FillMode = canvas.ImageFillContain
	img.SetMinSize(fyne.NewSize(100, 100))

	// Animation
	go func() {
		frame := 0
		for {
			select {
			case <-time.After(time.Duration(gifData.Delay[frame]*10) * time.Millisecond): // respects gif delays
				frame = (frame + 1) % len(gifData.Image) // Loop on frames
				img.Image = gifData.Image[frame]
				img.Refresh()
			case <-stopChan: // Stops the animation if called
				return
			}
		}
	}()

	return img
}
