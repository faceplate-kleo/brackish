package brackish

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Player struct {
	Name     string `yaml:"name"`
	Nickname string `yaml:"nickname"`
	Wins     int
	Losses   int
}

func NewPlayer(name, nickname string) *Player {
	return &Player{
		Name:     name,
		Nickname: nickname,
		Wins:     0,
		Losses:   0,
	}
}

func (p *Player) Win() {
	p.Wins++
}

func (p *Player) Lose() bool {
	p.Losses++
	if p.Losses >= 2 {
		fmt.Printf("Player %s is eliminated!\n", p.Nickname)
		return true
	}
	return false
}

func LoadPlayerFile(filename string) ([]Player, error) {
	rawData, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	players := make([]Player, 0)

	err = yaml.Unmarshal(rawData, &players)
	if err != nil {
		return nil, err
	}

	return players, nil
}
