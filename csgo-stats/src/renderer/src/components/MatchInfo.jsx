/* eslint-disable prettier/prettier */
/* eslint-disable react/prop-types */
import '../assets/last_three.css'
import { Link } from 'react-router-dom'

function MatchInfo({ data, userId, index }) {
  const match = data[index]

  if (!match) {
    return <p>Loading...</p>
  }

  const mapName = match.match_info.map.slice(3)
  const userIndex = getUserIndex(match, userId)
  const userStats = match.match_stats[userIndex]
  const [teamScoreA, teamScoreB] =
    match.team_info[0].team_id === userStats.team_id
      ? [match.team_info[0].final_score, match.team_info[1].final_score]
      : [match.team_info[1].final_score, match.team_info[0].final_score]

  return (
    <div className="last-three-match">
      <Link to={`/match/${match.match_info.match_id}`} className="custom-link">
        <div className={`${match.match_info.map} last-three-match-top`}>
          <div className="last-three-info">
            <div className="last-three-top">
              <h2 className="capitalize">
                {mapName}
                <span className={parseInt(teamScoreA) > parseInt(teamScoreB) ? 'winner' : 'loser'}>
                  {' '}
                  {teamScoreA}
                </span>
                <span>/{teamScoreB}</span>
              </h2>
            </div>
            <div className="last-three-bottom">
              <h3>{convertUnixToYMD(match.match_info.finished)}</h3>
            </div>
          </div>
        </div>
        <div className="last-three-match-bottom">
          <span className="last-three-stat">{userStats.kills} Kills</span>
          <span className="last-three-stat">
            {(userStats.kills / userStats.deaths).toFixed(2)} K/D Ratio
          </span>
          <span className="last-three-stat">{userStats.adr} ADR</span>
        </div>
      </Link>
    </div>
  )
}

function convertUnixToYMD(unixTime) {
  const date = new Date(unixTime * 1000) // Convert Unix time to milliseconds
  const year = date.getFullYear()
  const month = date.getMonth() + 1 // Months are zero-based, so add 1
  const day = date.getDate()

  const formattedDate = `${year}-${month.toString().padStart(2, '0')}-${day
    .toString()
    .padStart(2, '0')}`
  return formattedDate
}

function getUserIndex(match, userId) {
  return match.match_stats.findIndex((player) => player.player_id === userId)
}

export default MatchInfo
