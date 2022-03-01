package storage

import "github.com/sirupsen/logrus"

type InProcessWatcher struct {
	add    []Callback
	delete []Callback
}

func NewInProcessWatcher() *InProcessWatcher {
	logrus.Debug("creating in-process watcher")
	return &InProcessWatcher{
		add:    []Callback{},
		delete: []Callback{},
	}
}

func (w *InProcessWatcher) OnAdd(cb Callback) {
	w.add = append(w.add, cb)
}

func (w *InProcessWatcher) OnDelete(cb Callback) {
	w.delete = append(w.delete, cb)
}

func (w *InProcessWatcher) OnReconnect(cb func()) {
	// noop because the inprocess watcher can't disconnect
}

func (w *InProcessWatcher) EmitAdd(device *Device) {
	// This also triggers on updates which influences performance with big callbacks for many active devices
	// As the InProcessWatcher is only used for in-memory databases for development, this is not a problem
	for _, cb := range w.add {
		cb(device)
	}
}

func (w *InProcessWatcher) EmitDelete(device *Device) {
	for _, cb := range w.delete {
		cb(device)
	}
}
