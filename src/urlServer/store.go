package urlServer

type Store interface {
	Get(key, url *string) error
	Put(url, key *string) error
}
