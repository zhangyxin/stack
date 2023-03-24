package stack

import "sync"

type FStack struct {
	data  []interface{}
	top   uint64
	cap   uint64
	mutex sync.Mutex
}

type IStack interface{
	// Iterator() (i *Iterator.Iterator)
	Size() uint64
	Clear()
	Empty() bool
	Pop()
	Top() interface{}
	Push(e interface{})
}

func New()(s *FStack)  {
	return &FStack{
		data: make([]interface{},1,1),
		top: 0,
		cap: 1,
		mutex: sync.Mutex{},
	};
}

func (s*FStack) Size() uint64 {
	if(s == nil){
		s = New();
	}
	return s.top;
}

func (s *FStack) Clear() {
    if s == nil {
        s = New()
    }
    s.mutex.Lock()
    s.data = make([]interface{}, 0, 0)
    s.top = 0
    s.cap = 1
    s.mutex.Unlock()
}

func (s *FStack) Empty() (b bool) {
    if s == nil {
        return true
    }
    return s.Size() == 0
}

func (s *FStack) Push(e interface{}) {
    if s == nil {
        s = New()
    }
    s.mutex.Lock()
    if s.top < s.cap {
        //还有冗余,直接添加
        s.data[s.top] = e
    } else {
        //冗余不足,需要扩容
        if s.cap <= 65536 {
            //容量翻倍
            if s.cap == 0 {
                s.cap = 1
            }
            s.cap *= 2
        } else {
            //容量增加2^16
            s.cap += 65536
        }
        //复制扩容前的元素
        tmp := make([]interface{}, s.cap, s.cap)
        copy(tmp, s.data)
        s.data = tmp
        s.data[s.top] = e
	}
    s.top++
    s.mutex.Unlock()
}

func (s *FStack) Top() (e interface{}) {
    if s == nil {
        return nil
    }
    if s.Empty() {
        return nil
    }
    s.mutex.Lock()
    e = s.data[s.top-1]
    s.mutex.Unlock()
    return e
}

func (s *FStack) Pop() {
    if s == nil {
        s = New()
        return
    }
    if s.Empty() {
        return
    }
    s.mutex.Lock()
    s.top--
    if s.cap-s.top >= 65536 {
        //容量和实际使用差值超过2^16时,容量直接减去2^16
        s.cap -= 65536
        tmp := make([]interface{}, s.cap, s.cap)
        copy(tmp, s.data)
        s.data = tmp
    } else if s.top*2 < s.cap {
        //实际使用长度是容量的一半时,进行折半缩容
        s.cap /= 2
        tmp := make([]interface{}, s.cap, s.cap)
        copy(tmp, s.data)
        s.data = tmp
    }
    s.mutex.Unlock()
}