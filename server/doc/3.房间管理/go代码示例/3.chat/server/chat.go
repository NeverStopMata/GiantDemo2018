package main

type ChatInfo struct {
	who *EchoTask
	msg []byte
}

var chanChat chan *ChatInfo = make(chan *ChatInfo, 10)
var chanAdd chan *EchoTask = make(chan *EchoTask, 10)
var chanRemove chan *EchoTask = make(chan *EchoTask, 10)
var users []*EchoTask

func doChat() {
	for {
		select {
		case info := <-chanChat:
			{
				for _, v := range users {
					if v != info.who {
						v.AsyncSend(info.msg, 0)
					}
				}
			}
		case who := <-chanAdd:
			{
				users = append(users, who)
			}
		case who := <-chanRemove:
			{
				for k, v := range users {
					if v == who {
						users = append(users[:k], users[k+1:]...)
						break
					}
				}
			}
		}
	}
}
