package main

import (
	"compress/gzip"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	dem "github.com/markus-wa/demoinfocs-golang/v3/pkg/demoinfocs"
	events "github.com/markus-wa/demoinfocs-golang/v3/pkg/demoinfocs/events"
	"github.com/tidwall/gjson"
)

type ServeData struct {
	Matches []Match `json:"matches"`
}

type PlayerStats struct {
	ID                 string `json:"player_id"`
	TeamID             string `json:"team_id"`
	Name               string `json:"nickname"`
	Kills              string `json:"kills"`
	Assists            string `json:"assists"`
	Deaths             string `json:"deaths"`
	KDRatio            string `json:"kd_ratio"`
	MVPs               string `json:"mvps"`
	TripleKills        string `json:"triple_kills"`
	QuadroKills        string `json:"quadro_kills"`
	PentaKills         string `json:"penta_kills"`
	HeadshotPercentage string `json:"headshot_percentage"`
	ADR                string `json:"adr"`
	CounterStrafing    string `json:"counter_strafing"`
	FlashAssists       string `json:"flash_assists"`
	EnemiesFlashed     string `json:"enemies_flashed"`
	TeamFlashed        string `json:"team_flashed"`
	HEDamage           string `json:"he_damage"`
	FireDamage         string `json:"fire_damage"`
	FlashesThrown      string `json:"flashes_thrown"`
}

type Match struct {
	MatchInfo MatchInfo     `json:"match_info"`
	Stats     []PlayerStats `json:"match_stats"`
	TeamInfo  []TeamInfo    `json:"team_info"`
}

type MatchInfo struct {
	MatchID  string `json:"match_id"`
	Map      string `json:"map"`
	Winner   string `json:"winner"`
	Finished string `json:"finished"`
}

type TeamInfo struct {
	TeamID     string `json:"team_id"`
	FinalScore string `json:"final_score"`
}

type FlashEventPair struct {
	Attacker string
	Round    int
}

type RequestData struct {
	UserID string `json:"user_id"`
}

type UserMatches struct {
	UserID  string   `json:"user_id"`
	Matches [][]byte `json:"matches"`
}

var checkMatchesMutex sync.Mutex

func main() {
	http.HandleFunc("/check-matches", func(w http.ResponseWriter, r *http.Request) {
		// Handle CORS headers for preflight and actual requests
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		checkMatchesMutex.Lock()
		defer checkMatchesMutex.Unlock()

		if r.Method == http.MethodOptions {
			// Preflight request, respond with OK
			w.WriteHeader(http.StatusOK)
			return
		} else if r.Method == http.MethodPost {
			// Parse the JSON request body
			var requestData RequestData
			err := json.NewDecoder(r.Body).Decode(&requestData)
			if err != nil {
				http.Error(w, "Error decoding JSON", http.StatusBadRequest)
				return
			}

			filename := requestData.UserID + ".json"

			/////////////////////////////////////////////////////////////////////
			//Get stored data to compare against last matches from faceitAPI
			//Get stored matches from disk
			dataJson, err := os.ReadFile(filename)
			if err != nil {
				fmt.Println(err)
			}

			// Create a variable to hold the unmarshaled JSON data
			var data ServeData

			// Unmarshal the JSON data into the variable
			err = json.Unmarshal(dataJson, &data)
			if err != nil {
				fmt.Println(err)
			}

			/////////////////////////////////////////////////////////////////////

			// Take requestData.UserID and make a faceit api request to get a list of matches
			matchList := getMatchList(requestData.UserID)
			// Var to store the matches missing from stored data
			var missingMatches []string

			for _, faceitMatch := range matchList {
				found := false
				for _, storedMatch := range data.Matches {
					if faceitMatch == storedMatch.MatchInfo.MatchID {
						found = true
						break
					}
				}
				if !found {
					missingMatches = append(missingMatches, faceitMatch)
					fmt.Println("Added " + faceitMatch + " to the needed match list")
				}
			}

			var jsonResponse ServeData

			//Append stored matches to response
			jsonResponse.Matches = append(jsonResponse.Matches, data.Matches...)

			//Extract data from missing matches and append to response variable
			for i, match := range missingMatches {
				matchData := createMatchData(match)
				jsonResponse.Matches = append(jsonResponse.Matches, matchData)
				fmt.Println(strconv.Itoa(i+1) + " out of " + strconv.Itoa(len(missingMatches)))

				// Convert players slice to JSON
				playersJSON, err := json.Marshal(jsonResponse)
				if err != nil {
					fmt.Println("Error converting players to JSON:", err)
					continue // Skip writing and move to the next iteration
				}

				// Save JSON data to the file
				outputFile, err := os.Create(filename)
				if err != nil {
					fmt.Println("Error creating output file:", err)
					continue // Skip writing and move to the next iteration
				}
				defer outputFile.Close()

				_, err = outputFile.Write(playersJSON)
				if err != nil {
					fmt.Println("Error writing JSON data to file:", err)
				}

				fmt.Println("JSON data saved to players.json")
			}

			dataJsonFinal, err := os.ReadFile(filename)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println("JSON data sent")

			//Send the JSON response
			w.WriteHeader(http.StatusOK)
			w.Write(dataJsonFinal)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/get-matches", func(w http.ResponseWriter, r *http.Request) {
		// Handle CORS headers for preflight and actual requests
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			// Preflight request, respond with OK
			w.WriteHeader(http.StatusOK)
			return
		} else if r.Method == http.MethodPost {

			// Parse the JSON request body
			var requestData RequestData
			err := json.NewDecoder(r.Body).Decode(&requestData)
			if err != nil {
				http.Error(w, "Error decoding JSON", http.StatusBadRequest)
				return
			}

			filename := requestData.UserID + ".json"

			/////////////////////////////////////////////////////////////////////
			//Get stored data to compare against last matches from faceitAPI
			//Get stored matches from disk
			dataJson, err := os.ReadFile(filename)
			if err != nil {
				fmt.Println(err)
			}

			//Send the JSON response
			w.WriteHeader(http.StatusOK)
			w.Write(dataJson)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getMatchList(userID string) []string {
	//user := "52e19c3d-471c-44ac-af67-a394da815a37"
	url := "https://open.faceit.com/data/v4/players/" + userID + "/history?game=csgo&offset=0&limit=10"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", "Bearer 8f1de641-442d-4e79-9795-505a59bafca8")

	// Create a custom transport with TLS settings
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // Disabling SSL certificate validation
		},
	}
	client := &http.Client{Transport: tr}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		fmt.Println("API request failed for matches")
	}

	// Read the response body into a byte slice
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
	}

	// Parse JSON using gjson
	jsonStr := gjson.ParseBytes(bodyBytes)

	matches := jsonStr.Get("items").Array()

	var matchList []string

	for _, match := range matches {
		if match.Get("status").String() == "finished" {
			matchList = append(matchList, match.Get("match_id").String())
		}

	}

	return matchList
}

func createMatchData(matchID string) Match {
	matchInfo, teamInfo, allPlayerStats := getFaceitData(matchID)
	allDemoStats := handleDemo(matchID)

	// Loop through allPlayerStats and update with demoStats
	for i := range allPlayerStats {
		allPlayerStats[i].ADR = strconv.Itoa(allDemoStats["ADR"][allPlayerStats[i].Name])
		allPlayerStats[i].CounterStrafing = strconv.Itoa(allDemoStats["CounterStrafing"][allPlayerStats[i].Name])
		allPlayerStats[i].FlashAssists = strconv.Itoa(allDemoStats["FlashAssists"][allPlayerStats[i].Name])
		allPlayerStats[i].EnemiesFlashed = strconv.Itoa(allDemoStats["EnemiesFlashed"][allPlayerStats[i].Name])
		allPlayerStats[i].TeamFlashed = strconv.Itoa(allDemoStats["TeamFlashed"][allPlayerStats[i].Name])
		allPlayerStats[i].HEDamage = strconv.Itoa(allDemoStats["HEDamage"][allPlayerStats[i].Name])
		allPlayerStats[i].FireDamage = strconv.Itoa(allDemoStats["FireDamage"][allPlayerStats[i].Name])
		allPlayerStats[i].FlashesThrown = strconv.Itoa(allDemoStats["FlashesThrown"][allPlayerStats[i].Name])
	}

	//Create match struct and asign
	var match Match
	match.MatchInfo = matchInfo
	match.TeamInfo = teamInfo
	match.Stats = allPlayerStats

	return match
}

func getFaceitData(matchID string) (MatchInfo, []TeamInfo, []PlayerStats) {
	url := "https://open.faceit.com/data/v4/matches/" + matchID + "/stats"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", "Bearer 8f1de641-442d-4e79-9795-505a59bafca8")

	// Create a custom transport with TLS settings
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // Disabling SSL certificate validation
		},
	}
	client := &http.Client{Transport: tr}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		fmt.Println("API request failed")
	}

	// Read the response body into a byte slice
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
	}

	// Parse JSON using gjson
	jsonStr := gjson.ParseBytes(bodyBytes)

	teamInfo, playerData := getPlayerData(jsonStr)

	matchInfo := getMatchInfo(jsonStr)

	return matchInfo, teamInfo, playerData
}

func getPlayerData(jsonStr gjson.Result) ([]TeamInfo, []PlayerStats) {
	teams := jsonStr.Get("rounds.0.teams").Array()

	//Create arrays of both structs
	var players []PlayerStats
	var teamInfo []TeamInfo

	for _, team := range teams {
		//Create TeamInfo struct
		aTeam := TeamInfo{}
		//Assign team values
		aTeam.TeamID = team.Get("team_id").String()
		aTeam.FinalScore = team.Get("team_stats.Final Score").String()
		//Append to created struct
		teamInfo = append(teamInfo, aTeam)

		//Get playerData array from team object
		playersData := team.Get("players").Array()
		for _, playerData := range playersData {
			player := PlayerStats{}
			player.ID = playerData.Get("player_id").String()
			player.Name = playerData.Get("nickname").String()
			player.TeamID = team.Get("team_id").String()
			player.Kills = playerData.Get("player_stats.Kills").String()
			player.Assists = playerData.Get("player_stats.Assists").String()
			player.Deaths = playerData.Get("player_stats.Deaths").String()
			player.KDRatio = playerData.Get("player_stats.K/D Ratio").String()
			player.MVPs = playerData.Get("player_stats.MVPs").String()
			player.TripleKills = playerData.Get("player_stats.Triple Kills").String()
			player.QuadroKills = playerData.Get("player_stats.Quadro Kills").String()
			player.PentaKills = playerData.Get("player_stats.Penta Kills").String()
			player.HeadshotPercentage = playerData.Get("player_stats.Headshots %").String()
			players = append(players, player)
		}
	}

	return teamInfo, players
}

func getMatchInfo(jsonStr gjson.Result) MatchInfo {
	data := jsonStr.Get("rounds.0")

	var matchInfo MatchInfo

	matchInfo.MatchID = data.Get("match_id").String()
	matchInfo.Map = data.Get("round_stats.Map").String()
	matchInfo.Winner = data.Get("round_stats.Winner").String()

	url := "https://open.faceit.com/data/v4/matches/" + matchInfo.MatchID
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", "Bearer 8f1de641-442d-4e79-9795-505a59bafca8")

	// Create a custom transport with TLS settings
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // Disabling SSL certificate validation
		},
	}
	client := &http.Client{Transport: tr}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		fmt.Println("API request failed")
	}

	// Read the response body into a byte slice
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
	}

	// Parse JSON using gjson
	matchRequest := gjson.ParseBytes(bodyBytes)

	matchInfo.Finished = matchRequest.Get("finished_at").String()

	return matchInfo
}

func handleDemo(matchID string) map[string]map[string]int {
	//demoURL := "https://demos-us-central1.faceit-cdn.net/csgo/1-dfedca35-3934-4b58-ba3f-1179dbd3277a-1-1.dem.gz"
	demoURL := "https://demos-us-central1.faceit-cdn.net/csgo/" + matchID + "-1-1.dem.gz"

	// Get the filename from the URL
	tokens := strings.Split(demoURL, "/")
	filename := tokens[len(tokens)-1]

	// Create the output file
	outputFile, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating output file:", err)
	}

	// Download the .gz file
	response, err := http.Get(demoURL)
	if err != nil {
		fmt.Println("Error downloading file:", err)
	}

	// Copy the downloaded data to the output file
	_, err = io.Copy(outputFile, response.Body)
	if err != nil {
		fmt.Println("Error copying data:", err)
	}

	// Close the output file
	if err := outputFile.Close(); err != nil {
		fmt.Println("Error closing output file:", err)
	}

	// Open the .gz file
	compressedDemo, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening .gz file:", err)
	}

	// Create a reader for the .gz file
	gzReader, err := gzip.NewReader(compressedDemo)
	if err != nil {
		fmt.Println("Error creating gzip reader:", err)
	}

	// Create the output file for the extracted content
	extractedFilename := strings.TrimSuffix(filename, ".gz")
	extractedDemo, err := os.Create(extractedFilename)
	if err != nil {
		fmt.Println("Error creating extracted file:", err)
	}

	// Copy the extracted data to the output file
	_, err = io.Copy(extractedDemo, gzReader)
	if err != nil {
		fmt.Println("Error copying extracted data:", err)
	}
	// Close the compressedDemo
	if err = compressedDemo.Close(); err != nil {
		fmt.Println("Error closing file:", err)
	}

	// Delete the compressedDemo .gz file
	if err := os.Remove(filename); err != nil {
		fmt.Println("Error deleting downloaded .gz file:", err)
	}

	// Process the demo and get data
	demoStats := extractDemoData(extractedFilename)

	if err = extractedDemo.Close(); err != nil {
		fmt.Println("Error closing file:", err)
	}
	// Delete the extracted file
	if err := os.Remove(extractedFilename); err != nil {
		fmt.Println("Error deleting extracted file:", err)
	}

	fmt.Println("Extraction and deletion complete.")
	return demoStats
}

func extractDemoData(demoURL string) map[string]map[string]int {
	f, err := os.Open("./" + demoURL)
	if err != nil {
		fmt.Println("Failed to open demo file")
	}
	defer f.Close()

	p := dem.NewParser(f)
	defer p.Close()

	//Find start of match to skip warmup and knife round
	matchStart := 0
	matchStarted := false

	p.RegisterEventHandler(func(e events.MatchStart) {
		matchStart++
		if matchStart == 3 {
			matchStarted = true
		}
	})

	///////////////////////////////////////////////////

	//Keep Round count for average calculations
	numberOfRounds := 1
	p.RegisterEventHandler(func(e events.RoundEndOfficial) {
		numberOfRounds++
	})
	/////////////////////////////////////////////////////////////

	//Calculate Total Damage to find ADR && Calculate Fire and Frag Damage
	playerDamageTotals := make(map[string]int)
	playerFragDamage := make(map[string]int)
	playerFireDamage := make(map[string]int)

	p.RegisterEventHandler(func(e events.PlayerHurt) {
		if e.Attacker != nil && e.Attacker.Team != e.Player.Team && matchStarted {
			playerDamageTotals[e.Attacker.Name] += e.HealthDamageTaken
			if e.Weapon.String() == "HE Grenade" {
				playerFragDamage[e.Attacker.Name] += e.HealthDamageTaken
			}
			if e.Weapon.String() == "Molotov" || e.Weapon.String() == "Incendiary Grenade" {
				playerFireDamage[e.Attacker.Name] += e.HealthDamageTaken
			}
		}
	})
	//////////////////////////////////////////////////////////

	// Calculate Total Flashes Thrown
	playerFlashThrown := make(map[string]int)

	p.RegisterEventHandler(func(e events.GrenadeEventIf) {
		if e.Base().GrenadeType.String() == "Flashbang" && matchStarted {
			playerFlashThrown[e.Base().Thrower.Name]++
		}
	})

	// Calculate enemies flashed
	playerEnemiesFlashed := make(map[string]int)
	playerTeamFlashed := make(map[string]int)
	//flashedEvents := make(map[string]string)

	flashEvents := make(map[string]FlashEventPair)

	p.RegisterEventHandler(func(e events.PlayerFlashed) {
		if e.Attacker != nil && e.Player.IsAlive() && matchStarted {
			if e.Player.IsBlinded() && e.Player.FlashDurationTimeRemaining() > 1100*time.Millisecond {
				if e.Attacker.Team != e.Player.Team {
					playerEnemiesFlashed[e.Attacker.Name]++
					flashEvents[e.Player.Name] = FlashEventPair{Attacker: e.Attacker.Name, Round: numberOfRounds}
				}
				if e.Attacker.Team == e.Player.Team {
					playerTeamFlashed[e.Attacker.Name]++
				}
			}
		}
	})
	////////////////////////////////////////////////////////////////

	// Calculate flash assists
	playerFlashAssists := make(map[string]int)
	p.RegisterEventHandler(func(e events.Kill) {
		if e.Killer != nil && e.Killer.Team != e.Victim.Team && matchStarted {
			if e.Victim.IsBlinded() {
				if flashEvents[e.Victim.Name].Round == numberOfRounds {
					playerFlashAssists[flashEvents[e.Victim.Name].Attacker]++
				}
			}
		}
	})

	///////////////////////////////////////////////////////////////

	// Calculate Counterstrafing
	playerGoodStrafing := make(map[string]int)
	playerStrafingTotal := make(map[string]int)

	p.RegisterEventHandler(func(e events.WeaponFire) {
		//Valid Weapons
		validWeaponNames := []string{
			"AK-47",
			"FAMAS",
			"M4A4",
			"M4A1",
			"Galil AR",
			"SSG 08",
			"AWP",
			"AUG",
			"SG 553",
			"SCAR-20",
			"G3SG1",
		}
		//Max speed for each rifle
		weaponMaxSpeed := map[string]int{
			"AK-47":    73,
			"FAMAS":    74,
			"M4A4":     77,
			"M4A1":     77,
			"Galil AR": 73,
			"SSG 08":   78,
			"AWP":      68,
			"AUG":      74,
			"SG 553":   71,
			"SCAR-20":  73,
			"G3SG1":    73,
		}

		// Check if e.Weapon.String() matches any of the valid weapon names
		weaponName := e.Weapon.String()
		isValidWeapon := false

		for _, validName := range validWeaponNames {
			if weaponName == validName {
				isValidWeapon = true
				break
			}
		}
		// if match is started and the shot was fired from a valid weapon
		if isValidWeapon && matchStarted {
			velocity := e.Shooter.Velocity()
			velocityMagnitude2D := math.Sqrt(velocity.X*velocity.X + velocity.Y*velocity.Y)
			// if the player is not ducking and the player is moving
			if e.Shooter != nil && !e.Shooter.IsDucking() && velocityMagnitude2D > 0 {
				// Loop through all players on the opposing team
				for _, player := range p.GameState().Participants().Playing() {
					// If the shooter has spotted a player from the enemy team
					if player.Team != e.Shooter.Team && player.IsAlive() && player.IsSpottedBy(e.Shooter) {
						if int(velocityMagnitude2D) < weaponMaxSpeed[e.Weapon.String()] {
							playerGoodStrafing[e.Shooter.Name]++
							playerStrafingTotal[e.Shooter.Name]++
							break
						} else {
							playerStrafingTotal[e.Shooter.Name]++
							break
						}
					}
				}
			}
		}
	})

	//////////////////////////////////////////////////////////////

	// Parse to end
	err = p.ParseToEnd()
	if err != nil {
		fmt.Println("Failed to parse demo")
	}

	// Output the adr for each player as a map
	playerDamageMap := make(map[string]int)
	for playerName, totalDamage := range playerDamageTotals {
		playerDamageMap[playerName] = int(float64(totalDamage) / float64(numberOfRounds))
	}

	playerCounterStrafing := make(map[string]int)
	for player, shots := range playerStrafingTotal {
		playerCounterStrafing[player] = int(math.Ceil(float64(playerGoodStrafing[player]) / float64(shots) * 100))
	}

	demoStats := map[string]map[string]int{
		"ADR":             playerDamageMap,
		"FlashesThrown":   playerFlashThrown,
		"EnemiesFlashed":  playerEnemiesFlashed,
		"TeamFlashed":     playerTeamFlashed,
		"FlashAssists":    playerFlashAssists,
		"CounterStrafing": playerCounterStrafing,
		"HEDamage":        playerFragDamage,
		"FireDamage":      playerFireDamage,
	}

	return demoStats
}
