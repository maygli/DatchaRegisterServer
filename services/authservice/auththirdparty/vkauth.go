package auththirdparty

import (
	"datcha/servercommon"
	"datcha/services/authservice/authcommon"
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

type VkAuthService struct {
	VkClientId     string `json:"vk_client_id" env:"${SERVER_NAME}_AUTH_VK_CLIENT_ID"`
	VkClientSecret string `json:"vk_client_secret" env:"${SERVER_NAME}_AUTH_VK_CLIENT_SECRET"`
	ServerAddress  string `json:"server_address" env:"${SERVER_NAME}_SERVER_ADDRESS"`
	vkLoginConfig  oauth2.Config
}

const (
	AUTH_VK_CONFIGURATION_PATH = "$.auth.vk"
	VK_REDIRECT_URL            = "https://id.vk.com/authorize?response_type=code&client_id=%s&code_challenge=%s&code_challenge_method=s256&redirect_uri=%s&state=%s&scope=email vk_id.personal_info"
	AUTH_VK_CALLBACK_ENDPOINT  = "/auth/vk/callback"
	VK_EMAIL_SCOPE             = "email"
	VK_PROFILE_SCOPE           = "vkid.personal_info"
	VK_AUTH_URL                = "https://id.vk.com/authorize"
	VK_TOKEN_URL               = "https://id.vk.com/oauth2/auth"
	VK_USER_PROFILE_DATA_URL   = "https://id.vk.com/oauth2/user_info"
	VK_GET_USER_METHOD         = http.MethodPost
	VK_USER_KEY                = "user"
)

func NewVkAuthService(cfgReader *servercommon.ConfigurationReader) (*VkAuthService, error) {
	service := VkAuthService{}
	err := cfgReader.ReadConfiguration(&service, AUTH_VK_CONFIGURATION_PATH)
	if err != nil {
		return &service, err
	}
	service.vkLoginConfig = oauth2.Config{
		RedirectURL: service.ServerAddress + AUTH_VK_CALLBACK_ENDPOINT,
		ClientID:    service.VkClientId,
		Scopes:      []string{VK_EMAIL_SCOPE, VK_PROFILE_SCOPE},
		// Oauth2 config for VK is incorrect - so we have set endpoints by 'hand'
		Endpoint: oauth2.Endpoint{
			AuthURL:  VK_AUTH_URL,
			TokenURL: VK_TOKEN_URL,
		},
	}
	return &service, nil
}

func (service VkAuthService) generateVerifier() string {
	return authcommon.GenerateRandonString(authcommon.VERIFIER_LENGHT)
}

func (service VkAuthService) vkAuthLoginHandle(w http.ResponseWriter, r *http.Request) {
	verifier := service.generateVerifier()
	slog.Info(fmt.Sprintf("Code verifier=%s", verifier))
	http.SetCookie(w, authcommon.GenerateVerifierCockie(verifier))
	codeChallenge := oauth2.S256ChallengeFromVerifier(verifier)
	//	redirect_uri := service.ServerAddress + AUTH_VK_CALLBACK_ENDPOINT
	slog.Info(fmt.Sprintf("Code challenge=%s", codeChallenge))
	state := authcommon.GenerateRandonString(authcommon.STATE_LENGTH)
	http.SetCookie(w, authcommon.GenerateStateCockie(state))
	authUrl := service.vkLoginConfig.AuthCodeURL(state,
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		oauth2.SetAuthURLParam("code_challenge", codeChallenge))
	//	url := fmt.Sprintf(VK_REDIRECT_URL, service.VkClientId, code_challenge, redirect_uri, state)
	http.Redirect(w, r, authUrl, http.StatusSeeOther)
}

func (service VkAuthService) vkAuthCallbackHandle(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get(authcommon.STATE_NAME)
	if state == "" {
		slog.Error("Get vk oauth callback without state")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	err := authcommon.VerifyCookieValue(r, authcommon.STATE_COOKIE_NAME, state)
	if err != nil {
		slog.Error(fmt.Sprintf("Incorrect state value. Error: %s", err.Error()))
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	code := r.URL.Query().Get("code")
	if code == "" {
		slog.Error("Get vk outh callback without code")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	deviceId := r.URL.Query().Get("device_id")
	if code == "" {
		slog.Error("Get vk outh callback without device id")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	verifier, err := r.Cookie(authcommon.VERIFIER_COOKIE_NAME)
	if err != nil {
		slog.Error(fmt.Sprintf("Vk auth. Error get verifier from cookie. Error: %s", err.Error()))
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	// Vk also send deviceId. Now we don't need it
	slog.Info(fmt.Sprintf("Code=%s", code))
	token, err := service.vkLoginConfig.Exchange(r.Context(), code,
		oauth2.SetAuthURLParam("code_verifier", verifier.Value),
		oauth2.SetAuthURLParam("device_id", deviceId))
	if err != nil {
		slog.Error(fmt.Sprintf("Get vk outh token failed. Error: %s", err.Error()))
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	slog.Info(fmt.Sprintf("Get token: %+v", token))
	user, err := service.getUserDataFromVk(token)
	if err != nil {
		slog.Error(fmt.Sprintf("Get vk user data failed. Error: %s", err.Error()))
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	slog.Info(fmt.Sprintf("User = %v", user))
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (service VkAuthService) getUserDataFromVk(token *oauth2.Token) (VkUser, error) {
	data := url.Values{}
	data.Set(servercommon.ACCESS_TOKEN_KEY, token.AccessToken)
	data.Set(servercommon.CLIENT_ID_KEY, service.VkClientId)
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
	mux.HandleFunc("GET /auth/vk/login", service.vkAuthLoginHandle)
	mux.HandleFunc("GET "+AUTH_VK_CALLBACK_ENDPOINT, service.vkAuthCallbackHandle)
}
