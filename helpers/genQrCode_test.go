package helpers

import (
	"encoding/json"
	"fmt"
	"runtime"
	"testing"

	"github.com/skip2/go-qrcode"
)

var unsupportedImgName = "../test/mock/img/unsupportedImg.jpg"
var supportedImgName = "../test/mock/img/supportedImg.png"

func setup(file string) {
	user := User{"c85b4c8e-07ab-4c02-849d-71d495d6f905", "Foobar", "+11234567890"}
	usr, err := json.Marshal(user)

	if err != nil {
		fmt.Println("Error occured while generating QR code")
		return
	}

	qrcode.WriteFile(string(usr), qrcode.Medium, 256, supportedImgName)
}

// TODO: Files not deleting probably due to the test case finishing before the file can be deleted
// func teardown(file string) {
// 	os.Remove(file)
// }

func TestString_QrCodeToBase64(t *testing.T) {
	// setup(unsupportedImgName)

	t.Run("Should return an error if file path doesn't exist", func(t *testing.T) {
		got := QrCodeToBase64("invalidFilePath")
		os := runtime.GOOS
		var expected string
		if os == "windows" {
			expected = "open invalidFilePath: The system cannot find the file specified."
		} else if os == "linux" {
			expected = "open invalidFilePath: no such file or directory"
		}

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

	t.Run("Should successfully return a base64 string", func(t *testing.T) {
		setup(supportedImgName)
		got := QrCodeToBase64(supportedImgName)
		expected := "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAQAAAAEAAQMAAABmvDolAAAABlBMVEX///8AAABVwtN+AAACv0lEQVR42uyZQW6sMBBEy2LhpY/QN7EvhoCIi3lu4iN46QVyfXUzyfBzAYwU/9V8vUjpoagqd/B3/s4DjyfJHDJcXcssZENzfOn/1gcBBfA5vOoS9gIBSET/qsA0EkD6FzMWTNxJtgToSNxHAyqw6CeZ6Nl0ihEBp49fVAUN/tXc4wDAxlx5yCYQJES4X4q6G9Df/FVj+Kqz7OUQn/VZ/Ho3bwbsuKr/IIccAFr0+ZfV3QyoHpiRuAVy1w9ktym2gQC21CJioDntVEy0iT18i/YJAMDsO2LogWUqAJIZiOOBYQBPnxuCWrEKZio+wzWEjnkgoPjsWR03rGUqh+jXjBi28JO89wM6B1xNZJ1lk0OTF75jreuDAF9a8j2oGc+YOBGJbPH6LAYASMKxB6qMhQVJkW/7GATQKTQveti5l4nNpoD72OAIAGCetWHGbGqwlFuu0Xw3APHWYeoSNkteJIu0JXy73BMAX5gb4Pj11rX6SXO8yP5+wKbIputNduoU6g/h66fkDAB4aooBoYcDK3WKlnSKTzTfD1j2nu1wol7FPN/V7NLUhgdUMEh2U5tlP8fUqr5gJIBe/eK8QaqxtfPVWz998n4AYiWHGrRCDTU1PbJfXG4AwIKtQ01sk/PV04vYx+XuB/T60FI9M2wyf2iupbrWy4rjAUDTMdXnWIDDBKOX4PDJ7vsBaqupiXbn5SFe20PEEi5WfDtQtHjZikMvYkCD756XyjsE0HQK04NQ9YAW7WeeBbzfzZW7SkbOMS+xOADw3oFEWy7MYruG6DOWkYD3MtkkUyD6dSNp6f1xuRGA773cYhamjtacBTFGAmyZzNO1xOpB/F8wzwGyBp1Qw0Nrev9/cTcGUE+HmHicrYcfZAzg/RcQx81WoLa379f1wgDAuUzWsqgOcfawBlyS9wHA3/k7Q51/AQAA///4Z1Qi/poq9AAAAABJRU5ErkJggg=="

		if got != expected {
			t.Errorf("got %q expected %q", got, expected)
		}
	})

	// teardown(unsupportedImgName)
}

func TestString_GenQrCode(t *testing.T) {
	// setup(unsupportedImgName)

	t.Run("Should successfully return a base64 string", func(t *testing.T) {
		user := User{"c85b4c8e-07ab-4c02-849d-71d495d6f905", "Foobar", "+11234567890"}
		got := GenQrCode(user)
		expected := "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAQAAAAEAAQMAAABmvDolAAAABlBMVEX///8AAABVwtN+AAACv0lEQVR42uyZQW6sMBBEy2LhpY/QN7EvhoCIi3lu4iN46QVyfXUzyfBzAYwU/9V8vUjpoagqd/B3/s4DjyfJHDJcXcssZENzfOn/1gcBBfA5vOoS9gIBSET/qsA0EkD6FzMWTNxJtgToSNxHAyqw6CeZ6Nl0ihEBp49fVAUN/tXc4wDAxlx5yCYQJES4X4q6G9Df/FVj+Kqz7OUQn/VZ/Ho3bwbsuKr/IIccAFr0+ZfV3QyoHpiRuAVy1w9ktym2gQC21CJioDntVEy0iT18i/YJAMDsO2LogWUqAJIZiOOBYQBPnxuCWrEKZio+wzWEjnkgoPjsWR03rGUqh+jXjBi28JO89wM6B1xNZJ1lk0OTF75jreuDAF9a8j2oGc+YOBGJbPH6LAYASMKxB6qMhQVJkW/7GATQKTQveti5l4nNpoD72OAIAGCetWHGbGqwlFuu0Xw3APHWYeoSNkteJIu0JXy73BMAX5gb4Pj11rX6SXO8yP5+wKbIputNduoU6g/h66fkDAB4aooBoYcDK3WKlnSKTzTfD1j2nu1wol7FPN/V7NLUhgdUMEh2U5tlP8fUqr5gJIBe/eK8QaqxtfPVWz998n4AYiWHGrRCDTU1PbJfXG4AwIKtQ01sk/PV04vYx+XuB/T60FI9M2wyf2iupbrWy4rjAUDTMdXnWIDDBKOX4PDJ7vsBaqupiXbn5SFe20PEEi5WfDtQtHjZikMvYkCD756XyjsE0HQK04NQ9YAW7WeeBbzfzZW7SkbOMS+xOADw3oFEWy7MYruG6DOWkYD3MtkkUyD6dSNp6f1xuRGA773cYhamjtacBTFGAmyZzNO1xOpB/F8wzwGyBp1Qw0Nrev9/cTcGUE+HmHicrYcfZAzg/RcQx81WoLa379f1wgDAuUzWsqgOcfawBlyS9wHA3/k7Q51/AQAA///4Z1Qi/poq9AAAAABJRU5ErkJggg=="

		if got != expected {
			t.Errorf("got %q expected %q", got, expected)
		}
	})

	// teardown(unsupportedImgName)
}
