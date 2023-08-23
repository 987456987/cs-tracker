import { AreaChart, Area, CartesianGrid, XAxis, YAxis, Tooltip } from 'recharts'

function UserStats({ data, userId }) {
  const adrData = data.map((match) => ({
    ADR: match.match_stats.find((player) => player.player_id === userId).adr,
    date: convertUnixToYMD(match.match_info.finished)
  }))

  function CustomTooltip({ active, payload }) {
    if (active && payload && payload.length) {
      const dataPoint = payload[0].payload
      return (
        <div className="custom-tooltip">
          <p>{`Date: ${dataPoint.date}`}</p>
          <p>{`ADR: ${dataPoint.ADR}`}</p>
        </div>
      )
    }
    return null
  }

  function renderChart() {
    return (
      <div className="chart-container">
        <AreaChart width={600} height={300} data={adrData}>
          <Area type="monotone" dataKey="ADR" stroke="#7096b3" fill="#7096b3" />
          <CartesianGrid stroke="gray" strokeDasharray="3 3" />
          <XAxis dataKey="date" hide={true} />
          <YAxis
            domain={[
              Math.min(...adrData.map((entry) => entry.ADR)) - 10,
              Math.max(...adrData.map((entry) => entry.ADR)) + 10
            ]}
          />
          <Tooltip content={<CustomTooltip />} />
        </AreaChart>
      </div>
    )
  }

  return <>{renderChart()}</>
}

function convertUnixToYMD(unixTime) {
  const date = new Date(unixTime * 1000)
  const year = date.getFullYear()
  const month = date.getMonth() + 1
  const day = date.getDate()

  const formattedDate = `${year}-${month.toString().padStart(2, '0')}-${day
    .toString()
    .padStart(2, '0')}`
  return formattedDate
}

export default UserStats
