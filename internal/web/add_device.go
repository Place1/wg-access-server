package web

import (
	"encoding/json"
	"net/http"

	"github.com/place1/wireguard-access-server/internal/services"
	"github.com/place1/wireguard-access-server/internal/storage"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type AddDeviceRequest struct {
	Name      string `json:"name"`
	PublicKey string `json:"publicKey"`
}

type AddDeviceResponse struct {
	Device *storage.Device `json:"device"`
}

func AddDevice(devices *services.DeviceManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		req := AddDeviceRequest{}
		if err := decoder.Decode(&req); err != nil {
			logrus.Error(errors.Wrap(err, "unable to decode request body"))
			http.Error(w, "bad request payload", http.StatusBadRequest)
			return
		}

		device, err := devices.AddDevice(req.Name, req.PublicKey)
		if err != nil {
			logrus.Error(errors.Wrap(err, "unable to add device"))
			http.Error(w, "failed to add the new device", http.StatusInternalServerError)
			return
		}

		response, err := json.Marshal(AddDeviceResponse{
			Device: device,
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
