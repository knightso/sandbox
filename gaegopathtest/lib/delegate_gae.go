// +build appengine
package lib

import (
	"fmt"
)

func init() {
	fmt.Println("init for gae!")
	delegate = func() string {
		return "Hello GAE!"
	}
}
