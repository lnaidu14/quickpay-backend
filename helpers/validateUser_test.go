package helpers

import (
	"quickpay/main/types"
	"testing"
)

var ExampleUser = types.User{Id: "c85b4c8e-07ab-4c02-849d-71d495d6f905", Username: "Foobar", Ph: "+11234567890"}

func TestUser_ValidateUser(t *testing.T) {
	t.Run("Should return a success if user is valid", func(t *testing.T) {
		got := ValidateUser(ExampleUser)
		expected := "User is valid"

		if got != expected {
			t.Errorf("got %q expected %q", got, expected)
		}
	})

	// t.Run("Should throw an error if 'Id' is not of type 'uuid4'", func(t *testing.T) {
	// 	user := types.User{Id: "1234", Username: "Foobar", Ph: "+11234567890"}
	// 	got := ValidateUser(user)
	// 	expected := "Key: 'User.Id' Error:Field validation for 'Id' failed on the 'uuid4' tag"

	// 	if got != expected {
	// 		t.Errorf("got %q expected %q", got, expected)
	// 	}
	// })

	t.Run("Should throw an error if 'Username' is below the 'min' requirement", func(t *testing.T) {
		user := types.User{Id: "c85b4c8e-07ab-4c02-849d-71d495d6f905", Username: "Fo", Ph: "+11234567890"}
		got := ValidateUser(user)
		expected := "Key: 'User.Username' Error:Field validation for 'Username' failed on the 'min' tag"

		if got != expected {
			t.Errorf("got %q expected %q", got, expected)
		}
	})

	t.Run("Should throw an error Username is above the 'max' requirement", func(t *testing.T) {
		user := types.User{Id: "c85b4c8e-07ab-4c02-849d-71d495d6f905", Username: "abcdefghijklmnopqrstuvwxyz", Ph: "+11234567890"}
		got := ValidateUser(user)
		expected := "Key: 'User.Username' Error:Field validation for 'Username' failed on the 'max' tag"

		if got != expected {
			t.Errorf("got %q expected %q", got, expected)
		}
	})

	// t.Run("Should throw an error if Ph does not follow the 'E.164' format", func(t *testing.T) {
	// 	user := types.User{Id: "c85b4c8e-07ab-4c02-849d-71d495d6f905", Username: "Foobar", Ph: "1234"}
	// 	got := ValidateUser(user)
	// 	expected := "Key: 'User.Ph' Error:Field validation for 'Ph' failed on the 'e164' tag"

	// 	if got != expected {
	// 		t.Errorf("got %q expected %q", got, expected)
	// 	}
	// })
}
