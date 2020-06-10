package nsqd

// 小根堆
// inFlightPqueue队列按Message中的pri字段进行排序的（pri也是时间戳，是投递消息的超时时间）
type inFlightPqueue []*Message

func newInFlightPqueue(capacity int) inFlightPqueue {
	return make(inFlightPqueue, 0, capacity)
}

func (pq inFlightPqueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *inFlightPqueue) Push(x *Message) {
	n := len(*pq)
	c := cap(*pq)
	if n+1 > c { // 超过队列capacity则将capacity扩充两倍，
		npq := make(inFlightPqueue, n, c*2)
		copy(npq, *pq)
		*pq = npq
	}
	*pq = (*pq)[0 : n+1]
	x.index = n
	(*pq)[n] = x
	pq.up(n) //小根堆末尾插入数据，需往上调整
}

func (pq *inFlightPqueue) Pop() *Message {
	n := len(*pq)
	c := cap(*pq)
	pq.Swap(0, n-1)
	pq.down(0, n-1)
	if n < (c/2) && c > 25 {
		npq := make(inFlightPqueue, n, c/2)
		copy(npq, *pq)
		*pq = npq
	}
	x := (*pq)[n-1] //弹出旧的根节点
	x.index = -1
	*pq = (*pq)[0 : n-1]
	return x
}

func (pq *inFlightPqueue) Remove(i int) *Message {
	n := len(*pq)
	if n-1 != i {
		pq.Swap(i, n-1)
		pq.down(i, n-1) //down完后此时(*pq)[:n-2]为小根堆，
		pq.up(i)        // 相当于(*pq)[:n-2]为小根堆push (*pq)[n-2]元素后调整
	}
	x := (*pq)[n-1]
	x.index = -1
	*pq = (*pq)[0 : n-1]
	return x
}

// 小根堆的根节点，即pri最小的节点，若其pri小于等于max则弹出，否则返回 nil 和 pri-max
// 实例：消息 timeout 时间戳为 pri，传入当前时间戳，判断是否有 timeout 消息则弹出，否则返回还需等待时间
func (pq *inFlightPqueue) PeekAndShift(max int64) (*Message, int64) {
	if len(*pq) == 0 {
		return nil, 0
	}

	x := (*pq)[0]
	if x.pri > max {
		return nil, x.pri - max
	}
	pq.Pop()

	return x, 0
}

func (pq *inFlightPqueue) up(j int) {
	for {
		i := (j - 1) / 2 // parent
		if i == j || (*pq)[j].pri >= (*pq)[i].pri {
			break
		}
		pq.Swap(i, j)
		j = i
	}
}

// 节点i数据变化，需要往i下调整
func (pq *inFlightPqueue) down(i, n int) {
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && (*pq)[j1].pri >= (*pq)[j2].pri {
			j = j2 // = 2*i + 2  // right child
		}
		if (*pq)[j].pri >= (*pq)[i].pri {
			break
		}
		pq.Swap(i, j)
		i = j
	}
}
