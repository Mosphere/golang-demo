package main

import (
	"fmt"
	"net"
	"os"
)

func main(){
	listener, err := net.Listen("tcp", "127.0.0.1:8007")
	if err != nil{
		fmt.Println("net.Listen err: ", err)
		return
	}
	defer listener.Close()

	for{
		//阻塞等待客户端请求
		conn, err := listener.Accept()
		if err != nil{
			fmt.Println("listener.Accept err: ", err)
			return
		}

		//读取客户端发送过来的文件名
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)

		if err != nil{
			fmt.Println("os.Create err: ", err)
			return
		}
		_, err = conn.Write([]byte("ok"))
		if err != nil{
			fmt.Println("conn.Write err: ", err)
			return
		}

		fileName := string(buf[:n])
		go Receive(conn, fileName)
	}

}

func Receive(conn net.Conn, filename string){
	f, err := os.Create(filename)
	if err != nil{
		fmt.Println("os.Create err :", err)
	}
	defer f.Close()

	//读取客户端发送过来的文件内容，并将其写到本地文件

	for{
		buf := make([]byte, 4096)
		n, err := conn.Read(buf)
		if n == 0{
			fmt.Println("文件接收完毕")
			return
		}
		if err != nil{
			fmt.Println("conn.Read err :", err)
			return
		}
		_, _ = f.Write(buf[:n])
	}

}
