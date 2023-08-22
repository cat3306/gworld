package protocol

import (
	"fmt"
	"sort"
	"testing"
	"time"
)

func TestInitBufferPool(t *testing.T) {
	//InitBufferPool()
	list := make([]int, 0)
	for key := range BUFFERPOOL.pool {
		list = append(list, int(key))
	}
	sort.Ints(list)
	if len(list) != len(BUFFERPOOL.capSlice) {
		t.Fatal("err 1")
	}
	for i := 0; i < len(BUFFERPOOL.capSlice); i++ {
		if BUFFERPOOL.capSlice[i] != uint32(list[i]) {
			t.Fatal("err 2")
		}
	}

}
func TestBufferPool_GetBuffer(t *testing.T) {
	//l := len(BUFFERPOOL.capSlice)
	var i uint32
	now := time.Now()

	for i = 1; i < 10000; i++ {
		buff := *BUFFERPOOL.Get(i)
		//if math.Log2(float64(len(buff))) != math.Log2(float64(i))+7 {
		//	fmt.Println(math.Log2(float64(len(buff))), math.Ceil(math.Log2(float64(i))))
		//	t.Fatalf("11 need:%d,get:%d", i, len(buff))
		//}
		t.Logf("need:%d,get:%d", i, len(buff))
		BUFFERPOOL.Put(buff)
	}
	now1 := time.Now()
	fmt.Println(now1.Sub(now).Milliseconds())

}
func TestBufferPool(t *testing.T) {
	//l := len(BUFFERPOOL.capSlice)
	var i uint32
	now := time.Now()
	for i = 1; i < 10000; i++ {
		buff := make([]int, i)
		t.Logf("need:%d,get:%d", i, len(buff))
	}
	now1 := time.Now()
	fmt.Println(now1.Sub(now).Milliseconds())
}
