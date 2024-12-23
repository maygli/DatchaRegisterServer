package auththirdparty

import (
	"datcha/servercommon"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
)

type VkUser struct {
	UserId    string `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Avatar    string `json:"avatar"`
	Email     string `json:"email"`
}

type VkData struct {
	User             VkUser `json:"user"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	State            string `json:"state"`
}

type VkAuthConfigure struct {
	VkClientId       string `json:"vk_client_id" env:"${SERVER_NAME}_AUTH_VK_CLIENT_ID"`
	VkClientSecret   string `json:"vk_client_secret" env:"${SERVER_NAME}_AUTH_VK_CLIENT_SECRET"`
	ServerAddress    string `json:"server_address" env:"${SERVER_NAME}_SERVER_ADDRESS"`
	LoginEndPoint    string `json:"vk_login_endpoint" env:"${SERVER_NAME}_VK_LOGIN_ENDPOINT" default:"/auth/vk/login"`
	CallbackEndPoint string `json:"vk_callback_endpoint" env:"${SERVER_NAME}_VK_CALLBACK_ENDPOINT" default:"/auth/vk/callback"`
}

type VkAuthService struct {
	AuthThirdPartyBase
	vkLoginConfig oauth2.Config
}

const (
	VK_EMAIL_SCOPE           = "email"
	VK_PROFILE_SCOPE         = "vkid.personal_info"
	VK_AUTH_URL              = "https://id.vk.com/authorize"
	VK_TOKEN_URL             = "https://id.vk.com/oauth2/auth"
	VK_USER_PROFILE_DATA_URL = "https://id.vk.com/oauth2/user_info"
	VK_GET_USER_METHOD       = http.MethodPost
	VK_USER_KEY              = "user"
)

func NewVkAuthService(config *VkAuthConfigure, processor CompleteAuthProcessor) *VkAuthService {
	service := VkAuthService{
		AuthThirdPartyBase: AuthThirdPartyBase{
			LoginEndPoint:    config.LoginEndPoint,
			CallbackEndPoint: config.CallbackEndPoint,
			AuthProcessor:    processor,
		},
	}
	service.vkLoginConfig = oauth2.Config{
		RedirectURL: config.ServerAddress + service.GetCallbackEndpoint(service.GetServiceName()),
		ClientID:    config.VkClientId,
		Scopes:      []string{VK_EMAIL_SCOPE, VK_PROFILE_SCOPE},
		// Oauth2 config for VK is incorrect - so we have set endpoints by 'hand'
		Endpoint: oauth2.Endpoint{
			AuthURL:  VK_AUTH_URL,
			TokenURL: VK_TOKEN_URL,
		},
	}
	return &service
}

func (service VkAuthService) GetServiceName() string {
	return servercommon.VK_SERVICE_NAME
}

func (service VkAuthService) generateVerifier() string {
	return GenerateRandonString(VERIFIER_LENGHT)
}

func (service VkAuthService) vkAuthLoginHandle(w http.ResponseWriter, r *http.Request) {
	verifier := service.generateVerifier()
	slog.Info(fmt.Sprintf("Code verifier=%s", verifier))
	http.SetCookie(w, GenerateVerifierCockie(verifier))
	codeChallenge := oauth2.S256ChallengeFromVerifier(verifier)
	//	redirect_uri := service.ServerAddress + AUTH_VK_CALLBACK_ENDPOINT
	slog.Info(fmt.Sprintf("Code challenge=%s", codeChallenge))
	state := GenerateRandonString(STATE_LENGTH)
	http.SetCookie(w, GenerateStateCockie(state))
	authUrl := service.vkLoginConfig.AuthCodeURL(state,
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		oauth2.SetAuthURLParam("code_challenge", codeChallenge))
	http.Redirect(w, r, authUrl, http.StatusSeeOther)
}

func vkToAuthUser(user VkUser) AuthThirdPartyUser {
	resUser := AuthThirdPartyUser{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Avatar:    user.Avatar,
		UserId:    user.UserId,
		Service:   servercommon.VK_SERVICE_NAME,
	}
	return resUser
}

func (service VkAuthService) vkAuthCallbackHandle(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get(STATE_NAME)
	if state == "" {
		err := errors.New("Get vk oauth callback without state")
		service.ProcessError(err, w, r)
		return
	}
	err := VerifyCookieValue(r, STATE_COOKIE_NAME, state)
	if err != nil {
		slog.Error(fmt.Sprintf("Incorrect state value. Error: %s", err.Error()))
		service.ProcessError(err, w, r)
		return
	}
	code := r.URL.Query().Get("code")
	if code == "" {
		err = errors.New("Get vk outh callback without code")
		slog.Error(err.Error())
		service.ProcessError(err, w, r)
		return
	}
	deviceId := r.URL.Query().Get("device_id")
	if code == "" {
		err = errors.New("Get vk outh callback without device id")
		slog.Error(err.Error())
		service.ProcessError(err, w, r)
		return
	}
	verifier, err := r.Cookie(VERIFIER_COOKIE_NAME)
	if err != nil {
		slog.Error(fmt.Sprintf("Vk auth. Error get verifier from cookie. Error: %s", err.Error()))
		service.ProcessError(err, w, r)
		return
	}
	// Vk also send deviceId. Now we don't need it
	slog.Info(fmt.Sprintf("Code=%s", code))
	token, err := service.vkLoginConfig.Exchange(r.Context(), code,
		oauth2.SetAuthURLParam("code_verifier", verifier.Value),
		oauth2.SetAuthURLParam("device_id", deviceId))
	if err != nil {
		slog.Error(fmt.Sprintf("Get vk outh token failed. Error: %s", err.Error()))
		service.ProcessError(err, w, r)
	}
	slog.Info(fmt.Sprintf("Get token: %+v", token))
	user, err := service.getUserDataFromVk(token)
	if err != nil {
		slog.Error(fmt.Sprintf("Get vk user data failed. Error: %s", err.Error()))
		service.ProcessError(err, w, r)
	}
	slog.Info(fmt.Sprintf("User = %v", user))
	vkUser := vkToAuthUser(user)
	err = service.ProcessSuccess(vkUser, w, r)
	if err != nil {
		slog.Error(err.Error())
	}
}

func (service VkAuthService) getUserDataFromVk(token *oauth2.Token) (VkUser, error) {
	data := url.Values{}
	data.Set(servercommon.ACCESS_TOKEN_KEY, token.AccessToken)
	data.Set(servercommon.CLIENT_ID_KEY, service.vkLoginConfig.ClientID)
	vkData := VkData{}
	err := servercommon.JsonRequest(VK_GET_USER_METHOD, VK_USER_PROFILE_DATA_URL, data, &vkData)
	if err != nil {
		return VkUser{}, err
	}
	if vkData.Error != "" {
		slog.Error(fmt.Sprintf("error read Vk user data.Error: %s", vkData.ErrorDescription))
		return VkUser{}, errors.New(vkData.Error)
	}
	return vkData.User, err
}

func (service VkAuthService) RegisterHandlers(mux *http.ServeMux) {
	mux.HandleFunc("GET "+service.GetLoginEndpoint(service.GetServiceName()), service.vkAuthLoginHandle)
	mux.HandleFunc("GET "+service.GetCallbackEndpoint(service.GetServiceName()), service.vkAuthCallbackHandle)
}
