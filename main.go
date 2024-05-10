package brackish

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var loadNum int
	var playerFile string

	flag.IntVar(&loadNum, "load", -1, "bracket file to load")
	flag.StringVar(&playerFile, "players", "./players.yaml", "yaml file containing player information")
	flag.Parse()

	bracket := Bracket{}
	var err error

	if loadNum != -1 {
		bracket, err = LoadBracket(loadNum)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	players, err := LoadPlayerFile(playerFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	bracket.setStateFromSlice(players)
}
