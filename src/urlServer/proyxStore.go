package urlServer

import (
	"log"
	"net/rpc"
)

type proxyStore struct {
	client   *rpc.Client
	urlStore *urlStore
}

func NewProxyStore(addr string) *proxyStore {
	client, err := rpc.DialHTTP("tcp", addr)
	if err != nil {
		log.Println("Error constructing ProxyStore:", err)
	}
	return &proxyStore{client: client, urlStore: NewUrlStore("")}
}

func (s *proxyStore) Get(key, url *string) error {
	if err := s.urlStore.Get(key, url); err == nil {
		return nil
	}
	if err := s.client.Call("store.Get", key, url); err != nil {
		return err
	}
	s.urlStore.set(key, url)
	return nil
}

func (s *proxyStore) Put(url, key *string) error {
	if err := s.client.Call("store.Put", url, key); err != nil {
		log.Println("插入失败")
		return err
	}
	s.urlStore.set(key, url)
	return nil
}
