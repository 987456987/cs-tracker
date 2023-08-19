#### When a player leaves a faceit game for a certain amount of time the api no longer returns any stats on that player
#### My only real option is to either completely ignore those matches which would be annoying or it is to generate all that data myself

# Data I Get From Faceit

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

    matchInfo.MatchID = data.Get("match_id").String()
	matchInfo.Map = data.Get("round_stats.Map").String()
	matchInfo.Winner = data.Get("round_stats.Winner").String()


## What I actually need to genereate
- player.kills
- player.assists
- player.deaths
- player.KDRatio
- player.multikills
- player.headshotpercentage

Most of these I already have event handlers that i can exract the data from. 
Everything to include these stats into the struct to serve is already laid out, I just need to fill it all in


## Faceit Match info api
- player_id
- nickname
- faceit level
- team_id
- Team that won
- MatchID
- Map played

Get this data from faceit api to create playerStats struct and then fill in the rest from the handleDemo function.
This allows me to keep the way I create the structs and manage the data but dont have to worry about getting all the data from the demo which may have more edge cases.