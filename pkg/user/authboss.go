package user

import (
	"regexp"
	"saverbate/pkg/handler"

	goredislib "github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"

	"github.com/spf13/viper"

	"github.com/volatiletech/authboss/v3"
	"github.com/volatiletech/authboss/v3/defaults"

	abrenderer "github.com/volatiletech/authboss-renderer"
)

// This pattern is useful in real code to ensure that
// we've got the right interfaces implemented.
var (
	assertUser   = &User{}
	assertStorer = &Storer{}

	_ authboss.User            = assertUser
	_ authboss.AuthableUser    = assertUser
	_ authboss.ConfirmableUser = assertUser
	// _ authboss.LockableUser    = assertUser
	// _ authboss.RecoverableUser = assertUser
	_ authboss.ArbitraryUser = assertUser

	_ authboss.CreatingServerStorer   = assertStorer
	_ authboss.ConfirmingServerStorer = assertStorer
	// _ authboss.RecoveringServerStorer  = assertStorer
	_ authboss.RememberingServerStorer = assertStorer
)

// InitAuthBoss configures and returns Authboss instance
func InitAuthBoss(db *sqlx.DB, redis *goredislib.Client) (*authboss.Authboss, error) {
	ab := authboss.New()
	ab.Config.Paths.RootURL = viper.GetString("rootURL")

	userSessionStore := NewSessionStore()
	ab.Config.Storage.SessionState = userSessionStore.SessionStorer
	ab.Config.Storage.CookieState = userSessionStore.CookieStorer

	ab.Config.Storage.Server = NewStorer(db, redis)

	ab.Config.Core.MailRenderer = abrenderer.NewEmail("/auth", "ab_views")

	ab.Config.Core.ViewRenderer = handler.NewHTML("/auth", "web/templates/ab_views")
	ab.Config.Modules.RegisterPreserveFields = []string{"email", "name"}
	ab.Config.Modules.LogoutMethod = "GET"

	ab.Config.Modules.ResponseOnUnauthed = authboss.RespondRedirect

	defaults.SetCore(&ab.Config, false, true)

	emailRule := defaults.Rules{
		FieldName: "email", Required: true,
		MatchError: "Must be a valid e-mail address",
		MustMatch:  regexp.MustCompile(`.*@.*\.[a-z]+`),
		MaxLength:  1024,
	}
	passwordRule := defaults.Rules{
		FieldName: "password", Required: true,
		MinLength: 8,
	}
	nameRule := defaults.Rules{
		FieldName: "name", Required: true,
		MinLength: 2,
		MaxLength: 36,
	}

	ab.Config.Core.BodyReader = defaults.HTTPBodyReader{
		ReadJSON: false,
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
