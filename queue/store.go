package queue

type Store struct {
	s map[string][]byte
}

func NewStore() *Store {
	return &Store{s: map[string][]byte{}}
}

func (s *Store) Get(key string) []byte {
	return s.s[key]
}

func (s *Store) Set(key string, b []byte) {
	s.s[key] = b
}
