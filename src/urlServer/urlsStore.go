package urlServer

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

type urlStore struct {
	urls     map[string]string
	pairChan chan record
	mu       sync.RWMutex
	fileName string
}
type record struct {
	Key, Url string
}

func (s *urlStore) recover() error {
	file, err := os.Open(s.fileName)
	defer file.Close()
	if err != nil {
		log.Println("找不到恢复文件", err)
		return err
	}
	de := json.NewDecoder(file)
	r := &record{}
	var err0 error
	for err0 == nil {
		if err0 = de.Decode(r); err0 == nil {
			s.set(&r.Key, &r.Url)
		}
	}
	if err == io.EOF {
		return nil
	}
	return err
}

func NewUrlStore(fileName string) *urlStore {
	s := &urlStore{urls: make(map[string]string)}
	if fileName != "" {
		s.fileName = fileName
		s.pairChan = make(chan record, 1000)
		if err := s.recover(); err != nil {
			log.Println("恢复错误：", err)
		}
		go s.loopSave()
	}
	go s.loopPrint()
	return s
}

func (s *urlStore) Get(key, url *string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if u, ok := s.urls[*key]; ok {
		*url = u
		return nil
	}
	return errors.New("不存在键对应的值")
}
func (s *urlStore) loopPrint() {
	for {
		time.Sleep(5e9)
		s.mu.RLock()
		fmt.Println(s.urls)
		s.mu.RUnlock()
	}
}
func (s *urlStore) set(key, url *string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, present := s.urls[*key]; present {
		return errors.New("键已经存在")
	}
	s.urls[*key] = *url
	return nil
}

func (s *urlStore) count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.urls)
}

func (s *urlStore) Put(url, key *string) error {
	for {
		*key = genKey(s.count())
		if s.set(key, url) == nil {
			s.pairChan <- record{*key, *url}
			return nil
		}
	}
	return errors.New("不应该到这里")
}

func (s *urlStore) loopSave() {
	f, err := os.OpenFile(s.fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Println("本地存储时打开文件失败:", err)
	}
	en := json.NewEncoder(f)
	for {
		pair := <-s.pairChan
		if err0 := en.Encode(pair); err0 != nil {
			log.Println("encode失败")
		}
	}
}
