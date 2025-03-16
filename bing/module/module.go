package module

func Invoke(functions ...func() error) (err error) {
	for _, fn := range functions {
		if err = fn(); err != nil {
			break
		}
	}
	return
}
