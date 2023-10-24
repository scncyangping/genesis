package file

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	c := 'A'
	for index := range []int{1, 2, 3, 4, 5, 6} {
		for i2 := range []int{1, 2, 3, 4, 5} {
			fmt.Println(fmt.Sprintf("%c%d", c+int32(index), i2))
		}
	}
}
