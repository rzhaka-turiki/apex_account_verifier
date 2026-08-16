package model

type ApexPlayer struct {
	UID           string
	Name          string
	Tag           string
	Platform      string
	Level         int
	LevelPrestige int
}

func (a ApexPlayer) TotalLevel() int {
	return a.Level a.LevelPrestige*500
}