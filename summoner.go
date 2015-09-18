package riot

import "strings"
import "fmt"

//Summoner - Player data
type Summoner struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	ProfileIconID int    `json:"profileIconId"`
	RevisionDate  int64  `json:"revisionDate"`
	SummonerLevel int64  `json:"summonerLevel"`
}

func SummonerByName(region string, name ...string) (summoners map[string]Summoner, err error) {
	names := strings.Join(name, ",")
	summoners = make(map[string]Summoner)
	args := "api_key=" + apikey
	url := fmt.Sprintf("https://%v.%v/lol/%v/v1.4/summoner/by-name/%v?%v", region, uriAPIBase, region, names, args)
	fmt.Println("[DEBUG] ", url)
	err = getData(url, &summoners)
	if err != nil {
		return
	}
	return
}
