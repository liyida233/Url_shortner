package store

type Store interface {
	Save(shortUrl, longUrl, userId, customAlias string) error
	Get(shortUrl string) (string, error)
	Exists(shortUrl string) bool
}
type VisitTracker interface {
	IncrementVisitCount(shortUrl string)
}
