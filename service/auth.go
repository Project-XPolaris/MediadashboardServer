package service

import (
	"fmt"
	"mediadashboard/database"
	"mediadashboard/plugins"

	"github.com/allentom/harukap/plugins/youauth"
	"github.com/rs/xid"
)

const YouAuthProvider = "youauth"

func GenerateYouAuthToken(code string) (string, string, error) {
	tokens, err := plugins.DefaultYouAuthOauthPlugin.Client.GetAccessToken(code)
	if err != nil {
		return "", "", err
	}
	return LinkWithYouAuthToken(tokens)
}
func GenerateYouAuthTokenByPassword(username string, rawPassword string) (string, string, error) {
	authResult, err := plugins.DefaultYouAuthOauthPlugin.Client.GrantWithPassword(username, rawPassword)
	if err != nil {
		return "", "", err
	}
	return LinkWithYouAuthToken(authResult)
}
func LinkWithYouAuthToken(tokens *youauth.GenerateTokenResponse) (string, string, error) {
	currentUserResponse, err := plugins.DefaultYouAuthOauthPlugin.Client.GetCurrentUser(tokens.AccessToken)
	if err != nil {
		return "", "", err
	}
	// check if user exists
	uid := fmt.Sprintf("%d", currentUserResponse.Id)
	historyOauth := make([]database.Oauth, 0)
	err = database.Instance.Where("uid = ?", uid).
		Where("provider = ?", YouAuthProvider).
		Preload("User").
		Find(&historyOauth).Error
	if err != nil {
		return "", "", err
	}
	var user *database.User
	if len(historyOauth) == 0 {
		username := xid.New().String()
		// create new user
		user = &database.User{
			Uid:      xid.New().String(),
			Username: username,
		}
		err = database.Instance.Create(&user).Error
		if err != nil {
			return "", "", err
		}
	} else {
		user = historyOauth[0].User
	}

	oauthRecord := database.Oauth{
		Uid:          fmt.Sprintf("%d", currentUserResponse.Id),
		UserId:       user.ID,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		Provider:     YouAuthProvider,
	}
	err = database.Instance.Create(&oauthRecord).Error
	if err != nil {
		return "", "", err
	}
	return tokens.AccessToken, currentUserResponse.Username, nil
}
func refreshToken(accessToken string) (string, error) {
	tokenRecord := database.Oauth{}
	err := database.Instance.Where("access_token = ?", accessToken).First(&tokenRecord).Error
	if err != nil {
		return "", err
	}
	token, err := plugins.DefaultYouAuthOauthPlugin.Client.RefreshAccessToken(tokenRecord.RefreshToken)
	if err != nil {
		return "", err
	}
	err = database.Instance.Delete(&tokenRecord).Error
	if err != nil {
		return "", err
	}
	newOauthRecord := database.Oauth{
		UserId:       tokenRecord.UserId,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}
	err = database.Instance.Create(&newOauthRecord).Error
	if err != nil {
		return "", err
	}
	return token.AccessToken, nil
}

func CheckToken() {

}
