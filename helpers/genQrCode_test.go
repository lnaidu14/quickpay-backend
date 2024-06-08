package helpers

import (
	"testing"
)

var unsupportedImgName = "../test/mock/img/unsupportedImg.jpg"

// func setup(file string) {
// 	os.Create(file)
// }

// TODO: Files not deleting probably due to the test case finishing before the file can be deleted
// func teardown(file string) {
// 	os.Remove(file)
// }

func TestString_QrCodeToBase64(t *testing.T) {
	// setup(unsupportedImgName)

	t.Run("Should return an error if file path doesn't exist", func(t *testing.T) {
		got := QrCodeToBase64("invalidFilePath")
		expected := "open invalidFilePath: The system cannot find the file specified."

		if got != expected {
			t.Errorf("got %q expected %q", got, expected)
		}
	})

	t.Run("Should return an error if the type of image is not supported", func(t *testing.T) {
		got := QrCodeToBase64(unsupportedImgName)
		expected := "Image type not supported"

		if got != expected {
			t.Errorf("got %q expected %q", got, expected)
		}
	})

	// teardown(unsupportedImgName)
}
