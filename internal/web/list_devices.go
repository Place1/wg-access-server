package web

import (
	"encoding/json"
	"net/http"

	"github.com/alexedwards/scs/v2"

	"github.com/pkg/errors"
	"github.com/place1/wireguard-access-server/internal/services"
	"github.com/place1/wireguard-access-server/internal/storage"
	"github.com/sirupsen/logrus"
)

type ListDeviceRequest struct{}

type ListDeviceResponse struct {
	Items []*storage.Device `json:"items"`
}

func ListDevices(session *scs.SessionManager, devices *services.DeviceManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := session.GetString(r.Context(), "auth/subject")

		devices, err := devices.ListDevices(user)
		if err != nil {
			logrus.Error(errors.Wrap(err, "failed to list devices"))
			http.Error(w, "failed to list devices", http.StatusInternalServerError)
			return
		}

		response, err := json.Marshal(ListDeviceResponse{
			Items: devices,
		})
		if err != nil {
			logrus.Error(errors.Wrap(err, "failed to marshal response"))
			http.Error(w, "failed to marshal response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
		return
	}
}
