/* eslint-disable prettier/prettier */
/* eslint-disable react/prop-types */
import '../assets/last_three.css';

function MatchInfo({ data, userId, index }) {
  const match = data[index];

  if (!match) {
    return <p>Loading...</p>;
  }

  const mapName = match.match_info.map.slice(3);
  const userIndex = getUserIndex(match, userId);
  const userStats = match.match_stats[userIndex];
  const [teamScoreA, teamScoreB] =
    match.team_info[0].team_id === userStats.team_id
      ? [match.team_info[0].final_score, match.team_info[1].final_score]
      : [match.team_info[1].final_score, match.team_info[0].final_score];

  return (
    <div className="last-three-match">
      <div className={`${match.match_info.map} last-three-match-top`}>
        <div className="last-three-info">
          <div className="last-three-top">
            <h2 className="capitalize">
              {mapName}
              <span
                className={
                  parseInt(teamScoreA) > parseInt(teamScoreB)
                    ? 'winner'
                    : 'loser'
                }
              >
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
        <span className="last-three-stat">{userStats.kd_ratio} K/D Ratio</span>
        <span className="last-three-stat">{userStats.adr} ADR</span>
      </div>
    </div>
  );
}

function convertUnixToYMD(unixTime) {
  const date = new Date(unixTime * 1000); // Convert Unix time to milliseconds
  const year = date.getFullYear();
  const month = date.getMonth() + 1; // Months are zero-based, so add 1
  const day = date.getDate();

  const formattedDate = `${year}-${month.toString().padStart(2, '0')}-${day.toString().padStart(2, '0')}`;
  return formattedDate;
}

function getUserIndex(match, userId) {
  for (let i = 0; i < match.match_stats.length; i++) {
    if (userId == match.match_stats[i].player_id) {
      return i;
    }
  }
}

export default MatchInfo;
