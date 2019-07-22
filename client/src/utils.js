// create a named cookie
/**
 * @param {sring} cookieName
 * @param {string} cookieValue //token
 * @param {string} hourToExpire
 */
export const createCookie = (cookieName, cookieValue, hourToExpire) => {
  const date = new Date()
  date.setTime(date.getTime() + hourToExpire * 60 * 60 * 1000)
  document.cookie = `${cookieName} = ${cookieValue}; expires = ${date.toGMTString()}`
}

export const deleteCookie = (name) => {
  document.cookie = name + '=; expires=Thu, 01 Jan 1970 00:00:01 GMT;'
}
