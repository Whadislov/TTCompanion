package my_functions

import (
	mt "github.com/Whadislov/TTCompanion/internal/my_types"
	"github.com/google/uuid"
	"testing"
)

func TestNewUser(t *testing.T) {
	uID := uuid.New()
	d := mt.Database{
		Users: map[uuid.UUID]*mt.User{uID: {
			ID:           uID,
			Name:         "u1",
			Email:        "a1@a1.com",
			PasswordHash: "bvdfvdvver",
			CreatedAt:    "c",
		}},
	}

	expectedLen := 2
	expectedError1 := "username cannot be empty"
	expectedError2 := "username must be valid (only letters and figures are allowed, spaces are not allowed)"
	expectedError3 := "username is already taken"
	expectedError4 := "email cannot be empty"
	expectedError5 := "email must be valid. Example: abc@def.com"
	expectedError6 := "email is already used"
	expectedError7 := "password cannot be empty"
	expectedError8 := "password must be valid (spaces are not allowed)"
	expectedError9 := "passwords do not match"
	expectedError10 := "string is too long"

	t.Run("New user", func(t *testing.T) {
		password := "b"
		_, err1 := NewUser("", "", "", "", &d)
		_, err2 := NewUser(" ", "", "", "", &d)
		_, err3 := NewUser("u1", "", "", "", &d)
		_, err4 := NewUser("u2", "", "", "", &d)
		_, err5 := NewUser("u2", "1", "", "", &d)
		_, err6 := NewUser("u2", "a1@a1.com", "", "", &d)
		_, err7 := NewUser("u2", "a2@a2.com", "", "", &d)
		_, err8 := NewUser("u2", "a2@a2.com", " ", "", &d)
		_, err9 := NewUser("u2", "a2@a2.com", password, "", &d)
		_, err10 := NewUser("ddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd", "a2@a2.com", password, "", &d)
		_, err11 := NewUser("u2", "a2@a2.com", password, password, &d)

		if err1 == nil {
			t.Errorf("Expected error %v, got %v", expectedError1, err1)
		}
		if err2 == nil {
			t.Errorf("Expected error %v, got %v", expectedError2, err2)
		}
		if err3 == nil {
			t.Errorf("Expected error %v, got %v", expectedError3, err3)
		}
		if err4 == nil {
			t.Errorf("Expected error %v, got %v", expectedError4, err4)
		}
		if err5 == nil {
			t.Errorf("Expected error %v, got %v", expectedError5, err5)
		}
		if err6 == nil {
			t.Errorf("Expected error %v, got %v", expectedError6, err6)
		}
		if err7 == nil {
			t.Errorf("Expected error %v, got %v", expectedError7, err7)
		}
		if err8 == nil {
			t.Errorf("Expected error %v, got %v", expectedError8, err8)
		}
		if err9 == nil {
			t.Errorf("Expected error %v, got %v", expectedError9, err9)
		}
		if err10 == nil {
			t.Errorf("Expected error %v, got %v", expectedError10, err9)
		}
		if err11 != nil || len(d.Users) != expectedLen {
			t.Errorf("Expected error %v, got %v", nil, err6)
			t.Errorf("Expected length %v, got %v", 2, len(d.Users))
		}
	})
}
