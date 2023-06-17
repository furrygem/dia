package users

import "time"

var settings *Settings

type Settings struct {
	SetActiveAfterUserCreation   bool          `json:"set_active_after_user_creation" yaml:"set_active_after_user_creation"`
	AccessTokenExpirationTimeout time.Duration `json:"access_token_expiration_timeout" yaml:"access_token_expiration_timeout"`
}

func NewSettings() *Settings {
	if settings != nil {
		return settings
	}
	settings = &Settings{
		SetActiveAfterUserCreation: true,
	}
	return settings
}

func getSettings() *Settings {
	return settings
}
