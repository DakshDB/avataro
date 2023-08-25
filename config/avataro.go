package config

import (
	"fmt"
	"log"
	"os"
)

type avatarBodyParts []string
type colorSet map[string]avatarBodyParts
type avatarSet map[string]colorSet

var Backgrounds = map[string][]string{}
var Sets = map[string]avatarSet{}

func InitializeConfig() {
	getBackgroundDetails()
	getSets()
}

// getBackgroundDetails function
func getBackgroundDetails() {
	//	Read all folders in backgrounds folder
	var sets []string
	entries, err := os.ReadDir("./config/backgrounds")
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range entries {
		fmt.Println(e.Name())
		sets = append(sets, e.Name())
	}

	//	For each folder, read all files
	for _, set := range sets {
		var files []string
		entries, err := os.ReadDir("./config/backgrounds/" + set)
		if err != nil {
			log.Fatal(err)
		}

		for _, e := range entries {
			fmt.Println(e.Name())
			files = append(files, e.Name())
		}

		Backgrounds[set] = files
	}
}

// getSets function
func getSets() {
	// Read all folders in sets folder
	var _sets []string
	entries, err := os.ReadDir("./config/sets")
	if err != nil {
		log.Fatal(err)
	}

	// Set Folder
	for _, e := range entries {
		fmt.Println(e.Name())
		_sets = append(_sets, e.Name())
	}

	// Get all files in each set folder
	for _, set := range _sets {
		var colors []string
		entries, err := os.ReadDir("./config/sets/" + set)
		if err != nil {
			log.Fatal(err)
		}

		for _, e := range entries {
			fmt.Println(e.Name())
			colors = append(colors, e.Name())
		}

		// Get body parts for each color
		for _, color := range colors {
			var bodyParts []string
			entries, err := os.ReadDir("./config/sets/" + set + "/" + color)
			if err != nil {
				log.Fatal(err)
			}

			for _, e := range entries {
				fmt.Println(e.Name())
				bodyParts = append(bodyParts, e.Name())
			}

			// Get each body parts files
			for _, bodyPart := range bodyParts {
				var files []string
				entries, err := os.ReadDir("./config/sets/" + set + "/" + color + "/" + bodyPart)
				if err != nil {
					log.Fatal(err)
				}

				for _, e := range entries {
					fmt.Println(e.Name())
					files = append(files, e.Name())
				}

				if Sets[set] == nil {
					Sets[set] = make(avatarSet)
				}
				if Sets[set][color] == nil {
					Sets[set][color] = make(colorSet)
				}
				if Sets[set][color][bodyPart] == nil {
					Sets[set][color][bodyPart] = make(avatarBodyParts, 0)
				}
				Sets[set][color][bodyPart] = files
			}
		}
	}

}
