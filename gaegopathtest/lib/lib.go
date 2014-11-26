package lib

var delegate func() string

func Hello() string {
	if delegate == nil {
		return ""
	}
	return delegate()
}
