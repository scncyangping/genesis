package util

import (
	"fmt"
	"genesis/pkg/util/snowflake"
	"strconv"
	"testing"
)

func TestShort(t *testing.T) {
	c := BinaryConvert
	m := make(map[string]bool)
	cm := make(map[string]bool)
	//oneId := int64(0)
	for i := 0; i <= 1000001; i++ {

		id := snowflake.NextIdInt64()
		//if i == 1000000 {
		//	fmt.Println("oneId: ", oneId)
		//	id = oneId
		//}
		if m[strconv.FormatInt(id, 10)] == true {
			t.Errorf("duplicated id %d", id)
			fmt.Println("err")
			return
		} else {
			m[strconv.FormatInt(id, 10)] = true
		}
		eid := c.Encode(id)
		if cm[eid] == true {
			fmt.Println("err")
			t.Errorf("duplicated eid %s", eid)
			return
		} else {
			cm[eid] = true
		}
		if i == 0 {
			//oneId = id
			fmt.Println("start")
			fmt.Println("id ", id, "eid ", eid)
		}
		if i > 1000000 {
			fmt.Println("end")
			fmt.Println("id ", id, "eid ", eid)
		}
	}
}
