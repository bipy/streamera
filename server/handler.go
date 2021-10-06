package server

import (
    "bufio"
    "encoding/binary"
    "fmt"
    "io"
    "net"
    "streamera/common"
    "time"
)

func runListener(server *Server) {
    for {
        server.TCPListenerMutex.Lock()
        server.FrameMutex.Lock()
        conn, err := server.TCPListener.Accept()
        if err != nil {
            fmt.Println(common.Red("broken connection " + err.Error()))
            continue
        }
        go handleConn(server, conn)
    }
}

func handleConn(server *Server, conn net.Conn) {
    defer func(server *Server, conn net.Conn) {
        err := conn.Close()
        if err != nil {
            fmt.Println(common.Red("tcp conn close failed! " + err.Error()))
        }
    }(server, conn)

    defer server.TCPListenerMutex.Unlock()

    server.FrameMutex.Unlock()

    reader := bufio.NewReader(conn)

    timeStamp := make([]byte, common.TimeStampPacketSize)

    for {
        _, err := io.ReadFull(reader, timeStamp)
        if err != nil {
            fmt.Println(common.Red(conn.RemoteAddr().String() + " Down! " + err.Error()))
            break
        }
        curTime := time.Now().UnixMicro()
        recvTime := int64(binary.LittleEndian.Uint64(timeStamp))

        server.FrameMutex.Lock()
        _, err = io.ReadFull(reader, server.Frame)
        if err != nil {
            fmt.Println(common.Red(conn.RemoteAddr().String() + " Down! " + err.Error()))
            break
        }
        server.FrameMutex.Unlock()

        server.Counter.Mutex.Lock()
        server.Counter.PingRealTime = curTime - recvTime
        server.Counter.PktCount++
        server.Counter.Mutex.Unlock()
    }
}