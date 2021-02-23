package ws

/**
聊天室
*/
import (
	"container/list"
	"fmt"
	//"encoding/json"
	"github.com/gorilla/websocket"
	"time"
)

//订阅
type Subscription struct {
	Archive []Event      // 所有的事件
	New     <-chan Event // 新事件即将到来
}

//新事件
func NewEvent(ep EventType, user string, msg []byte) Event {
	return Event{ep, user, int(time.Now().Unix()), msg}
}

// 加入聊天室
func Join(user string, ws *websocket.Conn) {
	Subscribe <- Subscriber{Name: user, Conn: ws}
	//fmt.Println("离开聊天室", u ser, Subscribe)
}

// 离开
func Leave(user string) {
	Unsubscribe <- user
	//fmt.Println("离开聊天室", user, Unsubscribe)
}

// 用户
type Subscriber struct {
	Name string          //用户名
	Conn *websocket.Conn // 仅限WebSocket用户；否则为零。
}

var (
	// 用户新加入通道
	Subscribe = make(chan Subscriber, 10)
	// 用户退出通道
	Unsubscribe = make(chan string, 10)
	// 在此发送活动以发布它们。
	Publish = make(chan Event, 10)
	// 长轮询等候名单。
	WaitingList = list.New()
	//用户列表
	Subscribers = list.New()
)

// 此函数处理所有传入的chan消息。
func ChatRoom() {
	for {
		select {
		case sub := <-Subscribe:
			if !isUserExist(Subscribers, sub.Name) {
				Subscribers.PushBack(sub) // 将用户添加到列表末尾。
				// 发布JOIN事件。
				Publish <- NewEvent(EVENT_JOIN, sub.Name, nil)
				//fmt.Println("新用户:", sub.Name, ";WebSocket:", sub.Conn != nil)
			} else {
				//fmt.Println("老用户:", sub.Name, ";WebSocket:", sub.Conn != nil)
			}
		case event := <-Publish:
			// Notify waiting list.
			for ch := WaitingList.Back(); ch != nil; ch = ch.Prev() {
				ch.Value.(chan bool) <- true
				WaitingList.Remove(ch)
			}

			broadcastWebSocket(event)
			NewArchive(event)

			if event.Type == EVENT_MESSAGE {
				//cont := event.Content
				//var str string
				//if len(cont) > 100 {
				//	str = string(cont)[:100]
				//} else {
				//	str = string(cont)
				//}
				//fmt.Println("消息来自", event.User, ";消息内容:", str)
			}
		case unsub := <-Unsubscribe:
			for sub := Subscribers.Front(); sub != nil; sub = sub.Next() {
				if sub.Value.(Subscriber).Name == unsub {
					Subscribers.Remove(sub)
					// Clone connection.
					ws := sub.Value.(Subscriber).Conn
					if ws != nil {
						ws.Close()
						fmt.Println("WebSocket已关闭:", unsub)
					}
					Publish <- NewEvent(EVENT_LEAVE, unsub, nil) // Publish a LEAVE event.
					break
				}
			}
		}
	}
}

/*func init() {
	go ChatRoom()
}*/

//用户是否存在
func isUserExist(subscribers *list.List, user string) bool {
	for sub := subscribers.Front(); sub != nil; sub = sub.Next() {
		if sub.Value.(Subscriber).Name == user {
			return true
		}
	}
	return false
}

// broadcastWebSocket 向WebSocket用户广播消息。
func broadcastWebSocket(event Event) {
	//data, err := json.Marshal(string())
	//if err != nil {
	//	logs.Error("未能编组事件：", err)
	//	return
	//}
	//如果列表不为空，则第一个元素。
	//返回下一个list元素或nil。

	var sub *list.Element
	for sub = Subscribers.Front(); sub != nil; sub = sub.Next() {
		// 立即向WebSocket用户发送事件。
		s := sub.Value.(Subscriber)
		ws := s.Conn
		if ws != nil {
			if ws.WriteMessage(websocket.TextMessage, event.Content) != nil {
				// 用户断开连接
				Unsubscribe <- sub.Value.(Subscriber).Name
			}
			//if len(event.Tag) > 0 {
			//	switch event.Tag {
			//	case "manager":
			//		toUser(ws, sub, event, ManagerList, s.Name)
			//		break
			//	case "user":
			//		toUser(ws, sub, event, UserList, s.Name)
			//		break
			//	}
			//} else {
			//	if len(event.User) > 0 {
			//		switch event.User {
			//		case s.Name:
			//			if ws.WriteMessage(websocket.TextMessage, event.Content) != nil {
			//				// 用户断开连接
			//				Unsubscribe <- sub.Value.(Subscriber).Name
			//			}
			//			break
			//		}
			//	} else {
			//		if ws.WriteMessage(websocket.TextMessage, event.Content) != nil {
			//			// 用户断开连接
			//			Unsubscribe <- sub.Value.(Subscriber).Name
			//		}
			//	}
			//}
			// 广播消息
		}
	}
}

func toUser(ws *websocket.Conn, sub *list.Element, event Event, userList []string, user string) {
	for _, v := range userList {
		if v == user {
			if ws.WriteMessage(websocket.TextMessage, event.Content) != nil {
				// 用户断开连接
				Unsubscribe <- sub.Value.(Subscriber).Name
			}
		}
	}
}
