import React, { useState, useEffect } from "react"
import { Endpoints } from "../api"
import { deleteCookie } from "../utils"
import Errors from "../components/Errors"

const Session = ({ history }) => {
  const [user, setUser] = useState(null)
  const [isFetching, setIsFetching] = useState(false)
  const [errors, setErrors] = useState([])

  const headers = {
    Accept: "application/json",
    Authorization: document.cookie.split("token=")[1],
  }

  const getUserInfo = async () => {
    try {
      setIsFetching(true)
      const res = await fetch(Endpoints.session, {
        method: "GET",
        credentials: "include",
        headers,
      })

      if (!res.ok) logout()

      const { success, errors = [], user } = await res.json()
      setErrors(errors)
      if (!success) history.push("/login")
      setUser(user)
    } catch (e) {
      setErrors([e.toString()])
    } finally {
      setIsFetching(false)
    }
  }

  const logout = async () => {
    const res = await fetch(Endpoints.logout, {
      method: "GET",
      credentials: "include",
      headers,
    })

    if (res.ok) {
      deleteCookie("token")
      history.push("/login")
    }
  }

  useEffect(() => {
    getUserInfo()
  }, [])

  return (
    <div className="wrapper">
      <div>
        {isFetching ? (
          <div>fetching details...</div>
        ) : (
          <div>
            {user && (
              <div>
                <h1>Welcome, {user && user.name}</h1>
                <p>{user && user.email}</p>
                <br />
                <button onClick={logout}>logout</button>
              </div>
            )}
          </div>
        )}

        <Errors errors={errors} />
      </div>
    </div>
  )
}

export default Session
