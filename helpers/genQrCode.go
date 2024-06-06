package helpers

import (
	"fmt"

	"encoding/json"

	"github.com/skip2/go-qrcode"
)

var ExampleUser = User{"c85b4c8e-07ab-4c02-849d-71d495d6f905", "Lalit", "+11234567890"}

func GenQrCode(user User) {
	ValidateUser(ExampleUser)
	usr, err := json.Marshal(ExampleUser)

	qrcode.WriteFile(string(usr), qrcode.Medium, 256, "qr.png")
	if err != nil {
		fmt.Println("Error occured while generating QR code")
	}
	return (usr)
}
