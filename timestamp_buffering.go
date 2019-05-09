package tmbuf

import (
	"sync"
	"time"
	"fmt"
)

// TimestampBufferingTarget is the interface that wraps the basize operations
// 
// Initialize does some initialization operations on the resource.
// When Initialize encounters an error implementations must return false.
//
// Close does some resource recycling works which cannot be done
// by GC of Golang.
type Target interface {
	Initialize(param string) bool
	Reload() bool
	Close() bool
	GetBuffer() []interface{} 
}

type TargetCreator func() Target

// type TargetRef struct {
// 	Target Target
// 	// ref    *int32
// }

type TimestampBuffering struct {
	creator         TargetCreator
	buffer								Target
	mutex           sync.Mutex
	// refTarget       TargetRef

	reloadTimestamp int64
}


func newTimestampBuffering(f TargetCreator) *TimestampBuffering {
	d := new(TimestampBuffering)
	d.creator = f
	d.reloadTimestamp = 0
	return d
}


type dbmap []*TimestampBuffering

type Manager struct {
	targets dbmap
	mutex   sync.Mutex
}

func NewManager() *Manager {
	m := new(Manager)
	m.targets = []*TimestampBuffering {
		// new(TimestampBuffering),
		// new(TimestampBuffering),
	}
	// m.targets = append(m.targets,new(TimestampBuffering))
	return m
}

func (d *TimestampBuffering) reload() error {
	t := d.buffer 
	if t.Reload() == false {
		return fmt.Errorf("t.Reload() failed\n")
	}

	// d.reloadTimestamp = time.Now().Unix()

	// d.mutex.Lock()
	// defer d.mutex.Unlock()
	// d.refTarget.Release() // 将老对象释放掉

	return nil
}

func (d *TimestampBuffering) initialize(conf string) error {
	t := d.creator()
	if t.Initialize(conf) == false {
		return fmt.Errorf("t.Initialize(%v) failed\n", conf)
	}
	d.buffer = t
	return nil
}

func (d *TimestampBuffering) Get() Target {
	// d.mutex.Lock()
	// defer d.mutex.Unlock()
	// atomic.AddInt32(d.refTarget.ref, 1)
	return d.buffer
}


func (m *Manager) Add(name string, conf string, f TargetCreator) (err error) {
	d := newTimestampBuffering(f)
	d.initialize(conf)
	err = d.reload()
	if err == nil {
		// m.targets[name] = d
		// time.Now().Unix() 
		// now := time.Now().Unix()
		// index: = (now / 35);
		m.targets = append(m.targets,d)
	}

	return err
}

func (m *Manager) GetCurrentIndex(base int64) int {

	targetLen := int64(len(m.targets))
	now := time.Now().Unix()
	fmt.Printf("GetCurrentIndex now time: %v\n",now)
	// get now
	index := ((now / base) + targetLen + 1) % targetLen;
	if index >= targetLen {
		fmt.Errorf("get TimestampBuffering Current index error, index = %d, fix to:%d!\n",index,targetLen-1)
		index = targetLen-1
	}
	return int(index)
}

func (m *Manager) GetIdleIndex(base int64) int {

	targetLen := int64(len(m.targets))
	now := time.Now().Unix()
	fmt.Printf("GetIdleIndex now time: %v\n",now)
	// get now
	index := ((now / base) + targetLen ) % targetLen;
	if index >= targetLen {
		fmt.Errorf("get TimestampBuffering Idle index error, index = %d, fix to:%d!\n",index,0)
		index = 0
	}
	return int(index)
}

// time sycle base as user 
func (m *Manager) Get(base int64) *TimestampBuffering {
	// m.mutex.Lock()
	// defer m.mutex.Unlock()

	// get now
	index := m.GetCurrentIndex(base)
	fmt.Printf("get current index :%d\n",index)


	return m.targets[index]

	// return nil
}

// 需要错开当前在使用的, base 必须与 Get 的一致
func (m *Manager) Reload(base int64) error {
	index := m.GetIdleIndex(base)

	fmt.Printf("get idle index :%d\n",index)
	d := m.targets[index]

	if d == nil {
		return fmt.Errorf("Cannot find the Target by index [%v]\n", index)
	}

	return d.reload()
}
