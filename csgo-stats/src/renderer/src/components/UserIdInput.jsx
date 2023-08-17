/* eslint-disable prettier/prettier */
/* eslint-disable react/prop-types */
function UserIdInput({ handleFetchData, inputValue, setInputValue }) {
  const handleUserIdChange = (event) => {
    setInputValue(event.target.value)
  }
  return (
    <>
      <div>
        <p>Please enter your user ID:</p>
        <form onSubmit={handleFetchData}>
          <input type="text" value={inputValue} onChange={handleUserIdChange} />
          <button type="submit">Fetch Data</button>
        </form>
      </div>
    </>
  )
}

export default UserIdInput
