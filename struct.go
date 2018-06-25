package main

type PlanetInfo struct {
	allZones map[int][]int
}

type SelfInfo struct {
	Response struct {
		ActiveZoneGame     string `json:"active_zone_game"`
		ActiveZonePosition string `json:"active_zone_position"`
		ActivePlanet       string `json:"active_planet"`
		TimeOnPlanet       int    `json:"time_on_planet"`
		Score              string `json:"score"`
		Level              int    `json:"level"`
		NextLevelScore     string `json:"next_level_score"`
	} `json:"response"`
}

type Score struct {
	Response struct {
		OldScore       string `json:"old_score"`
		OldLevel       int    `json:"old_level"`
		NewScore       string `json:"new_score"`
		NewLevel       int    `json:"new_level"`
		NextLevelScore string `json:"next_level_score"`
	} `json:"response"`
}

type ZoneInfo struct {
	Response struct {
		ZoneInfo struct {
			CaptureProgress float64 `json:"capture_progress"`
			Captured        bool    `json:"captured"`
			Difficulty      int     `json:"difficulty"`
			Gameid          string  `json:"gameid"`
			Leader          struct {
				Accountid int    `json:"accountid"`
				Avatar    string `json:"avatar"`
				Name      string `json:"name"`
				URL       string `json:"url"`
			} `json:"leader"`
			TopClans []struct {
				Accountid int    `json:"accountid"`
				Avatar    string `json:"avatar"`
				Name      string `json:"name"`
				URL       string `json:"url"`
			} `json:"top_clans"`
			Type         int `json:"type"`
			ZonePosition int `json:"zone_position"`
		} `json:"zone_info"`
	} `json:"response"`
}

type PlanetDetail struct {
	GiveawayApps []int  `json:"giveaway_apps"`
	ID           string `json:"id"`
	State        struct {
		ActivationTime  int     `json:"activation_time"`
		Active          bool    `json:"active"`
		CaptureProgress float64 `json:"capture_progress"`
		Captured        bool    `json:"captured"`
		CloudFilename   string  `json:"cloud_filename"`
		CurrentPlayers  int     `json:"current_players"`
		Difficulty      int     `json:"difficulty"`
		GiveawayID      string  `json:"giveaway_id"`
		ImageFilename   string  `json:"image_filename"`
		LandFilename    string  `json:"land_filename"`
		MapFilename     string  `json:"map_filename"`
		Name            string  `json:"name"`
		Position        int     `json:"position"`
		Priority        int     `json:"priority"`
		TagIds          string  `json:"tag_ids"`
		TotalJoins      int     `json:"total_joins"`
	} `json:"state"`
	TopClans []struct {
		ClanInfo struct {
			Accountid int    `json:"accountid"`
			Avatar    string `json:"avatar"`
			Name      string `json:"name"`
			URL       string `json:"url"`
		} `json:"clan_info"`
		NumZonesControled int `json:"num_zones_controled"`
	} `json:"top_clans"`
	Zones []struct {
		CaptureProgress float64 `json:"capture_progress"`
		Captured        bool    `json:"captured"`
		Difficulty      int     `json:"difficulty"`
		Gameid          string  `json:"gameid"`
		Leader          struct {
			Accountid int    `json:"accountid"`
			Avatar    string `json:"avatar"`
			Name      string `json:"name"`
			URL       string `json:"url"`
		} `json:"leader"`
		TopClans []struct {
			Accountid int    `json:"accountid"`
			Avatar    string `json:"avatar"`
			Name      string `json:"name"`
			URL       string `json:"url"`
		} `json:"top_clans"`
		Type         int `json:"type"`
		ZonePosition int `json:"zone_position"`
	} `json:"zones"`
}

type Planets struct {
	Response struct {
		Planets []PlanetDetail `json:"planets"`
	} `json:"response"`
}
