package web

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/gorilla/mux"

	"github.com/place1/wireguard-access-server/internal/services"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func DeleteDevice(session *scs.SessionManager, devices *services.DeviceManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name, ok := vars["name"]
		if !ok {
			http.Error(w, "missing device name in path", http.StatusBadRequest)
			return
		}

		user := session.GetString(r.Context(), "auth/email")

		if err := devices.DeleteDevice(user, name); err != nil {
			logrus.Error(errors.Wrap(err, "failed to remove device"))
			http.Error(w, "failed to remove device", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
		return
	}
}
