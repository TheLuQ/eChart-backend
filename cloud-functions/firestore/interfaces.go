package firestore

type SheetStore interface {
	GetAllTitles() ([]SheetGroup, error)
	SearchByIds(ids []string) ([]Sheet, error)
	AddSheet(sheet *Sheet) error
	UpsertSheet(sheet *Sheet) error
	RemoveSheet(sheet *Sheet) error
}

type EventStore interface {
	AddEvent(event *Event) error
	GetShortEvents() ([]Event, error)
	GetEventDetails(ids []string) ([]Event, error)
}
