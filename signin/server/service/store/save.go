package store

type DBer interface {
	Get(key string) (string, error)
	Set(key, value string) error
}

func Get(key string) (string, error) {
	if _store != nil {
		return _store.Get(key)
	}
	return "", nil
}

func Set(key, value string) error {
	if _store != nil {
		if err := _store.Set(key, value); err != nil {
			return err
		}
	}
	return nil
}

func New(db DBer) {
	_store = &store{db}
}

type store struct {
	DBer
}

func (s *store) Init() error {
	return nil
}

func (s *store) Start() error {
	return nil
}

func (s *store) Stop() error {
	return nil
}

var _store *store
