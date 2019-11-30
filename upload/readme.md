有时候遇到个别文件或者图片资源需要上传至服务端又不想下载ftp等客户端工具，最近正也好在学golang，所以就以golang实现相关功能。
服务端架构:C/S, 服务端运行server.go 然后本地通过go run client.go pathfile将本地文件上传至server.go所在的服务器
