import { LineChart, Line, CartesianGrid, XAxis, YAxis, Tooltip } from 'recharts';

function UserStats({ data, userId }) {
  const adrData = data.map((match) => ({
    ADR: match.match_stats.find((player) => player.player_id === userId).adr,
    date: convertUnixToYMD(match.match_info.finished)
  }));

  function CustomTooltip({ active, payload, label }) {
    if (active && payload && payload.length) {
      const dataPoint = payload[0].payload;
      return (
        <div className="custom-tooltip">
          <p>{`Date: ${dataPoint.date}`}</p>
          <p>{`ADR: ${dataPoint.ADR}`}</p>
        </div>
      );
    }

    return null;
  }

  function renderChart() {
    return (
      <LineChart width={600} height={300} data={adrData}>
        <Line type="monotone" dataKey="ADR" stroke="#8884d8" />
        <CartesianGrid stroke="#ccc" />
        <XAxis dataKey="date" />
        <YAxis tickFormatter={(value) => convertUnixToYMD(value)} />
        <Tooltip content={<CustomTooltip />} />
      </LineChart>
    );
  }

  return <>{renderChart()}</>;
}

function convertUnixToYMD(unixTime) {
  const date = new Date(unixTime * 1000);
  const year = date.getFullYear();
  const month = date.getMonth() + 1;
  const day = date.getDate();

  const formattedDate = `${year}-${month.toString().padStart(2, '0')}-${day
    .toString()
    .padStart(2, '0')}`;
  return formattedDate;
}

export default UserStats;
