package main

import (
    "bufio"
    "fmt"
    "net"
)

var quitSemaphore chan bool

func main() {
    var tcpAddr *net.TCPAddr
    tcpAddr, _ = net.ResolveTCPAddr("tcp", "127.0.0.1:9999")

    conn, _ := net.DialTCP("tcp", nil, tcpAddr)
    defer conn.Close()
    fmt.Println("connected!")

//    go onMessageRecived(conn)
//声明一个管道用于接收解包的数据
 	go  handleConnection(conn)
    // 控制台聊天功能加入
    for {
        var msg string
        fmt.Scan(&msg)
        if msg == "quit" {
            break
        }
//        b := []byte(msg + "\n")	
//		b, _ := Encode(msg + "\n")
		b := Packet([]byte (msg))
        conn.Write(b)
                 

//		SendMessage([]byte(msg + "\n"),conn)
    }
}

// func SendMessage(unpackMessage []byte, conn *net.TCPConn){
// 		b := Packet(unpackMessage)
//        conn.Write(b)
// }
 
 func onMessageRecived(conn *net.TCPConn) {
 	tmpBuffer := make([]byte, 0)
 	buffer := make([]byte, 1024)
 	 
    reader := bufio.NewReader(conn)
    //声明一个管道用于接收解包的数据
    readerChannel := make(chan []byte, 16)
    go readerData(readerChannel,conn)
    for {
    	
//        n, err := conn.Read(buffer)
		n, err :=  reader.Read(buffer)
        if err != nil {
            Log(conn.RemoteAddr().String(), " connection error: ", err)
            return
        }

        tmpBuffer = Unpack(append(tmpBuffer, buffer[:n]...), readerChannel)
    }
//    
}
 
func handleConnection(conn *net.TCPConn) {
    //声明一个临时缓冲区，用来存储被截断的数据
    tmpBuffer := make([]byte, 0)

	reader := bufio.NewReader(conn)
    //声明一个管道用于接收解包的数据
    readerChannel := make(chan []byte, 16)
    go readerData(readerChannel,conn)

    

    buffer := make([]byte, 1024)
    for {
    	
//        n, err := conn.Read(buffer)
		n, err :=  reader.Read(buffer)
        if err != nil {
            Log(conn.RemoteAddr().String(), " connection error: ", err)
            return
        }

        tmpBuffer = Unpack(append(tmpBuffer, buffer[:n]...), readerChannel)
    }
}


func readerData(readerChannel chan []byte, conn *net.TCPConn) {
    for {
        select {
        	case data := <-readerChannel:
            	Log(string(data))
        }
    }
}

func Log(v ...interface{}) {
    fmt.Println(v...)
}




