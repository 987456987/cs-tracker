import { useParams } from 'react-router-dom'

function Match() {
  const { matchId } = useParams()
  return (
    <>
      <div>{matchId}</div>
    </>
  )
}

export default Match
