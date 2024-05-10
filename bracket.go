package brackish

import (
	"fmt"
	"gopkg.in/yaml.v3"
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

func (b *Bracket) setStateFromSlice(players []Player) {
	//assume players are already shuffled

	teams := make([]Team, 0)

	for i := range len(players) - 1 {
		playerA := players[i]

		if i == len(players)-1 {
			b.ByPlayers = append(b.ByPlayers, playerA)
			break
		}

		playerB := players[i+1]

		teams = append(teams, Team{playerA, playerB})
		i++
	}

	if len(teams)%2 != 0 {
		byTeam := teams[len(teams)-1]
		teams = teams[:len(teams)-1]

		b.ByPlayers = append(b.ByPlayers, byTeam.PlayerA)
		b.ByPlayers = append(b.ByPlayers, byTeam.PlayerB)
	}

	for _, byPlayer := range b.ByPlayers {
		fmt.Printf("Player %s gets a by this round!\n", byPlayer.Nickname)
	}
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
