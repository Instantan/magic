package magic

// Must panics if the passed error is not nil
func Must[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}
	return val
}
