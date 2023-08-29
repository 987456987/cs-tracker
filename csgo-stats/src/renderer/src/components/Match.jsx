import { useParams } from 'react-router-dom'

function Match({ data }) {
  const { matchId } = useParams()
  console.log(matchId)

  const currentMatch = data.find((match) => match.match_info.match_id === matchId)

  const playerStats = currentMatch.match_stats

  const teamA = []
  const teamB = []
  playerStats.forEach((player) => {
    if (player.team_id == currentMatch.team_info[0].team_id) {
      teamA.push(player)
    }
    if (player.team_id == currentMatch.team_info[1].team_id) {
      teamB.push(player)
    }
  })

  console.log(currentMatch)
  if (currentMatch) console.log(currentMatch.match_stats)
  return currentMatch ? (
    <>
      <div className="match-list">
        <table>
          <thead>
            <tr>
              <th>Player</th>
              <th>Kills</th>
              <th>Assists</th>
              <th>Deaths</th>
              <th>ADR</th>
              <th>Counter Strafe</th>
            </tr>
          </thead>
          <tbody>
            {teamA.map((player) => (
              <tr key={player.player_id}>
                <td>{player.nickname}</td>
                <td>{player.kills}</td>
                <td>{player.assists}</td>
                <td>{player.deaths}</td>
                <td>{player.adr}</td>
                <td>{player.counter_strafing}</td>
              </tr>
            ))}
            {teamB.map((player) => (
              <tr key={player.player_id}>
                <td>{player.nickname}</td>
                <td>{player.kills}</td>
                <td>{player.assists}</td>
                <td>{player.deaths}</td>
                <td>{player.adr}</td>
                <td>{player.counter_strafing}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </>
  ) : (
    <></>
  )
}

export default Match
