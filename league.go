package riot

import (
	"errors"
	"fmt"
	"strconv"
)

//League - This object contains league information.
type League struct {
	Entries       []LeagueItem `json:"entries"`
	Name          string       `json:"name"`
	ParticipantID string       `json:"participantId"`
	Queue         string       `json:"queue"`
	Tier          string       `json:"tier"`
}

//LeagueItem - This object contains league participant information representing a summoner or team.
type LeagueItem struct {
	Division         string     `json:"division"`
	IsFreshBlood     bool       `json:"isFreshBlood"`
	IsHotStreak      bool       `json:"isHotStreak"`
	IsInactive       bool       `json:"isInactive"`
	IsVeteran        bool       `json:"isVeteran"`
	LeaguePoints     int        `json:"leaguePoints"`
	MiniSeries       MiniSeries `json:"miniSeries"`
	PlayerOrTeamID   string     `json:"playerOrTeamId"`
	PlayerOrTeamName string     `json:"playerOrTeamName"`
	Wins             int        `json:"wins"`
}

//MiniSeries - This object contains mini series information.
type MiniSeries struct {
	Losses   int    `json:"losses"`
	Progress string `json:"progress"`
	Target   int    `json:"target"`
	Wins     int    `json:"wins"`
}

//LeagueEntry - Returns all league entries for specified summoners and summoners' teams.
func LeagueEntry(region string, summonerID ...int64) (leagues map[int64][]League, err error) {
	if len(summonerID) > 10 {
		return leagues, errors.New("[LeagueEntry] Maximux number of summoners requested must be lower than 10")
	}

	leagues = make(map[int64][]League)
	fullData := make(map[string][]League)
	IDs := int64ArrayArgs(summonerID).String()

	args := "api_key=" + apikey
	url := fmt.Sprintf("https://%v.%v/lol/%v/v2.5/league/by-summoner/%v/entry?%v", region, uriAPIBase, region, IDs, args)
	err = getData(url, &fullData)
	if err != nil {
		return nil, err
	}

	for k, v := range fullData {
		id, err := strconv.ParseInt(k, 10, 64)
		if err != nil {
			return nil, err
		}
		leagues[id] = v
	}

	return leagues, err
}
