package authinternalservice

import "net/http"

func (server *AuthInternalService) logoutPostHandle(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(ACCESS_COOCKIE_NAME)
	if err != nil {
		return
	}
	cookie.MaxAge = DELETE_COOKIE_AGE
	http.SetCookie(w, cookie)
}
