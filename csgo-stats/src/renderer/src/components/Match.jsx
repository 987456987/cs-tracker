import { useParams } from 'react-router-dom'

function Match({ data }) {
  const { matchId } = useParams()
  console.log(matchId)

  const currentMatch = data.find((match) => match.match_info.match_id === matchId)

  console.log(currentMatch)
  return currentMatch ? (
    <>
      <div className="scoreboard">
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
            {currentMatch.match_stats.map((player) => {
              if (player.team_id === currentMatch.team_info[0].team_id) {
                return (
                  <tr key={player.player_id}>
                    <td>{player.nickname}</td>
                    <td>{player.kills}</td>
                    <td>{player.assists}</td>
                    <td>{player.deaths}</td>
                    <td>{player.adr}</td>
                    <td>{player.counter_strafing}</td>
                  </tr>
                )
              }
              return null
            })}
          </tbody>
        </table>
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
            {currentMatch.match_stats.map((player) => {
              if (player.team_id === currentMatch.team_info[1].team_id) {
                return (
                  <tr key={player.player_id}>
                    <td>{player.nickname}</td>
                    <td>{player.kills}</td>
                    <td>{player.assists}</td>
                    <td>{player.deaths}</td>
                    <td>{player.adr}</td>
                    <td>{player.counter_strafing}</td>
                  </tr>
                )
              }
              return null
            })}
          </tbody>
        </table>
      </div>
    </>
  ) : (
    <></>
  )
}

export default Match
