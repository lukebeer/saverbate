package user

import (
	"regexp"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"github.com/volatiletech/authboss-renderer"
	"github.com/volatiletech/authboss/v3"
	"github.com/volatiletech/authboss/v3/defaults"
)

// InitAuthBoss configures and returns Authboss instance
func InitAuthBoss(db *sqlx.DB) (*authboss.Authboss, error) {
	ab := authboss.New()
	ab.Config.Paths.RootURL = viper.GetString("rootURL")

	userSessionStore := NewSessionStore()
	ab.Config.Storage.SessionState = userSessionStore.SessionStorer
	ab.Config.Storage.CookieState = userSessionStore.CookieStorer

	ab.Config.Storage.Server = NewStorer(db)

	ab.Config.Core.ViewRenderer = abrenderer.NewHTML("/auth", "web/templates/ab_views")
	ab.Config.Modules.RegisterPreserveFields = []string{"email", "name"}
	ab.Config.Modules.LogoutMethod = "GET"

	defaults.SetCore(&ab.Config, true, true)

	emailRule := defaults.Rules{
		FieldName: "email", Required: true,
		MatchError: "Must be a valid e-mail address",
		MustMatch:  regexp.MustCompile(`.*@.*\.[a-z]+`),
		MaxLength:  1024,
	}
	passwordRule := defaults.Rules{
		FieldName: "password", Required: true,
		MinLength: 4,
	}
	nameRule := defaults.Rules{
		FieldName: "name", Required: true,
		MinLength: 3,
		MaxLength: 255,
	}

	ab.Config.Core.BodyReader = defaults.HTTPBodyReader{
		ReadJSON: true,
		Rulesets: map[string][]defaults.Rules{
			"register": {emailRule, passwordRule, nameRule},
			//"recover_end": {passwordRule},
		},
		Confirms: map[string][]string{
			"register": {"password", authboss.ConfirmPrefix + "password"},
			//"recover_end": {"password", authboss.ConfirmPrefix + "password"},
		},
		Whitelist: map[string][]string{
			"register": {"email", "name", "password"},
		},
	}

	err := ab.Init()
	if err != nil {
		return nil, err
	}

	return ab, nil
}
