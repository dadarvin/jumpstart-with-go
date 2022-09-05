package image

import (
	"bufio"
	"encoding/base64"
	_ "image/jpeg"
	"os"
)

// fgetbase64 Gets base64 string of an existing JPEG file
// fungsi utk mengamfile file berdasarkan nama user, utk diconversi kebase64cide dan dikirim ke clien
func Fgetbase64(fileName string) (string, error) {
	filename := "assets/Pict_" + fileName + ".png"
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
