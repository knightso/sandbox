// +build appengine
package lib

func init() {
	delegate = func() string {
		return "Hello GAE!"
	}
}
