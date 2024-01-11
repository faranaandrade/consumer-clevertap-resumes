package location

import (
	"sync"
	"time"

	"github.com/occmundial/go-common/logger"
)

// defaultLocation: Nombre de la localidad predeterminada para datos regionales (fecha)
const defaultLocation = "America/Mexico_City"

var (
	once     sync.Once
	location *Location
)

type Location struct {
	log      Logger
	Location *time.Location
}

func GetLocation(log *logger.Log) *Location {
	once.Do(func() {
		location = NewLocation(log)
	})
	return location
}

func NewLocation(log *logger.Log) *Location {
	l := &Location{log: log}
	l.Location = l.loadLocation(defaultLocation)
	return l
}

func (l *Location) DateNow() time.Time {
	return l.DateToLocation(time.Now())
}

func (l *Location) DateToLocation(date time.Time) time.Time {
	return date.In(l.Location)
}

func (l *Location) loadLocation(name string) *time.Location {
	location, err := time.LoadLocation(name)
	if err != nil {
		l.log.Error("locations", "loadLocation", err)
		return location
	}
	return location
}
