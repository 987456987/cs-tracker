import { LineChart, Line, CartesianGrid, XAxis, YAxis, Tooltip } from 'recharts';

function UserStats({ data, userId }) {
  const adrData = data.map((match) => ({
    ADR: match.match_stats.find((player) => player.player_id === userId).adr
  }));

  function renderChart() {
    return (
      <LineChart width={600} height={300} data={adrData}>
        <Line type="monotone" dataKey="ADR" stroke="#8884d8" />
        <CartesianGrid stroke="#ccc" />
        <YAxis />
        <Tooltip label={"ADR"}/>
      </LineChart>
    );
  }

  return <>{renderChart()}</>;
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

export default UserStats;