package main // 定义包名为 main，表示这是一个可独立运行的程序

import (
	"bufio" // 缓冲读写，提供了带缓冲的 I/O 操作
	"flag"  // 命令行参数解析
	"fmt"   // 格式化 I/O，提供打印输出功能
	"net"   // 网络相关功能，提供 TCP/IP 网络支持
	"os"    // 操作系统相关功能，提供文件和输入输出接口
)

// 从服务器读取消息的函数，持续运行
func read(conn net.Conn) {
	reader := bufio.NewReader(conn) // 创建一个新的读取器，用于从服务器连接中读取数据
	for {
		message, _ := reader.ReadString('\n') // 读取一行数据，直到遇到换行符 '\n'，忽略错误
		fmt.Print("\rReceived: ", message)    // 打印从服务器收到的消息，使用 '\r' 覆盖当前行
		fmt.Print("Enter message: ")          // 提示用户输入消息
	}
}

// 向服务器发送消息的函数，持续运行
func write(conn net.Conn) {
	reader := bufio.NewReader(os.Stdin) // 创建一个新的读取器，用于从标准输入（键盘）读取用户输入
	for {
		fmt.Print("Enter message: ")       // 提示用户输入消息
		text, _ := reader.ReadString('\n') // 读取用户输入的一行数据，直到遇到换行符 '\n'，忽略错误
		conn.Write([]byte(text))           // 将用户输入的消息转换为字节数组，发送给服务器，忽略错误
	}
}

func main() {
	// 定义一个命令行参数，用于指定服务器的 IP 地址和端口，默认值为 "127.0.0.1:8030"
	addrPtr := flag.String("ip", "127.0.0.1:8030", "IP:port string to connect to")
	flag.Parse() // 解析命令行参数

	conn, _ := net.Dial("tcp", *addrPtr)            // 尝试建立到服务器的 TCP 连接，忽略错误
	defer conn.Close()                              // 在程序结束时关闭连接
	fmt.Println("Connected to server at", *addrPtr) // 打印连接成功的信息

	// 异步读取和显示来自服务器的消息
	go read(conn) // 启动一个新的 Goroutine，执行 read 函数，读取服务器发送的消息

	// 开始获取和发送用户输入的消息
	write(conn) // 调用 write 函数，读取用户输入并发送给服务器
}
