package imageo

import (
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"

	"github.com/nfnt/resize"
)

// CreateThumbnail ...
func CreateThumbnail(path, img, toPath string, maxWidth, maxHeight uint) {
	file, err := os.Open(fmt.Sprintf("%s%s", path, img))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	imgFile, _, err := image.Decode(file)
	if err != nil {
		log.Fatal(err, img)
	}

	newImage := resize.Thumbnail(maxWidth, maxHeight, imgFile, resize.Lanczos3)

	if _, err := os.Stat(toPath); os.IsNotExist(err) {
		os.Mkdir(toPath, 0777)
	}

	toimg, err := os.Create(fmt.Sprintf("%s%s", toPath, img))
	if err != nil {
		log.Fatal(err, img)
	}
	defer toimg.Close()

	// Encode uses a Writer, use a Buffer if you need the raw []byte
	err = jpeg.Encode(toimg, newImage, nil)
	if err != nil {
		log.Fatal(err, img)
	}
}

// GetImageDimensions ...
func GetImageDimensions(fileName, filePath string) (int, int) {
	file, err := os.Open(filePath + fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	defer file.Close()

	image, _, err := image.DecodeConfig(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", filePath+fileName, err)
	}

	return image.Height, image.Width
}
