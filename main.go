package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gizak/termui"
)

var (
	client     *http.Client
	APIURL     = "https://community.steam-api.com/%s/v0001/"
	selfInfo   SelfInfo
	textPos    int
	dMap       = map[int]int{1: 600, 2: 1200, 3: 2400, 4: 2400}
	dName      = map[int]string{1: "easy", 2: "medium", 3: "hard", 4: "boss"}
	pMap       map[string]PlanetDetail
	nextPlanet string
	nextZone   int
	difficulty int
	cache      time.Time
	token      string
)

func init() {
	client = &http.Client{}
	pMap = make(map[string]PlanetDetail)
	res, err := ioutil.ReadFile("token.txt")
	if err != nil {
		fmt.Println("Couldn't read token", err)
		os.Exit(1)
	}
	token = string(res)
}

func spost(path string, form url.Values) []byte {
	form.Add("access_token", token)
	req, err := http.NewRequest("POST", fmt.Sprintf(APIURL, path), strings.NewReader(form.Encode()))
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("Origin", "https://steamcommunity.com")
	req.Header.Set("Referer", "https://steamcommunity.com/saliengame/play")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3464.0 Safari/537.36")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	//	log.Println(string(body))
	return body
}

func sget(path string, params string) []byte {
	req, err := http.NewRequest("GET", fmt.Sprintf(APIURL, path)+params, nil)
	req.Header.Set("Origin", "https://steamcommunity.com")
	req.Header.Set("Referer", "https://steamcommunity.com/saliengame/play")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3464.0 Safari/537.36")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	//	log.Println(string(body))
	return body
}

func joinPlanet(id string) {
	printStatus("joining planet " + id)
	spost("ITerritoryControlMinigameService/JoinPlanet", url.Values{"id": []string{id}})
}

func leavePlanet() {
	if selfInfo.Response.ActivePlanet != "" {
		printStatus("leaving planet " + selfInfo.Response.ActivePlanet)
		spost("IMiniGameService/LeaveGame", url.Values{"gameid": []string{selfInfo.Response.ActivePlanet}})
	}
}

func leaveZone() {
	if selfInfo.Response.ActiveZoneGame != "" {
		printStatus("leaving zone " + selfInfo.Response.ActiveZoneGame)
		spost("IMiniGameService/LeaveGame", url.Values{"gameid": []string{selfInfo.Response.ActiveZoneGame}})
	}
}

func joinZone(pos int) error {
	printStatus(fmt.Sprintf("joining zone %d", pos))
	res := spost("ITerritoryControlMinigameService/JoinZone", url.Values{"zone_position": []string{strconv.Itoa(pos)}})
	var zoneinfo ZoneInfo
	json.Unmarshal(res, &zoneinfo)
	if zoneinfo.Response.ZoneInfo.ZonePosition == 0 {
		printStatus(fmt.Sprintf("ERROR: zone %d failed", pos))
		return fmt.Errorf("failed")
	}
	printStatus(fmt.Sprintf("OK: zone %d joined", pos))
	printZoneCapture(strconv.Itoa(int(math.Trunc(zoneinfo.Response.ZoneInfo.CaptureProgress*100))) + "%")
	return nil
}

func getPlanets() Planets {
	res := sget("ITerritoryControlMinigameService/GetPlanets", "?active_only=1")
	var planets Planets
	json.Unmarshal(res, &planets)
	for _, p := range planets.Response.Planets {
		if p.State.Captured == false {
			printStatus(p.ID + " " + p.State.Name + " --")
		}
	}
	return planets
}

func getPlanet(ID string) PlanetInfo {
	res := sget("ITerritoryControlMinigameService/GetPlanet", "?id="+ID)
	var planets Planets
	json.Unmarshal(res, &planets)
	planet := planets.Response.Planets[0]
	pMap[ID] = planet
	allZones := make(map[int][]int)
	for _, z := range planet.Zones {
		if z.ZonePosition == 0 {
			continue
		}
		if !z.Captured && z.CaptureProgress < 0.95 {
			allZones[z.Difficulty] = append(allZones[z.Difficulty], z.ZonePosition)
		}
		if !z.Captured && z.Type == 4 {
			allZones[z.Type] = append(allZones[z.Type], z.ZonePosition)
		}
	}
	return PlanetInfo{allZones}
}

func getUncapturedPlanets() string {
	info := make(map[string]PlanetInfo)
	planets := getPlanets()
	for _, planet := range planets.Response.Planets {
		if planet.State.Captured == false {
			planetinfo := getPlanet(planet.ID)
			info[planet.ID] = planetinfo
		}
	}
	type result struct {
		ID    string
		Count int
	}
	results := make(map[int][]result)
	for i := 4; i > 0; i-- {
		for k, v := range info {
			results[i] = append(results[i], result{k, len(v.allZones[i])})
		}
	}
	sort := make(map[int]int)
	ordered := make(map[int]string)
	for i := 4; i > 0; i-- {
		for _, result := range results[i] {
			if result.Count > sort[i] {
				sort[i] = result.Count
				ordered[i] = result.ID
			}
		}
	}
	for i := 4; i > 0; i-- {
		if sort[i] > 0 {
			printStatus(fmt.Sprintf("search returning planet %s with %d zones of difficulty %d", ordered[i], sort[i], i))
			printZonesLeft(strconv.Itoa(sort[i]))
			return ordered[i]
		}
	}
	return ""
}

func refreshPlanetInfo() PlanetInfo {
	if selfInfo.Response.ActivePlanet != "" {
		return getPlanet(selfInfo.Response.ActivePlanet)
	}
	return PlanetInfo{}
}

func getSelfInfo(p bool) {
	res := spost("ITerritoryControlMinigameService/GetPlayerInfo", url.Values{})
	json.Unmarshal(res, &selfInfo)
	if p {
		printStatus(fmt.Sprintf("Planet/Zone: %s %s Level: %d Score: %s Nextscore: %s", selfInfo.Response.ActivePlanet, selfInfo.Response.ActiveZonePosition,
			selfInfo.Response.Level, selfInfo.Response.Score, selfInfo.Response.NextLevelScore))
		printScore(selfInfo.Response.Score)
		printLevel(strconv.Itoa(selfInfo.Response.Level))
		printPlanet(selfInfo.Response.ActivePlanet)
	}
}

func reportScore(addscore int) error {
	res := spost("ITerritoryControlMinigameService/ReportScore", url.Values{"score": []string{strconv.Itoa(addscore)}})
	var score Score
	json.Unmarshal(res, &score)

	if score.Response.NewScore == "" {
		printStatus("ERROR: score update failed")
		return fmt.Errorf("score update failed")
	}
	printScore(score.Response.NewScore)

	if score.Response.OldLevel != score.Response.NewLevel {
		printLevel(strconv.Itoa(score.Response.NewLevel))
	}
	return nil
}

func getNext() (string, int, int) {
	nextPlanet := getUncapturedPlanets()
	nextZone := 0
	difficulty := 0
	info := getPlanet(nextPlanet)
	for k := 4; k > 0; k-- {
		if len(info.allZones[k]) > 0 {
			t := fmt.Sprintf("difficulty: %d available: %#v", k, info.allZones[k])
			printStatus(t)
			nextZone = info.allZones[k][0]
			difficulty = k
			break
		}
	}
	cache = time.Now()
	return nextPlanet, nextZone, difficulty
}

func loop() error {
	printZone("?")
	printPlanet("?")
	printDifficulty("?")
	printCapture("?")
	printNextGrind("?")
	printNextLevel("?")
	printZoneCapture("?")
	//	printZonesLeft("?")
	getSelfInfo(true)
	nextscore, _ := strconv.Atoi(selfInfo.Response.NextLevelScore)
	currentscore, _ := strconv.Atoi(selfInfo.Response.Score)
	todoscore := nextscore - currentscore

	t := fmt.Sprintf("%.2f %% done for next level\n", (float64(currentscore)/float64(nextscore))*100)
	perc := math.Trunc(float64(currentscore) / float64(nextscore) * 100)
	updateGauge(0)
	updateNextLevelGauge(int(perc))
	printStatus(t)
	leaveZone()
	leavePlanet()

	// only lookup if we're starting up, otherwise we have stuff in cache
	if time.Since(cache) > time.Second*20 {
		nextPlanet, nextZone, difficulty = getNext()
	}
	planetID := nextPlanet
	if planetID == "" {
		printStatus("ERROR: no planetID found")
		return fmt.Errorf("no planetID found")
	}
	joinPlanet(planetID)
	time.Sleep(time.Millisecond * 500)
	getSelfInfo(false)
	if selfInfo.Response.ActivePlanet == planetID {
		printStatus("OK: planet " + planetID + " joined")
		printPlanet(planetID)
	} else {
		printStatus("ERROR: planet " + planetID + " join failed")
		time.Sleep(time.Second * 5)
		return fmt.Errorf("planetjoin failed")
	}
	retryjoincount := 0
retryjoin:
	if retryjoincount == 5 {
		printStatus("ERROR: retryjoin to high, returning. Sleeping 10s")
		time.Sleep(time.Second * 10)
		return fmt.Errorf("retryjoin failed")
	}
	err := joinZone(nextZone)
	if err != nil {
		time.Sleep(time.Second * 5)
		retryjoincount++
		goto retryjoin
	}
	printZone(strconv.Itoa(nextZone))
	printDifficulty(dName[difficulty])
	printCapture(strconv.Itoa(int(math.Trunc(pMap[planetID].State.CaptureProgress*100))) + "%")
	ticker := time.NewTicker(1 * time.Second)
	go func(difficulty int) {
		i := 0
		for range ticker.C {
			i++
			perc := math.Trunc(float64(i) / float64(110) * 100)
			printNextGrind(strconv.Itoa(110-i) + "s")
			score := dMap[difficulty]
			printNextLevel(strconv.Itoa(todoscore/(score/110)-i) + "s")
			updateGauge(int(perc))

			// proactive next planet
			if i == 105 {
				nextPlanet, nextZone, difficulty = getNext()
			}
		}
	}(difficulty)
	time.Sleep(time.Second * 110)
	ticker.Stop()
	score := dMap[difficulty]
	printStatus(fmt.Sprintf("next level will be in %d seconds", todoscore/(score/110)))
	printNextLevel(strconv.Itoa(todoscore/(score/110)) + "s")
	retryscorecount := 0
retryscore:
	if retryscorecount == 5 {
		printStatus("ERROR: retryscore to high, returning. Sleeping 10s")
		time.Sleep(time.Second * 10)
		return fmt.Errorf("retryscore failed")
	}
	err = reportScore(score)
	if err != nil {
		time.Sleep(time.Second * 5)
		retryscorecount++
		goto retryscore
	}
	return nil
}

func main() {
	err := termui.Init()
	if err != nil {
		panic(err)
	}
	defer termui.Close()
	termui.Handle("/sys/kbd/q", func(termui.Event) {
		termui.StopLoop()
		os.Exit(0)
	})
	printText("Info", "Press q to exit", 40, 0, 15)
	go func() {
		for {
			err := loop()
			if err != nil {
				time.Sleep(time.Second * 10)
			}
			updateGauge(0)
			printStatus("OK: starting new loop")
		}
	}()
	termui.Loop()
}
