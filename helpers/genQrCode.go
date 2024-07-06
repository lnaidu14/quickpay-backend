package helpers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"quickpay/main/types"

	"github.com/skip2/go-qrcode"
)

func GenQrCode(user types.User) string {
	usr, err := json.Marshal(user)

	fileName := fmt.Sprintf("qr-%v.png", user.Username)

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
		return fmt.Sprintf("%s", err)
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
			return fmt.Sprintf("%s", e)
		}
	} else if base64Encoding == "" {
		return "Image type not supported"
	}
	return base64Encoding
}
