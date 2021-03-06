package ws

import (
	"container/list"
)

// 事件类型
type EventType int

const (
	EVENT_JOIN    = iota // 加入
	EVENT_LEAVE          // 离开
	EVENT_MESSAGE        //留言
)

type Event struct {
	Type EventType // JOIN：加入, LEAVE：离开, MESSAGE：留言
	User string    // 用户名
	//Tag       string    // 标签
	Timestamp int    // Unix timestamp (secs)
	Content   []byte // 留言内容
}

//归档
const archiveSize = 20

// 活动档案。
var archive = list.New()

// NewArchive 将新事件保存到存档列表。
func NewArchive(event Event) {
	if archive.Len() >= archiveSize {
		archive.Remove(archive.Front())
	}
	archive.PushBack(event)
}

// GetEvents 返回上次接收后的所有事件。
func GetEvents(lastReceived int) []Event {
	events := make([]Event, 0, archive.Len())
	for event := archive.Front(); event != nil; event = event.Next() {
		e := event.Value.(Event)
		if e.Timestamp > int(lastReceived) {
			events = append(events, e)
		}
	}
	return events
}
