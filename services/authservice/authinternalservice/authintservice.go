package authinternalservice

import (
	"datcha/datamodel"
	"datcha/repository/authrepository"
	"datcha/servercommon"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	INT_AUTH_PATH = "$.auth.internal"
)

const (
	HEADER_CONTENT_TYPE          string = "Content-Type"
	APPLICATION_JSON_TYPE        string = "application/json"
	MIN_PASSWORD_LENGHT          int    = 4
	MIN_USER_NAME_LEN            int    = 3
	ACCESS_COOCKIE_NAME          string = "token"
	DELETE_COOKIE_AGE            int    = -100
	ACCESS_SUBJECT               string = "access"
	REFRESH_SUBJECT              string = "refresh"
	CONFIRM_SUBJECT              string = "confirm"
	RESET_SUBJECT                string = "reset"
	BEARER                       string = "Bearer"
	AUTHORIZATION_HEADER         string = "Authorization"
	ACCESS_CONTROL_EXPOSE_HEADER string = "Access-Control-Expose-Headers"
	INVALID_USER_ID              uint   = 0
)

type AuthInternalService struct {
	repository           authrepository.AuthRepositorier
	AccessSecretKey      string `json:"access_secret" env:"{SERVER_NAME}_AUTH_ACCESS_SECRET" default:"qx04ZY06cB%kX%#mPLq@qEFaBdmP7kqu"`
	RefreshSecretKey     string `json:"refresh_secret" env:"{SERVER_NAME}_REFRESH_SECRET" default:"u/5-@TYm*yG4p0iw8S8;FqHG6z]hE=@*"`
	ConfirmSecretKey     string `json:"confirm_secret" env:"{SERVER_NAME}_CONFIRM_SECRET" default:"LH/*tX3#tYK/yxV6)kE48ghQ1J#1NL}@"`
	Issuer               string `json:"issuer" env:"{SERVER_NAME}_AUTH_ISSUER" default:"datchasmarthome"`
	AcessTokenAge        int    `json:"acceess_token_age" env:"{SERVER_NAME}_AUTH_ACCESS_TOKEN_AGE" default:"600"`
	RefreshTokenAge      int    `json:"refresh_token_age" env:"{SERVER_NAME}_AUTH_REFRESH_TOKEN_AGE" default:"5184000"`
	ConfirmTokenAge      int    `json:"confirm_token_age" env:"{SERVER_NAME}_AUTH_CONFIRM_TOKEN_AGE" default:"600"`
	ResetTokenAge        int    `json:"reset_token_age" env:"{SERVER_NAME}_AUTH_RESET_TOKEN_AGE" default:"600"`
	CookieAge            int
	ConfiramationEmail   string `json:"confirmation_email" env:"{SERVER_NAME}_AUTH_CONFIRMATION_EMAIL"`
	EmailPassword        string `json:"email_password" env:"{SERVER_NAME}_EMAIL_PASSOWRD"`
	SMTPServerURL        string `json:"smtp_server_url" env:"{SERVER_NAME}_SMTP_SERVER_URL"`
	SMTPServerPort       int    `json:"smtp_server_port" env:"{SERVER_NAME}_SMTP_SERVER_PORT"`
	EmailTemplatesFolder string `json:"email_templates_folder" env:"{SERVER_NAME}_EMAIL_TEMPLATES_FOLDER"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type StatusResponse struct {
	AuthStatus datamodel.AccountStatus `json:"auth_status"`
}

type UserData struct {
	UserName string `json:"username" schema:"name"`
	Password string `json:"password" schema:"psw"`
	Email    string `json:"email"  schema:"email"`
	Locale   string `json:"locale"  schema:"locale"`
}

type JwtAuthClaims struct {
	jwt.RegisteredClaims
	UserId uint `json:"user_id"`
}

func (claim JwtAuthClaims) GetUserId() (uint, error) {
	return claim.UserId, nil
}

func NewAuthServer(cfgReader *servercommon.ConfigurationReader, rep authrepository.AuthRepositorier) (*AuthInternalService, error) {
	server := AuthInternalService{}
	server.repository = rep
	err := cfgReader.ReadConfiguration(&server, INT_AUTH_PATH)
	if err != nil {
		return nil, err
	}
	return &server, nil
}

func (server *AuthInternalService) validateUserName(userName string) error {
	if userName == "" {
		log.Println("Error. User name is empty")
		return errors.New(servercommon.ERROR_NAME_EMPTY)
	}
	if len(userName) < MIN_USER_NAME_LEN {
		log.Println("Error. User name too short")
		return errors.New(servercommon.ERROR_NAME_TOO_SHORT)
	}
	return nil
}

func (server *AuthInternalService) getUserData(r *http.Request) (UserData, error) {
	userData := UserData{}
	err := servercommon.ProcessBodyData(r, &userData)
	return userData, err
}

func (server *AuthInternalService) generateToken(userId uint, lifeTime int,
	secret string, subject string) (string, error) {
	expTime := jwt.NumericDate{time.Now().Add(time.Duration(lifeTime) * time.Second)}
	claims := JwtAuthClaims{
		jwt.RegisteredClaims{
			Subject:   subject,
			Issuer:    server.Issuer,
			ExpiresAt: &expTime,
			IssuedAt:  &jwt.NumericDate{time.Now()},
		},
		userId,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Println("Error. Generation token failed. Error: " + err.Error())
		return "", errors.New(servercommon.ERROR_INTERNAL)
	}
	return tokenStr, nil
}

func (server *AuthInternalService) generateAuthCockie(userId uint) (*http.Cookie, error) {
	tokenStr, err := server.generateToken(userId, server.AcessTokenAge,
		server.AccessSecretKey, ACCESS_SUBJECT)
	if err != nil {
		log.Println("error parse login data. Error: " + err.Error())
		return nil, errors.New(servercommon.ERROR_INTERNAL)
	}
	cookie := http.Cookie{
		Name:     ACCESS_COOCKIE_NAME,
		Value:    tokenStr,
		HttpOnly: true,
		//Temporay. Need to work with http during debug. Must be true
		Secure:   false,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(server.CookieAge),
	}
	return &cookie, nil
}

func (server *AuthInternalService) SendAuthCoockie(w http.ResponseWriter, userId uint) {
	cookie, err := server.generateAuthCockie(userId)
	if err != nil {
		http.Error(w, servercommon.ERROR_INTERNAL, http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, cookie)
}

func (server *AuthInternalService) getUserFromRequest(r *http.Request) (*datamodel.User, error) {
	userId, err := server.getUserIdFromRequest(r)
	if err != nil {
		log.Println("Attemt of non autorised access. Error: " + err.Error())
		return &datamodel.User{}, err
	}
	user, err := server.repository.FindUser(userId)
	if err != nil {
		log.Println("Attemt of non autorised access. Error: " + err.Error())
		return &datamodel.User{}, err
	}
	return user, nil
}

func (server *AuthInternalService) getUserIdFromRequest(r *http.Request) (uint, error) {
	authCookie, err := r.Cookie(ACCESS_COOCKIE_NAME)
	if err != nil {
		return INVALID_USER_ID, errors.New(servercommon.ERROR_NOT_AUTHORISED)
	}
	return server.getUserIdFromToken(authCookie.Value, ACCESS_SUBJECT, server.AccessSecretKey)
}

func (server *AuthInternalService) getUserIdFromToken(tokenStr string, subjectStr string, secret string) (uint, error) {
	claims := JwtAuthClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return INVALID_USER_ID, errors.New(servercommon.ERROR_NOT_AUTHORISED)
	}
	if !token.Valid {
		return INVALID_USER_ID, errors.New(servercommon.ERROR_NOT_AUTHORISED)
	}
	subject, err := claims.GetSubject()
	if err != nil {
		return INVALID_USER_ID, errors.New(servercommon.ERROR_NOT_AUTHORISED)
	}
	if subject != subjectStr {
		return INVALID_USER_ID, errors.New(servercommon.ERROR_NOT_AUTHORISED)
	}
	userId, err := claims.GetUserId()
	if err != nil {
		return INVALID_USER_ID, errors.New(servercommon.ERROR_NOT_AUTHORISED)
	}
	tokenExpTime, err := claims.GetExpirationTime()
	if err != nil {
		log.Println("can't get token time. Error: " + err.Error())
		return INVALID_USER_ID, errors.New(servercommon.ERROR_NOT_AUTHORISED)
	}
	currTime := time.Now()
	if currTime.After(tokenExpTime.Time) {
		log.Println("token is expired")
		return INVALID_USER_ID, errors.New(servercommon.ERROR_NOT_AUTHORISED)
	}
	fmt.Println("Claim. UserId =" + strconv.FormatUint(uint64(userId), 10))
	return userId, nil
}

func (server *AuthInternalService) WriteAutorizationHeader(w http.ResponseWriter, userId uint) error {
	tokenData, err := server.generateToken(userId, server.RefreshTokenAge,
		server.RefreshSecretKey, REFRESH_SUBJECT)
	if err != nil {
		log.Println("Error. Can't generate refresh token. Error=" + err.Error())
		return errors.New(servercommon.ERROR_INTERNAL)
	}
	tokenStr := BEARER + " " + tokenData
	w.Header().Set(AUTHORIZATION_HEADER, tokenStr)
	w.Header().Set(ACCESS_CONTROL_EXPOSE_HEADER, AUTHORIZATION_HEADER)
	return nil
}

func (server *AuthInternalService) writeStatusResponse(w http.ResponseWriter, status datamodel.AccountStatus) error {
	resp := StatusResponse{
		AuthStatus: status,
	}
	return servercommon.SendJsonResponse(w, resp)
}

func (server *AuthInternalService) RegisterHandlers(r *http.ServeMux) {
	r.HandleFunc("PUT /login", server.loginPutHandle)
	r.HandleFunc("POST /logout", server.logoutPostHandle)
	r.HandleFunc("POST /register", server.registerPostHandle)
	r.HandleFunc("PUT /refresh", server.refreshPostHandle)
	r.HandleFunc("GET /confirm/{"+servercommon.CONFIRM_TOKEN_KEY+"}", server.confirmGetHandle)
	r.HandleFunc("PUT /confirm", server.confirmPutHandle)
	r.HandleFunc("GET /auth_status", server.statusGetHandle)
	r.HandleFunc("GET /refresh", server.refreshGetHandle)
}
