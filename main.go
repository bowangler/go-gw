package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"sync"
)

// 端口配置结构体
type PortConfig struct {
	Port         int    // 监听端口
	HeaderLength int    // 包头长度
	BodyLenStart int    // 包体长度字段起始位置
	BodyLenSize  int    // 包体长度字段字节数
	Example      string // 示例格式说明
}

// 全局端口配置
var portConfigs = []PortConfig{
	{
		Port:         52001,
		HeaderLength: 12,
		BodyLenStart: 8,  // 000110100155 中0155的起始位置
		BodyLenSize:  4,
		Example:      "000110100155<root>...</root>",
	},
	{
		Port:         52002,
		HeaderLength: 10,
		BodyLenStart: 6,  // 0010030122 中0122的起始位置
		BodyLenSize:  4,
		Example:      "0010030122<data>...</data>",
	},
}

func main() {
	var wg sync.WaitGroup

	// 启动多个端口监听
	for _, config := range portConfigs {
		wg.Add(1)
		go func(c PortConfig) {
			defer wg.Done()
			startTCPServer(c)
		}(config)
	}

	wg.Wait()
}

// 启动TCP服务器
func startTCPServer(config PortConfig) {
	addr := fmt.Sprintf(":%d", config.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Port %d listen failed: %v", config.Port, err)
	}
	defer listener.Close()

	log.Printf("Server started on port %d", config.Port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Port %d accept error: %v", config.Port, err)
			continue
		}

		go handleConnection(conn, config)
	}
}

// 处理客户端连接
func handleConnection(conn net.Conn, config PortConfig) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		// 读取包头
		header := make([]byte, config.HeaderLength)
		if _, err := io.ReadFull(reader, header); err != nil {
			if err != io.EOF {
				log.Printf("Port %d header read error: %v", config.Port, err)
			}
			return
		}
        log.Printf("Port %d header read: %v", config.Port, header)

		// 解析包体长度
		bodyLenStr := string(header[config.BodyLenStart : config.BodyLenStart+config.BodyLenSize])
		bodyLen, err := strconv.Atoi(bodyLenStr)
		if err != nil {
			log.Printf("Port %d invalid body length: %s", config.Port, bodyLenStr)
			return
		}

		// 读取包体
		body := make([]byte, bodyLen)
		if _, err := io.ReadFull(reader, body); err != nil {
			log.Printf("Port %d body read error: %v", config.Port, err)
			return
		}
        log.Printf("Port %d body read: %v", config.Port, body)

		// 业务处理
		processPacket(config.Port, header, body)

		// 示例响应
		response := fmt.Sprintf("Processed: %s%s", header, body)
		conn.Write([]byte(response))
	}
}

// 业务逻辑处理（示例）
func processPacket(port int, header, body []byte) {
	switch port {
	case 52001:
		// XML处理逻辑
		log.Printf("Port 52001 Received XML:\nHeader: %s\nBody: %s", header, body)
	case 52002:
		// 自定义数据处理逻辑
		log.Printf("Port 52002 Received Data:\nHeader: %s\nBody: %s", header, body)
	}
}
