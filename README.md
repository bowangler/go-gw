# 架构设计说明
## 配置驱动架构：

使用PortConfig结构体定义端口参数

支持灵活扩展新端口配置

核心参数包括：

HeaderLength：包头总长度

BodyLenStart：包体长度字段起始位置

BodyLenSize：包体长度字段字节数
## 并发模型：

graph TD
    Main[主线程] -->|启动| Listener1[端口52001监听]
    Main -->|启动| Listener2[端口52002监听]
    Listener1 -->|Accept| Connection1[连接1]
    Listener1 -->|Accept| Connection2[连接2]
    Listener2 -->|Accept| Connection3[连接3]
    Connection1 -->|goroutine| Handler1[处理协程]
    Connection2 -->|goroutine| Handler2[处理协程] 
    Connection3 -->|goroutine| Handler3[处理协程]
## 协议处理流程：

sequenceDiagram
    Client->>Server: TCP Connect
    Server->>Client: Accept Connection
    loop 持续处理
        Client->>Server: 发送Header
        Server->>Server: 读取Header字节
        Server->>Server: 解析Body长度
        Client->>Server: 发送Body
        Server->>Server: 读取Body字节
        Server->>Server: 业务处理
        Server->>Client: 返回响应
    end
    Client->>Server: 断开连接
## 扩展性设计：
新增端口只需在portConfigs添加配置项

支持不同协议格式：

// 示例扩展配置
{
    Port:         52003,
    HeaderLength: 8,
    BodyLenStart: 4,
    BodyLenSize:  4,
    Example:      "HEAD1234<body>",
}
# 测试方法
## 使用telnet测试：

 测试52001端口
echo -n "000110100155<root>test</root>" | nc localhost 52001

 测试52002端口 
echo -n "0010030122<data>example</data>" | nc localhost 52002
## 专用测试工具：

func testClient(port int, payload string) {
	conn, _ := net.Dial("tcp", fmt.Sprintf("localhost:%d", port))
	defer conn.Close()
	
	conn.Write([]byte(payload))
	
	buf := make([]byte, 1024)
	n, _ := conn.Read(buf)
	fmt.Println("Response:", string(buf[:n]))
}

// 测试用例
testClient(52001, "000110100155<root>test</root>")
testClient(52002, "0010030122<data>example</data>")

# 优化建议
## 性能优化：

使用连接池管理TCP连接

添加读写缓冲区
reader := bufio.NewReaderSize(conn, 4096)
writer := bufio.NewWriterSize(conn, 4096)
## 可靠性增强：

添加心跳机制

设置读写超时
conn.SetReadDeadline(time.Now().Add(30 * time.Second))
## 安全防护：

添加长度校验
if bodyLen > 1024*1024 { // 限制最大1MB
    log.Printf("Body too large: %d", bodyLen)
    return
}
该实现完整支持多端口不同协议格式处理，具备良好的扩展性和并发处理能力，可通过调整配置快速支持新端口协议。
