package middleware

import (
	"context"
	"net/http"

	"github.com/spf13/viper"
	"github.com/volatiletech/authboss/v3"
)

// ConfigDataInject is middleware for injecting configuration data
func ConfigDataInject() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		h := func(w http.ResponseWriter, r *http.Request) {
			var data authboss.HTMLData
			dataIntf := r.Context().Value(authboss.CTXKeyData)

			if dataIntf == nil {
				data = authboss.HTMLData{}
			} else {
				data = dataIntf.(authboss.HTMLData)
			}

			data.MergeKV("websockets_url", viper.GetString("websocketsURL"))

			r = r.WithContext(context.WithValue(r.Context(), authboss.CTXKeyData, data))
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(h)
	}
}
