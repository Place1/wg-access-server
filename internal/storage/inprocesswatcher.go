package storage

type InProcessWatcher struct {
	add    []Callback
	delete []Callback
}

func NewInProcessWatcher() *InProcessWatcher {
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
	for _, cb := range w.add {
		cb(device)
	}
}

func (w *InProcessWatcher) EmitDelete(device *Device) {
	for _, cb := range w.delete {
		cb(device)
	}
}
