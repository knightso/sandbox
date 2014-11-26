// +build !appengine

package lib

import "syscall"
import "fmt"

func init() {
	fmt.Println("init for normal go!")
	delegate = func() string {
		return fmt.Sprintf("Hello! pid:%d", syscall.Getpid())
	}
}
