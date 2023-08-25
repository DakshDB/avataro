package usecase

import (
	"avataro/domain"
	"bytes"
	"fmt"
	"image/jpeg"
)

type avataroUsecase struct {
}

func AvataroUsecase() domain.AvataroUsecase {
	return &avataroUsecase{}
}

func (u *avataroUsecase) GetAvataro(text string, set *string, backgroundSet *string) []byte {
	hash := getHash(text)
	fmt.Println("hash", hash)

	assets := getAssetFromHash(hash, backgroundSet, set, nil)
	fmt.Println("res", assets)

	image := generateImage(assets)

	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, image, nil)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	imageBytes := buf.Bytes()

	return imageBytes
}
