package main

import (
    "bufio"
    "fmt"
    "net"
    "protos" 
    "github.com/golang/protobuf/proto"
)

var quitSemaphore chan bool

func main() {
    var tcpAddr *net.TCPAddr
    tcpAddr, _ = net.ResolveTCPAddr("tcp", "127.0.0.1:9999")

    conn, _ := net.DialTCP("tcp", nil, tcpAddr)
//    defer conn.Close()
    fmt.Println("connected!")

    go onMessageRecived(conn)
//声明一个管道用于接收解包的数据
// 	go  handleConnection(conn)
    sendMessageToServer2(conn)
    for{
    	
    }
}


func sendMessageToServer2(conn *net.TCPConn){
  msg := &protos.Helloworld{
        Id:  proto.Int32(101),
        Str: proto.String("hello"),
    } //msg init

//    path := string("./log.txt")
//    f, err := os.Create(path)
//    if err != nil {
//        fmt.Printf("failed: %s\n", err)
//        return
//    }
//
//    defer f.Close()
	Log("send msg:", msg)
    buffer, err := proto.Marshal(msg) //SerializeToOstream
    if err != nil{
    	Log("proto.Marshal failed")
    }
//    f.Write(buffer)
	b := Packet(buffer) 
	conn.Write(b)
}

func sendMessageToServer1(conn *net.TCPConn){
	for {
        var msg string
        fmt.Scan(&msg)
        if msg == "quit" {
            break
        }
		b := Packet([]byte (msg))
        conn.Write(b)    
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
//            	 Log(string(data))
				 parseData(data)
        }
    }
}

//解析数据
func parseData(msgbuf []byte ){
	 	msg := &protos.Helloworld{}
	    proto.Unmarshal(msgbuf, msg) //unSerialize
	   //CheckError(err)
	   Log("the msg is", msg)
}




func Log(v ...interface{}) {
    fmt.Println(v...)
}




