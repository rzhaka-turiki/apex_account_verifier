package dto

type bridgeResponse struct {
	Global struct {
		UID           string `json:"uid"`
		Name          string `json:"name"`
		Tag           string `json:"tag"`
		Platform      string `json:"platform"`
		Level         int    `json:"level"`
		LevelPrestige int    `json:"levelPrestige"`

		Rank struct {
			RankScore int32  `json:"rankScore"`
			RankName  string `json:"rankName"`
			RankDiv   int32  `json:"rankDiv"`
		} `json:"rank"`
	} `json:"global"`

	Realtime struct {
		IsOnline     int    `json:"isOnline"`
		IsInGame     int    `json:"isInGame"`
		CurrentState string `json:"currentState"`
	} `json:"realtime"`
}
