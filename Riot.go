package riot

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	apikey string
	// ErrAPIKeyNotSet - Default error if API key is not set
	ErrAPIKeyNotSet = errors.New("[Riot.go] API Key is not set")
	smallRateChan   rateChan
	longRateChan    rateChan
)

const (
	uriAPIBase = "api.pvp.net/api"
	//BR - Brazil server name
	BR = "br"
	//EUNE - Europe NorthEast server name
	EUNE = "eune"
	//EUW - Europe West server name
	EUW = "euw"
	//KR - Korea server name
	KR = "kr"
	//LAN - Latin America North server name
	LAN = "lan"
	//LAS - Latin America South server name
	LAS = "las"
	//NA - North America server name
	NA = "na"
	//OCE - Oceania server name
	OCE = "oce"
	//RU - Russia server name
	RU = "ru"
	//TR - Turkey server name
	TR = "tr"
	//SEASON3 - League of Legends Season 3
	SEASON3 = "SEASON3"
	//SEASON4 - League of Legends Season 4
	SEASON4 = "SEASON4"
	//RANKEDSOLO5x5 - Ranked Solo 5s
	RANKEDSOLO5x5 = "RANKED_SOLO_5x5"
	//RANKEDTEAM3x3 - Ranked Team 3s
	RANKEDTEAM3x3 = "RANKED_TEAM_3x3"
	//RANKEDTEAM5x5 - Ranked Team 5s
	RANKEDTEAM5x5 = "RANKED_TEAM_5x5"
)

type rateChan struct {
	RateQueue   chan bool
	TriggerChan chan bool
}

//RiotError - structure containing the StatusCode for the reply from Riot's API
type RiotError struct {
	StatusCode int
}

func (err RiotError) Error() string {
	return fmt.Sprintf("Error: HTTP Status %d", err.StatusCode)
}

func SetAPIKey(key string) {
	apikey = key
}

func SetSmallRateLimit(numrequests int, pertime time.Duration) {
	smallRateChan = rateChan{
		RateQueue:   make(chan bool, numrequests),
		TriggerChan: make(chan bool),
	}
	go rateLimitHandler(smallRateChan, pertime)
}

func SetLongRateLimit(numrequests int, pertime time.Duration) {
	longRateChan = rateChan{
		RateQueue:   make(chan bool, numrequests),
		TriggerChan: make(chan bool),
	}
	go rateLimitHandler(longRateChan, pertime)
}

func rateLimitHandler(RateChan rateChan, pertime time.Duration) {
	returnChan := make(chan bool)
	go timeTriggerWatcher(RateChan.TriggerChan, returnChan)
	for {
		<-returnChan
		<-time.After(pertime)
		go timeTriggerWatcher(RateChan.TriggerChan, returnChan)
		length := len(RateChan.RateQueue)
		for i := 0; i < length; i++ {
			<-RateChan.RateQueue
		}
	}
}

func timeTriggerWatcher(timeTrigger chan bool, returnChan chan bool) {
	timeTrigger <- true
	returnChan <- true
}

func getData(url string, v interface{}) (err error) {
	checkRateLimiter(smallRateChan)
	checkRateLimiter(longRateChan)

	res, err := http.Get(url)
	if err != nil {
		return
	}
	checkTimeTrigger(smallRateChan)
	checkTimeTrigger(longRateChan)

	if res.StatusCode != http.StatusOK {
		return RiotError{StatusCode: res.StatusCode}
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, v)
	if err != nil {
		return
	}
	return
}

func checkRateLimiter(RateChan rateChan) {
	if RateChan.RateQueue != nil && RateChan.TriggerChan != nil {
		RateChan.RateQueue <- true
	}
}

func checkTimeTrigger(RateChan rateChan) {
	if RateChan.RateQueue != nil && RateChan.TriggerChan != nil {
		select {
		case <-RateChan.TriggerChan:
		default:
		}
	}
}

type int64ArrayArgs []int64

func (a int64ArrayArgs) String() string {
	str := make([]string, len(a))
	for k, v := range a {
		str[k] = strconv.FormatInt(v, 10)
	}
	return strings.Join(str, ",")
}

func IsKeySet() bool {
	return !(apikey == "")
}
