function MatchList({ data, userId }) {
  function getUserIndex(match) {
    return match.match_stats.findIndex((player) => player.player_id === userId)
  }

  function formatTimeAgo(unixTime) {
    const currentTime = new Date()
    const matchTime = new Date(unixTime * 1000) // Convert Unix time to milliseconds
    const timeDifference = currentTime - matchTime
    const minutesAgo = Math.floor(timeDifference / (1000 * 60))
    const hoursAgo = Math.floor(minutesAgo / 60)

    if (hoursAgo < 1) {
      return `${minutesAgo} minutes ago`
    } else if (hoursAgo < 24) {
      return `${hoursAgo} hours ago`
    } else if (hoursAgo < 48) {
      return `Yesterday`
    } else {
      return convertUnixToYMD(unixTime) // If more than 24 hours, use the original conversion
    }
  }

  return (
    <>
      <table border="1">
        <thead>
          <tr>
            <th>Map</th>
            <th>Date</th>
            <th>Kills</th>
            <th>Deaths</th>
            <th>ADR</th>
          </tr>
        </thead>
        <tbody>
          {data.map((match) => (
            <tr key={match.match_info.match_id}>
              <td>{match.match_info.map}</td>
              <td>{formatTimeAgo(match.match_info.finished)}</td>
              {getUserIndex(match) !== -1 ? (
                <>
                  <td>{match.match_stats[getUserIndex(match)].kills}</td>
                  <td>{match.match_stats[getUserIndex(match)].deaths}</td>
                  <td>{match.match_stats[getUserIndex(match)].adr}</td>
                </>
              ) : (
                <>
                  <td>N/A</td>
                  <td>N/A</td>
                  <td>N/A</td>
                </>
              )}
            </tr>
          ))}
        </tbody>
      </table>
    </>
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

export default MatchList
