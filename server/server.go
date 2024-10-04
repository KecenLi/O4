package main // 定义包名为 main，表示这是一个可独立运行的程序

import (
	"bufio" // 缓冲读写，提供了带缓冲的 I/O 操作
	"flag"  // 命令行参数解析
	"fmt"   // 格式化 I/O，提供打印输出功能
	"net"   // 网络相关功能，提供 TCP/IP 网络支持
)

// 定义一个消息结构体，用于存储发送者 ID 和消息内容
type Message struct {
	sender  int    // 发送者的客户端 ID
	message string // 消息内容
}

var clients = make(map[int]net.Conn) // 创建一个映射，用于存储客户端 ID 和对应的连接

// 接受新的客户端连接的函数，持续运行
func acceptConns(ln net.Listener, conns chan net.Conn) {
	for {
		conn, _ := ln.Accept() // 等待并接受新的客户端连接，忽略错误
		conns <- conn          // 将新的连接发送到通道 conns
	}
}

// 处理客户端消息的函数，为每个客户端启动一个协程
func handleClient(client net.Conn, clientid int, msgs chan Message) {
	reader := bufio.NewReader(client) // 创建一个新的读取器，从客户端连接中读取数据
	for {
		message, _ := reader.ReadString('\n') // 从客户端读取一行数据，直到遇到换行符 '\n'，忽略错误
		// 创建一个新的消息结构体，包含发送者 ID 和消息内容
		msg := Message{
			sender:  clientid,
			message: message,
		}
		msgs <- msg // 将消息发送到通道 msgs，供服务器主循环处理
	}
}

func main() {
	// 定义一个命令行参数，用于指定服务器监听的端口，默认值为 ":8030"
	portPtr := flag.String("port", ":8030", "port to listen on")
	flag.Parse() // 解析命令行参数

	ln, _ := net.Listen("tcp", *portPtr)                 // 创建一个 TCP 监听器，监听指定的端口，忽略错误
	defer ln.Close()                                     // 在程序结束时关闭监听器
	fmt.Println("Server is listening on port", *portPtr) // 打印服务器启动信息

	conns := make(chan net.Conn) // 创建一个通道，用于存储新的客户端连接
	msgs := make(chan Message)   // 创建一个通道，用于存储收到的消息
	var clientID int = 0         // 初始化客户端 ID

	go acceptConns(ln, conns) // 启动一个协程，持续接受新的客户端连接
	for {
		select {
		case conn := <-conns: // 如果有新的客户端连接
			clientID++                                       // 增加客户端 ID
			id := clientID                                   // 当前客户端的 ID
			clients[id] = conn                               // 将新的客户端连接存储到映射中
			fmt.Println("New client connected with ID:", id) // 打印新客户端连接信息
			go handleClient(conn, id, msgs)                  // 启动一个协程，处理该客户端的消息
		case msg := <-msgs: // 如果收到新的消息
			fmt.Printf("Message received from client %d: %s", msg.sender, msg.message) // 打印收到的消息
			// 将消息发送给除发送者之外的所有客户端
			for id, conn := range clients {
				if id != msg.sender {
					conn.Write([]byte(msg.message)) // 将消息发送给其他客户端，忽略错误
				}
			}
		}
	}
}
