package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

//go run client.go filename
func main(){
	//获取命令行参数
	list := os.Args
	if len(list) != 2{
		fmt.Println("参数格式不正确")
		return
	}

	filePath := list[1]
	fileInfo, err := os.Stat(filePath)
	if err != nil{
		fmt.Println("os.Stat err", err)
		return
	}

	conn, err := net.Dial("tcp", "118.24.148.138:8088")
	if err != nil{
		fmt.Println("net.Dial err", err)
		return
	}
	defer conn.Close()

	//向服务端发送文件名
	filename := fileInfo.Name()
	_, err = conn.Write([]byte(filename))
	if err != nil{
		fmt.Println("conn.Write err", err)
		return
	}

	//读取服务端回发的ok
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil{
		fmt.Println("conn.Read err", err)
		return
	}
	if string(buf[:n]) != "ok"{
		fmt.Println("服务端异常", err)
		return
	}

	fmt.Println("服务端已接收文件名")
	//将文件内容发送给服务端
	SendFile(conn, filePath)
}

func SendFile(conn net.Conn, filePath string){
	//读出本地文件内容，并将其发送给服务端
	file, err := os.Open(filePath)
	if err != nil{
		fmt.Println("os.Open err: ", err)
		return
	}
	defer file.Close()

	for{
		buf := make([]byte, 4096)
		n, err := file.Read(buf)
		if err != nil{
			if err == io.EOF{
				fmt.Println("文件已发送完成")
			}else{
				fmt.Println("file.Read err: ", err)
			}
			return
		}

		_, err = conn.Write(buf[:n])
		if err != nil{
			fmt.Println("conn.Write err: ", err)
			return
		}
	}
}
