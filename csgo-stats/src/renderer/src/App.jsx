/* eslint-disable prettier/prettier */
import { useState, useEffect } from 'react'
import LastThree from './components/LastThree'
import UserIdInput from './components/UserIdInput'

function formatDuration(durationMs) {
  if (durationMs < 1000) {
    return `${durationMs} ms`
  } else if (durationMs < 60000) {
    return `${(durationMs / 1000).toFixed(2)} seconds`
  } else {
    return `${(durationMs / 60000).toFixed(2)} minutes`
  }
}

function getData(requestData, setData, setFetchTime) {
  // Request a wake lock to prevent the computer from sleeping
  let wakeLock = null

  if ('wakeLock' in navigator) {
    navigator.wakeLock
      .request('screen')
      .then((lock) => {
        wakeLock = lock
        console.log('Locking Screen:')
      })
      .catch((error) => {
        console.error('Wake Lock request failed:', error)
      })
  }

  const startTime = performance.now()

  fetch('http://localhost:8080/get-matches', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(requestData)
  })
    .then((response) => response.json())
    .then((data) => {
      const endTime = performance.now()
      setFetchTime(endTime - startTime)
      setData(data.matches.sort((a, b) => b.match_info.finished - a.match_info.finished))

      // Release the wake lock
      if (wakeLock) {
        wakeLock.release()
        console.log('Releasing Screen:')
      }
    })
    .catch((error) => {
      console.error('Error fetching data:', error)
      setFetchTime(null)

      // Release the wake lock
      if (wakeLock) {
        wakeLock.release()
      }
    })
}

function checkData(requestData) {
  fetch('http://localhost:8080/check-matches', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(requestData)
  })
    .catch((error) => {
      console.error('Error fetching data:', error)
    })
}

function getUserID(userName) {
  const apiUrl = `https://open.faceit.com/data/v4/players?nickname=${userName}`
  const token = 'ca431399-60b3-4381-ae55-12b5657fde76'

  return fetch(apiUrl, {
    method: 'GET',
    headers: {
      accept: 'application/json',
      Authorization: `Bearer ${token}`
    }
  })
    .then((response) => response.json())
    .then((data) => {
      console.log(data)
      return data.player_id // Return the player ID
    })
    .catch((error) => {
      console.error('Error fetching data:', error)
      return null // Return null in case of an error
    })
}

function App() {
  const [data, setData] = useState([])
  const [fetchTime, setFetchTime] = useState(null)
  const [userId, setUserId] = useState('')

  useEffect(() => {
    const storedUserId = localStorage.getItem('user_id')
    if (storedUserId) {
      setUserId(storedUserId)
      const requestData = {
        user_id: storedUserId // Use the playerID obtained from getUserID
      }
      checkData(requestData)
      getData(requestData, setData, setFetchTime)
    }
  }, [])

  const handleFetchData = async (event) => {
    event.preventDefault()
    if (inputValue) {
      const playerID = await getUserID(inputValue)
      if (playerID) {
        localStorage.setItem('user_id', playerID)
        setUserId(playerID)

        const requestData = {
          user_id: playerID // Use the playerID obtained from getUserID
        }
        checkData(requestData)
        getData(requestData, setData, setFetchTime)
      }
    }
  }

  console.log(data)
  console.log('Fetch time:', fetchTime ? formatDuration(fetchTime) : 'N/A')

  return (
    <div className="App">
      {userId ? (
        <>
          <header className="App-header">
          </header>
          <LastThree data={data} userId={userId} />
        </>
      ) : (
        <UserIdInput handleFetchData={handleFetchData} />
      )}
    </div>
  )
}

export default App

//52e19c3d-471c-44ac-af67-a394da815a37
