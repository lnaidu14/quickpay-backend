package helpers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/skip2/go-qrcode"
)

var ExampleUser = User{"c85b4c8e-07ab-4c02-849d-71d495d6f905", "Lalit", "+11234567890"}

func GenQrCode(user User) string {
	usr, err := json.Marshal(ExampleUser)

	fileName := fmt.Sprintf("qr-%v.png", user.ID)

	qrcode.WriteFile(string(usr), qrcode.Medium, 256, fileName)
	if err != nil {
		fmt.Println("Error occured while generating QR code")
		return ""
	}
	str := QrCodeToBase64(fileName)
	return str

}

func QrCodeToBase64(path string) string {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	var base64Encoding string
	mimeType := http.DetectContentType(bytes)
	switch mimeType {
	case "image/jpeg":
		base64Encoding += "data:image/jpeg;base64,"
	case "image/png":
		base64Encoding += "data:image/png;base64,"
	default:
		base64Encoding = ""
	}
	base64Encoding += base64.StdEncoding.EncodeToString(bytes)
	if base64Encoding != "" {
		e := os.Remove(path)
		if e != nil {
			log.Fatal(e)
		}
	}
	return base64Encoding
}
