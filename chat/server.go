package main

import (
	"fmt"
	"net"
	"time"
)

type Client struct {
	C chan string
	Name string
	Addr string
}
var msg = make(chan string)
var msgMap = make(map[string]Client)
func main(){
	listener, err := net.Listen("tcp", "127.0.0.1:8090")
	if err != nil{
		fmt.Println("listen err :", err)
		return
	}
	defer listener.Close()
	go Manager()
	for{
		//阻塞等待客户端请求
		conn, err := listener.Accept()
		if err != nil{
			fmt.Println("listen err :", err)
			return
		}
		go HandleConnection(conn)
	}
}

//管理所有消息
func Manager(){
	for{
		text := <-msg
		for _,item := range msgMap{
			item.C<- text
		}
	}
}

//将消息分发给当前客户端
func WriteMsgToClient(conn net.Conn, clt Client){
	for text := range clt.C{
		_, _ = conn.Write([]byte(text + "\n"))
	}
}

func makeMsg(clt Client, msg string) string{
	return "[" + clt.Addr + "]" + "("+ clt.Name + ")" + ":" + msg
}
func HandleConnection(conn net.Conn){
	defer conn.Close()

	clientIP := conn.RemoteAddr().String()
	client := Client{make(chan string), clientIP, clientIP}
	fmt.Println(clientIP + "已上线")
	msgMap[clientIP] = client

	//将消息发送给当前客户端
	go WriteMsgToClient(conn, client)
	//登陆提示
	msg<- makeMsg(client, "login")
	quit := make(chan bool)	//客户端是否关闭状态
	isActive := make(chan bool)	//活跃状态
	//将客户端发送的消息同步到msg通道
	go func() {
		for{
			buf := make([]byte, 4096)
			n, err := conn.Read(buf)
			if n == 0{
				quit<- true
				fmt.Println(clientIP, "客户端已关闭")
				return
			}
			if err != nil{
				fmt.Println("conn.Read err :", err)
				return
			}

			msgText := string(buf[:n-1])	//windows下的netcat工具会在输入文本后多加个\n
			fmt.Println(len(msgText))
			if msgText == "who"{	//展示在线用户
				_, _ = conn.Write([]byte("online user list :\n"))
				for _, user := range msgMap{
					_, _ = conn.Write([]byte("[" + user.Addr + "]" + "(" + user.Name + ")\n"))
				}
			}else if len(msgText) > 7 && msgText[:7] == "rename|" {	//重命名
				newName := msgText[8:]
				client.Name = newName
				msgMap[clientIP] = client
				_, _ = conn.Write([]byte("rename successful\n"))
			}else{
				//将消息写入通道
				msg<- makeMsg(client, string(buf[:n]))
				fmt.Println(msgText)
			}
			isActive<- true
		}

	}()

	for{
		select{
		case <-quit :
			delete(msgMap, clientIP)
			msg<- makeMsg(client, " be logout for time out")
			return
		case <-isActive:

		case <-time.After(time.Second * 60):
			delete(msgMap, clientIP)
			msg<- makeMsg(client, " be logout for time out")
			return
		}
	}
}
