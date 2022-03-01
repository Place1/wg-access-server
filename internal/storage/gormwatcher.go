package storage

import (
	"reflect"
	"runtime"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

type GormWatcher struct {
	*gorm.Callback
	table string
}

func NewGormWatcher(db *gorm.DB, table string) *GormWatcher {
	logrus.Debug("creating gorm watcher")
	return &GormWatcher{
		Callback: db.Callback(),
		table:    table,
	}
}

func (w *GormWatcher) OnAdd(cb Callback) {
	name := runtime.FuncForPC(reflect.ValueOf(cb).Pointer()).Name()
	logrus.Debugf("OnAdd callback name %s", name)
	w.Callback.Create().Register(name, func(scope *gorm.Scope) {
		w.emit(cb, scope)

	})
}

func (w *GormWatcher) OnDelete(cb Callback) {
	name := runtime.FuncForPC(reflect.ValueOf(cb).Pointer()).Name()
	logrus.Debugf("OnDelete callback name %s", name)
	w.Callback.Delete().Register(name, func(scope *gorm.Scope) {
		w.emit(cb, scope)
	})
}

func (w *GormWatcher) OnReconnect(cb func()) {
	// noop because the watcher can't reconnect
}

func (w *GormWatcher) emit(cb Callback, scope *gorm.Scope) {
	if scope.TableName() == w.table {
		cb(*scope.Value.(**Device))
	}
}

func (w *GormWatcher) EmitAdd(device *Device) {
	// noop because we rely on gorm callback
}

func (w *GormWatcher) EmitDelete(device *Device) {
	// noop because we rely on gorm callback
}
