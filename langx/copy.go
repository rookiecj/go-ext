package langx

func Copy[T any](src T, options ...func(ptr *T)) T {
	dup := src
	for _, opt := range options {
		opt(&dup)
	}
	return dup
}
