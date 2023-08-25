package usecase

import (
	"avataro/config"
	"crypto/sha512"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/png"
	"math"
	"math/big"
	"os"
	"path"
	"sort"
	"strings"
)

func getImageFromFilePath(filePath string) (image.Image, error) {
	// Check if file exists
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		fmt.Println("File does not exist:", filePath)
		return nil, err
	}

	f, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Println("Error closing file:", err)
		}
	}(f)

	imageDecoded, _, err := image.Decode(f)
	if err != nil {
		fmt.Println("Error decoding image:", err)
		return nil, err
	}

	// Resize image
	imageDecoded = resize(imageDecoded, 256, 256)

	return imageDecoded, err
}

func resize(img image.Image, length int, width int) image.Image {
	//truncate pixel size
	minX := img.Bounds().Min.X
	minY := img.Bounds().Min.Y
	maxX := img.Bounds().Max.X
	maxY := img.Bounds().Max.Y
	for (maxX-minX)%length != 0 {
		maxX--
	}
	for (maxY-minY)%width != 0 {
		maxY--
	}
	scaleX := (maxX - minX) / length
	scaleY := (maxY - minY) / width

	imgRect := image.Rect(0, 0, length, width)
	resImg := image.NewRGBA(imgRect)
	draw.Draw(resImg, resImg.Bounds(), &image.Uniform{C: color.White}, image.Point{}, draw.Src)
	for y := 0; y < width; y += 1 {
		for x := 0; x < length; x += 1 {
			averageColor := getAverageColor(img, minX+x*scaleX, minX+(x+1)*scaleX, minY+y*scaleY, minY+(y+1)*scaleY)
			resImg.Set(x, y, averageColor)
		}
	}
	return resImg
}

func getAverageColor(img image.Image, minX int, maxX int, minY int, maxY int) color.Color {
	var averageRed float64
	var averageGreen float64
	var averageBlue float64
	var averageAlpha float64
	scale := 1.0 / float64((maxX-minX)*(maxY-minY))

	for i := minX; i < maxX; i++ {
		for k := minY; k < maxY; k++ {
			r, g, b, a := img.At(i, k).RGBA()
			averageRed += float64(r) * scale
			averageGreen += float64(g) * scale
			averageBlue += float64(b) * scale
			averageAlpha += float64(a) * scale
		}
	}

	averageRed = math.Sqrt(averageRed)
	averageGreen = math.Sqrt(averageGreen)
	averageBlue = math.Sqrt(averageBlue)
	averageAlpha = math.Sqrt(averageAlpha)

	averageColor := color.RGBA{
		R: uint8(averageRed),
		G: uint8(averageGreen),
		B: uint8(averageBlue),
		A: uint8(averageAlpha)}

	return averageColor
}

func getHash(text string) string {
	//	Calculate the SHA512 hash of the text
	hasher := sha512.New()
	hasher.Write([]byte(text))
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

func getAssetFromHash(hash string, backgroundSet *string, set *string, color *string) map[string]string {
	response := make(map[string]string)

	// Divide the hash into 16 parts
	hashParts := make([]string, 16)
	for i := 0; i < 16; i++ {
		hashParts[i] = hash[i*8 : i*8+8]
	}

	if backgroundSet == nil {
		// Get background set
		_backgroundSet := getBackgroundSet(hashParts[0])
		backgroundSet = &_backgroundSet
	}
	response["backgroundSet"] = *backgroundSet

	background := getBackground(hashParts[1], *backgroundSet)
	response["background"] = background

	if set == nil {
		// Get set
		_set := getSet(hashParts[2])
		set = &_set
	}
	response["set"] = *set

	if color == nil {
		// Get color
		_color := getColor(hashParts[3], *set)
		color = &_color
	}
	response["color"] = *color

	// Possible body parts
	bodyParts := config.Sets[*set][*color]
	bodyPartsNames := make([]string, 0)
	for key := range bodyParts {
		bodyPartsNames = append(bodyPartsNames, key)
	}
	sort.Strings(bodyPartsNames)

	// Get each body part based on the hash
	for bodyPartsIndex, bodyPartName := range bodyPartsNames {
		part := getBodyPart(hashParts[bodyPartsIndex+4], bodyParts[bodyPartName])
		response[bodyPartName] = part
	}

	return response
}

func getBackgroundSet(hexHash string) string {

	// Convert the hash to a number
	digitHash := new(big.Int)
	digitHash.SetString(hexHash, 16)

	backgrounds := config.Backgrounds
	// backgroundSetsCount is the number of keys in the backgrounds map
	backgroundSetsCount := len(backgrounds)
	// backgroundSetIndex is the index of the background set to use
	backgroundSetIndex := digitHash.Mod(digitHash, big.NewInt(int64(backgroundSetsCount))).Uint64()

	var backgroundsSets []string
	for key := range backgrounds {
		backgroundsSets = append(backgroundsSets, key)
	}

	// Sort the backgrounds sets alphabetically
	sort.Strings(backgroundsSets)

	return backgroundsSets[backgroundSetIndex]
}

func getBackground(hexHash string, backgroundSet string) string {
	// Convert the hash to a number
	digitHash := new(big.Int)
	digitHash.SetString(hexHash, 16)

	backgrounds := config.Backgrounds[backgroundSet]
	// backgroundsCount is the number of keys in the backgrounds map
	backgroundsCount := len(backgrounds)
	// backgroundIndex is the index of the background to use
	backgroundIndex := digitHash.Mod(digitHash, big.NewInt(int64(backgroundsCount))).Uint64()

	var backgroundsKeys []string
	for key := range backgrounds {
		backgroundsKeys = append(backgroundsKeys, backgrounds[key])
	}

	// Sort the backgrounds alphabetically
	sort.Strings(backgroundsKeys)

	return backgroundsKeys[backgroundIndex]
}

func getSet(hexHash string) string {
	// Convert the hash to a number
	digitHash := new(big.Int)
	digitHash.SetString(hexHash, 16)

	sets := config.Sets
	// setsCount is the number of keys in the sets map
	setsCount := len(sets)
	// setIndex is the index of the set to use
	setIndex := digitHash.Mod(digitHash, big.NewInt(int64(setsCount))).Uint64()

	var setsKeys []string
	for key := range sets {
		setsKeys = append(setsKeys, key)
	}

	// Sort the sets alphabetically
	sort.Strings(setsKeys)

	return setsKeys[setIndex]
}

func getColor(hexHash string, set string) string {
	// Convert the hash to a number
	digitHash := new(big.Int)
	digitHash.SetString(hexHash, 16)

	colors := config.Sets[set]
	// colorsCount is the number of keys in the colors map
	colorsCount := len(colors)
	// colorIndex is the index of the color to use
	colorIndex := digitHash.Mod(digitHash, big.NewInt(int64(colorsCount))).Uint64()

	var colorsKeys []string
	for key := range colors {
		colorsKeys = append(colorsKeys, key)
	}

	// Sort the colors alphabetically
	sort.Strings(colorsKeys)

	return colorsKeys[colorIndex]
}

func getBodyPart(hexHash string, bodyParts []string) string {
	// Convert the hash to a number
	digitHash := new(big.Int)
	digitHash.SetString(hexHash, 16)

	// bodyPartsCount is the number of keys in the bodyParts map
	bodyPartsCount := len(bodyParts)
	// bodyPartIndex is the index of the bodyPart to use
	bodyPartIndex := digitHash.Mod(digitHash, big.NewInt(int64(bodyPartsCount))).Uint64()

	return bodyParts[bodyPartIndex]
}

func generateImage(assets map[string]string) image.Image {

	backgroundImageSet := assets["backgroundSet"]
	backgroundImageFile := assets["background"]
	backgroundImagePath := path.Join("config/backgrounds/", backgroundImageSet, backgroundImageFile)
	fmt.Println("Background image path:", backgroundImagePath)
	backgroundImage, err := getImageFromFilePath(backgroundImagePath)
	if err != nil {
		fmt.Println("Error getting background image:", err)
		return nil
	}

	//	Get set and color
	set := assets["set"]
	_color := assets["color"]

	//	Get body parts
	bodyPartNames := make([]string, 0)
	for key := range assets {
		if key != "background" && key != "backgroundSet" && key != "set" && key != "color" {
			bodyPartNames = append(bodyPartNames, key)
		}
	}

	fmt.Println("Body part names:", bodyPartNames)

	//	Sort body parts based on the string after # in the body part name
	sort.Slice(bodyPartNames, func(i, j int) bool {
		part1 := strings.Split(bodyPartNames[i], "#")[1]
		part2 := strings.Split(bodyPartNames[j], "#")[1]

		sortedBodyParts := []string{part1, part2}
		sort.Strings(sortedBodyParts)

		return sortedBodyParts[0] == part1
	})

	fmt.Println("Body part names:", bodyPartNames)

	//	Generate the image
	for _, bodyPartName := range bodyPartNames {

		bodyPartFile := assets[bodyPartName]
		bodyPartImagePath := path.Join("config/sets/", set, _color, bodyPartName, bodyPartFile)
		fmt.Println("Body part image path for", bodyPartName, ":", bodyPartImagePath)
		bodyPartImagePath = strings.Replace(bodyPartImagePath, "{set}", set, 1)

		bodyPartImage, err := getImageFromFilePath(bodyPartImagePath)
		if err != nil {
			fmt.Println("Error getting body part image:", err)
			return nil
		}

		backgroundImage = drawImage(backgroundImage, bodyPartImage)
	}

	return backgroundImage
}

func drawImage(backgroundImage image.Image, bodyPartImage image.Image) image.Image {
	newImage := image.NewRGBA(backgroundImage.Bounds())

	draw.Draw(newImage, newImage.Bounds(), backgroundImage, image.Point{}, draw.Src)
	draw.Draw(newImage, newImage.Bounds(), bodyPartImage, image.Point{}, draw.Over)

	return newImage
}
