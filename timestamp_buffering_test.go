
package tmbuf_test

import (
	"testing"
	// "sync"
	// "strconv"
	"time"
	// "math/rand"
	"sync/atomic"
	"fmt"

	"github.com/bmizerany/assert"
	"github.com/yyd01245/tmbuf"
	
)


// Dict implement the interface of dbuf.DoubleBufferingTarget
type Dict struct {
	d string
	// data []string
	data []interface{}
	//业务自己的其他更复杂的数据结构
}

// 创建缓存
func newDict() tmbuf.Target {
	d := new(Dict)
	d.data = make([]interface{},0)
	// return tmbuf.Target(d)
	return d
}

var initializedCount int32
var closedCount int32

func (d *Dict) Initialize(conf string) bool {
	// 这个conf一般情况下是一个配置文件的路径
	// 这里我们简单的认为它只是一段数据
	c := atomic.AddInt32(&initializedCount, 1)
	_ = c
	d.d = conf
	//fmt.Printf("Dict.Initialize() called, count=%d\n", c)
	return true
}

func (d *Dict) Reload() bool {
	// 在这里做一些资源释放工作
	fmt.Printf("Reload dict conf:%v\n",d.d)
	fmt.Printf("Reload dict data:%v\n",d.data)
	d.data = append(d.data,"ddddd")
	// d.data = append(d.data,"dddd")
	d.data[0] = fmt.Sprintf("%d",time.Now().Unix())

	return true
	//fmt.Printf("Dict.Close() called, count=%d\n", c)
}

func (d *Dict) Close() bool {
	// 在这里做一些资源释放工作
	c := atomic.AddInt32(&closedCount, 1)
	_ = c
	return true
	//fmt.Printf("Dict.Close() called, count=%d\n", c)
}

func (d *Dict) GetBuffer() []interface{} {
	return d.data
}

func TestTimestampBuffering(t *testing.T) {
	m := tmbuf.NewManager()
	rc := m.Add("mydict", "The config for Dict1", newDict) // 初始化的时候引用计数为1
	assert.Equal(t, rc, nil)
	rc_2 := m.Add("mydict", "The config for Dict2", newDict) // 初始化的时候引用计数为1
	assert.Equal(t, rc_2, nil)

	var Inteval int64 = 5
	d := m.Get(Inteval)
	assert.NotEqual(t, d, nil)
	// assert.Equal(t, initializedCount, int32(1))
	// assert.Equal(t, closedCount, int32(0))
	dict := d.Get()
	fmt.Printf("---get buffer:%v\n",dict.GetBuffer())
	assert.Equal(t,dict.GetBuffer()[0],fmt.Sprintf("%d",time.Now().Unix()))
	// var wg sync.WaitGroup
	ReloadedCount := 1000
	// ReloadedCount := 2
	go func() {
		for i := 0; i < ReloadedCount; i++ {
			m.Reload(Inteval)
			time.Sleep(time.Duration(Inteval) * time.Second)
		}
	}()


	for i := 0; i < ReloadedCount; i++ {
		// 模拟一堆协程在同时使用Dict对象
		// wg.Add(1)
		// go func() {
			// defer wg.Done()
			d := m.Get(Inteval)
			dict := d.Get()

			fmt.Printf("---get buffer:%v, now:%v\n",dict.GetBuffer(),fmt.Sprintf("%d",time.Now().Unix()))
			// assert.Equal(t,dict.GetBuffer()[0],fmt.Sprintf("%d",time.Now().Unix()))
			assert.NotEqual(t, dict, nil)
			// assert.Equal(t, tg.Ref() >= 1, true)
			time.Sleep(time.Duration(Inteval) * time.Second)
		// }()

		// 模拟不定期的字典文件重新加载
		// wg.Add(1)
		// go func() {
		// 	defer wg.Done()
		// 	assert.Equal(t, m.Reload("mydict", "The new config for Dict" + strconv.Itoa(i)), true)
		// 	time.Sleep(time.Duration(rand.Intn(100000) + 1) * time.Microsecond)
		// }()
	}
	// wg.Wait()

	// assert.Equal(t, initializedCount, int32(ReloadedCount) + 1)
	// assert.Equal(t, closedCount, int32(ReloadedCount))

	// tg := d.Get()
	// defer tg.Release()
	// assert.Equal(t, tg.Ref(), int32(2))// 初始化的时候引用计数为1，Get 之后，引用计数又自动加1，因此这里为2。
}