package users

var settings *Settings

type Settings struct {
	SetActiveAfterUserCreation bool `json:"set_active_after_user_creation" yaml:"set_active_after_user_creation"`
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
