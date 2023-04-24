package services

import (
	"fmt"
	"net/http"

	"github.com/freifunkMUC/wg-access-server/internal/devices"
)

func HealthEndpoint(d *devices.DeviceManager) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := d.Ping(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "ping failed")
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "ok")
	})
}
