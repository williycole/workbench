package blog_snippets

import (
	"bytes"
	"fmt"

	"github.com/sahilm/fuzzy"
)

var gundamCsv = bytes.NewReader([]byte(
	"name,age,pilot\nAmuro Ray,15,EF\nChar Aznable,20,Z\n",
))

var pokemonCsv = bytes.NewReader([]byte(
	"name,type,power\nPikachu,Electric,55\nBulbasaur,Grass/Poison,49\nCharizard,Fire/Flying,84\n",
))

var warHammerCsv = bytes.NewReader([]byte(
	"name,faction,power\nRoboute Guilliman,Ultramarines,95\nAbaddon the Despoiler,Black Legion,98\nEldrad Ulthran,Craftworld Ulthwé,90\n",
))

func main() {
	const bold = "\033[1m%s\033[0m"
	pattern := "mnr"
	data := []string{"game.cpp", "moduleNameResolver.ts", "my name is_Ramsey"}

	matches := fuzzy.Find(pattern, data)

	for _, match := range matches {
		for i := 0; i < len(match.Str); i++ {
			if contains(i, match.MatchedIndexes) {
				fmt.Print(fmt.Sprintf(bold, string(match.Str[i])))
			} else {
				fmt.Print(string(match.Str[i]))
			}
		}
		fmt.Println()
	}
}

func contains(needle int, haystack []int) bool {
	for _, i := range haystack {
		if needle == i {
			return true
		}
	}
	return false
}
