package image

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"strings"
)

// Base64toPng convert base 64 to png format
func Base64toPng(fIdUser string, fPicture string) error {

	//fPicture adalah base64cie yg dikirim dari clien utk diubah jadi png
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(fPicture))

	m, formatString, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
		return err
	}
	bounds := m.Bounds()
	fmt.Println(bounds, formatString)

	//Encode from image format to writer
	//fUser=nama user yg dijadinakan file name (nama user unique)
	pngFilename := "Pict_" + fIdUser + ".png"

	f, err := os.OpenFile(pngFilename, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		log.Fatal(err)
		return err
	}

	err = png.Encode(f, m)
	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Println("Png file", pngFilename, "created")

	return nil
}

// fgetbase64 Gets base64 string of an existing JPEG file
// fungsi utk mengamfile file berdasarkan nama user, utk diconversi kebase64cide dan dikirim ke clien
func Fgetbase64(fileName string) (string, error) {

	var filename = "Pict_" + fileName + ".png"
	imgFile, err := os.Open(filename)

	if err != nil {
		return "", err
	}

	defer imgFile.Close()

	// create a new buffer base on file size
	fInfo, err := imgFile.Stat()
	if err != nil {
		return "", err
	}

	var size = fInfo.Size()
	buf := make([]byte, size)
	// read file content into buffer
	fReader := bufio.NewReader(imgFile)
	fReader.Read(buf)

	imgBase64Str := base64.StdEncoding.EncodeToString(buf)

	return imgBase64Str, err
}
