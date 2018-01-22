package main

type ChatInfo struct {
	who *EchoTask
	msg []byte
}

var chanChat chan *ChatInfo = make(chan *ChatInfo, 10)
var chanAdd chan *EchoTask = make(chan *EchoTask, 10)
var chanRemove chan *EchoTask = make(chan *EchoTask, 10)
var users map[int][]*EchoTask = make(map[int][]*EchoTask)
var curRoomId = 1

const DEFAULT_ROOM_MEMBER = 2

func doChat() {
	for {
		select {
		case info := <-chanChat:
			{
				if us, ok := users[info.who.roomId]; ok {
					for _, v := range us {
						if v != info.who {
							v.AsyncSend(info.msg, 0)
						}
					}
				}
			}
		case who := <-chanAdd:
			{
				if _, ok := users[curRoomId]; !ok {
					users[curRoomId] = make([]*EchoTask, 0)
				}

				if len(users[curRoomId]) >= DEFAULT_ROOM_MEMBER {
					curRoomId++
					users[curRoomId] = make([]*EchoTask, 0)
				}
				who.roomId = curRoomId
				users[curRoomId] = append(users[curRoomId], who)
			}
		case who := <-chanRemove:
			{
				if us, ok := users[who.roomId]; ok {
					for k, v := range us {
						if v == who {
							us = append(us[:k], us[k+1:]...)
							if len(us) == 0 {
								delete(users, who.roomId)
							}
							break
						}
					}
				}
			}
		}
	}
}
