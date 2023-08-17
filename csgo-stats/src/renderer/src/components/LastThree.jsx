/* eslint-disable prettier/prettier */
/* eslint-disable react/prop-types */
import MatchInfo from "./MatchInfo"

function LastThree({ data, userId }) {
  return (
    <div className="last-three">
      <MatchInfo data={data} userId={userId} index={1} />
      <MatchInfo data={data} userId={userId} index={2} />
      <MatchInfo data={data} userId={userId} index={3} />
    </div>
  )
}

export default LastThree

