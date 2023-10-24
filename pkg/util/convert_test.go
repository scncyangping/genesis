package util

import (
	"fmt"
	"testing"
)

func TestCamelize(t *testing.T) {
	fmt.Println(Camelize("sys_user", false))
	fmt.Println(Camelize("sys_user", false))
}
