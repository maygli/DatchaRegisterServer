package authentification

import "net/http"

func (server *AuthServer) logoutPostHandle(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(ACCESS_COOCKIE_NAME)
	if err != nil {
		return
	}
	cookie.MaxAge = DELETE_COOKIE_AGE
	http.SetCookie(w, cookie)
}
