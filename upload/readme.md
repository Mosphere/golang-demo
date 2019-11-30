有时候遇到个别文件或者图片资源需要上传至服务端又不想下载ftp等客户端工具，最近正也好在学golang，所以就以golang实现相关功能。
服务端架构:C/S, 服务端运行server.go 然后本地通过go run client.go pathfile将本地文件上传至server.go所在的服务器
实现逻辑主要如下：
### 1.客户端发送文件名给服务端
执行go client.go pathfile后,通过os.Args获取命令行参数
```
	list := os.Args
	if len(list) != 2{
		fmt.Println("参数格式不正确")
		return
	}
```
然后再提取文件名将其发送给服务端
```
filePath := list[1]
	fileInfo, err := os.Stat(filePath)
  conn, err := net.Dial("tcp", "ip:8007")
  defer conn.Close()
	//向服务端发送文件名
	filename := fileInfo.Name()
	_, err = conn.Write([]byte(filename))
  ```
  服务端获取客户端传送过来的文件名后创建相关文件
  ```
  f, err := os.Create(filename)
  ```
  ### 2.客户端发送文件内容给服务端
  ```
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
```
服务端接收文件内容
```
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
```
`注`
```
若客户端和服务端不在同一台主机，server.go中net.Listen中的address由`127.0.0.1`改为`0.0.0.0`,然后在安全组中放开8088端口允许外部访问
```
  
