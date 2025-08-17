package store

type Engine interface {
	Search(key string) (string, error)
	Save(key, value string) error
}

type store struct {
	e Engine
}

func Search(key string) (string, error) {
	if _store != nil {
		return _store.e.Search(key)
	}
	return "", nil
}

func Save(key, value string) error {
	if _store != nil {
		if err := _store.e.Save(key, value); err != nil {
			return err
		}
	}
	return nil
}

var _store *store
