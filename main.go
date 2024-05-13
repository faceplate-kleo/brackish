package main

import (
	"flag"
	"fmt"
	"github.com/faceplate-kleo/brackish/brackish"
	"os"
)

func main() {
	var loadNum int
	var playerFile string
    var show bool
    var showAndExit bool

	flag.IntVar(&loadNum, "load", -1, "bracket file to load")
	flag.StringVar(&playerFile, "players", "./playerlist.yaml", "yaml file containing player information")
    flag.BoolVar(&show, "show", false, "show the bracket")
    flag.BoolVar(&showAndExit, "showandexit", false, "show the bracket then immediately exit")
	flag.Parse()

	bracket := brackish.Bracket{}
	var err error

	if loadNum != -1 {
		bracket, err = brackish.LoadBracket(loadNum)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	players, err := brackish.LoadPlayerFile(playerFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	brackish.Shuffle(players, 50)

	bracket.SetStateFromSlice(players)
	round := 1

    if show || showAndExit{
        bracket.Show()
        if showAndExit {
            return
        }
    }

	winners := brackish.Team{}

	for {
		fmt.Printf("ROUND %d\n", round)

		remaining, err := bracket.Play()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if len(remaining) <= 2 {
			winners = brackish.Team{remaining[0], remaining[1]}
			break
		}
		round += 1
		brackish.Shuffle(remaining, 50)
		fmt.Println("############################################################\n\n\n\n\n\n\n\n")
		_ = bracket.Save()
		bracket.SetStateFromSlice(remaining)
	}

	fmt.Println("THE WINNERS:")
	fmt.Println(winners.PlayerA.Nickname)
	fmt.Println(winners.PlayerB.Nickname)
}
