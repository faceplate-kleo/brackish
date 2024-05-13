package brackish

import (
	"bufio"
	"fmt"
	"gopkg.in/yaml.v3"
	"math/rand"
	"os"
	"path"
)

type Bracket struct {
	State     []Match  `yaml:"state"`
	ByPlayers []Player `yaml:"byPlayers"`
	SaveNum   int      `yaml:"saveNum"`
}

type Team struct {
	PlayerA Player `yaml:"playerA"`
	PlayerB Player `yaml:"playerB"`
}

type Match struct {
	TeamA Team `yaml:"teamA"`
	TeamB Team `yaml:"teamB"`
}

func (b *Bracket) SetStateFromSlice(players []Player) {
	//assume players are already shuffled
	b.State = []Match{}
	b.ByPlayers = []Player{}

	teams := make([]Team, 0)

	for i := 0; i < len(players)-1; i += 2 {
		playerA := players[i]

		if i == len(players)-1 {
			b.ByPlayers = append(b.ByPlayers, playerA)
			i--
			break
		}

		playerB := players[i+1]

		teams = append(teams, Team{playerA, playerB})
	}

	if len(teams)%2 != 0 {
		byTeam := teams[len(teams)-1]
		teams = teams[:len(teams)-1]

		b.ByPlayers = append(b.ByPlayers, byTeam.PlayerA)
		b.ByPlayers = append(b.ByPlayers, byTeam.PlayerB)
	}

	for i := 0; i < len(teams)-1; i += 2 {
		TeamA := teams[i]
		TeamB := teams[i+1]

		b.State = append(b.State, Match{TeamA, TeamB})
	}

	for _, byPlayer := range b.ByPlayers {
		fmt.Printf("Player %s gets a by this round!\n", byPlayer.Nickname)
	}

	//for _, match := range b.State {
	//	fmt.Println(match)
	//}
}

func (b *Bracket) Save() error {
	stringData, err := yaml.Marshal(b)
	if err != nil {
		return err
	}

	cwd := os.Getenv("PWD")

	err = os.WriteFile(path.Join(cwd, fmt.Sprintf("brackishState-%d.yaml", b.SaveNum)), stringData, 0664)
	if err != nil {
		return err
	}

	b.SaveNum++
	return nil
}

func LoadBracket(saveNum int) (Bracket, error) {
	cwd := os.Getenv("PWD")
	fileBytes, err := os.ReadFile(path.Join(cwd, fmt.Sprintf("brackishState-%d.yaml", saveNum)))
	if err != nil {
		return Bracket{}, err
	}

	var NewBracket Bracket
	err = yaml.Unmarshal(fileBytes, &NewBracket)
	if err != nil {
		return Bracket{}, err
	}

	return NewBracket, nil
}

func (b *Bracket) Play() ([]Player, error) {
	if len(b.State) == 0 {
		return nil, fmt.Errorf("bracket is empty")
	}
	continuing := make([]Player, 0)
	for i, state := range b.State {
		fmt.Printf(
			"PLAY GAME %d: [%s & %s] vs [%s & %s]\n",
			i+1,
			state.TeamA.PlayerA.Nickname,
			state.TeamA.PlayerB.Nickname,
			state.TeamB.PlayerA.Nickname,
			state.TeamB.PlayerB.Nickname,
		)

		winner, loser, err := GetWinner(state.TeamA, state.TeamB)
		if err != nil {
			return nil, err
		}
		winner.PlayerA.Win()
		winner.PlayerB.Win()
		continuing = append(continuing, winner.PlayerA)
		continuing = append(continuing, winner.PlayerB)
		if !loser.PlayerA.Lose() {
			continuing = append(continuing, loser.PlayerA)
		}
		if !loser.PlayerB.Lose() {
			continuing = append(continuing, loser.PlayerB)
		}
	}

	return append(continuing, b.ByPlayers...), nil
}

func (b *Bracket) Show() {
    for _, match := range b.State {
        fmt.Println("-----------")
        fmt.Printf("%v &\n", match.TeamA.PlayerA.Nickname)
        fmt.Printf("%v\n", match.TeamA.PlayerB.Nickname)
        fmt.Println("  \u21D1")
        fmt.Println("  VS ---------->")
        fmt.Println("  \u21D3")
        fmt.Printf("%v &\n", match.TeamB.PlayerA.Nickname)
        fmt.Printf("%v\n", match.TeamB.PlayerB.Nickname)
        fmt.Println("-----------")
    }
}

func GetWinner(TeamA, TeamB Team) (loser Team, winner Team, err error) {
	fmt.Printf(
		"WHO WON?\n1. [%s & %s]\n2. [%s & %s]\n",
		TeamA.PlayerA.Nickname,
		TeamA.PlayerB.Nickname,
		TeamB.PlayerA.Nickname,
		TeamB.PlayerB.Nickname,
	)

	var choice int

	scanner := bufio.NewScanner(os.Stdin)

	for choice == 0 {
		fmt.Print("Enter 1 or 2: ")
		scanner.Scan()
		inputStr := scanner.Text()
		if inputStr == "1" {
			choice = 1
		} else if inputStr == "2" {
			choice = 2
		}
	}

	if choice == 1 {
		return TeamA, TeamB, nil
	}
	return TeamB, TeamA, nil
}

func Shuffle(players []Player, n int) {
	for range n {
		indA := rand.Intn(len(players))
		indB := rand.Intn(len(players))
		tmp := players[indA]
		players[indA] = players[indB]
		players[indB] = tmp
	}
}
