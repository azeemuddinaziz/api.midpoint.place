package db

import (
	"time"

	"github.com/championswimmer/api.midpoint.place/src/utils/applogger"
	"github.com/samber/lo"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	_numAttempts = 10
	_delay       = 1 * time.Second
	_openSync    = lo.Synchronize()
)

func attemptToOpen[T any](opener func() (*T, error)) *T {
	var t *T
	_, _, e := lo.AttemptWithDelay(_numAttempts, _delay, func(i int, d time.Duration) (err error) {
		applogger.Info("Opening attempt ", i)
		_openSync.Do(func() {
			t, err = opener()
		})
		return
	})
	if t == nil {
		applogger.Panic("Failed to open ", e.Error())
	}
	return t
}

func ProvidePostgresDB(dbUrl string, config *gorm.Config) *gorm.DB {
	applogger.Warn("App: Using postgres db")
	return attemptToOpen(func() (*gorm.DB, error) {
		return gorm.Open(postgres.Open(dbUrl), config)
	})
}

func ProvideSqliteDB(dbUrl string, config *gorm.Config) *gorm.DB {
	applogger.Warn("App: Using sqlite db")
	return lo.Must(gorm.Open(sqlite.Open(dbUrl), config))
}
