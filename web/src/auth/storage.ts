const ACCESS_TOKEN = 'accessToken'

export const setAccessToken = (value: string) => {
  localStorage.setItem(ACCESS_TOKEN, value)
}
