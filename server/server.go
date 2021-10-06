package server

import (
    "fmt"
    "gocv.io/x/gocv"
    "image"
    "image/color"
    "net"
    "streamera/common"
    "sync"
)

type Server struct {
    TCPListener      net.Listener
    TCPListenerMutex sync.Mutex
    Frame            []byte
    FrameMutex       sync.RWMutex
    Counter          *SpeedCounter
}

func NewServer(ip net.IP, port int) (*Server, error) {
    tcpListener, err := net.ListenTCP("tcp", &net.TCPAddr{
        IP:   ip,
        Port: port,
    })
    if err != nil {
        fmt.Println(common.Red("generate tcp server failed!"))
        return nil, err
    }

    counter := NewSpeedCounter()

    server := &Server{
        TCPListener:      tcpListener,
        TCPListenerMutex: sync.Mutex{},
        Frame:            make([]byte, common.FramePacketSize),
        FrameMutex:       sync.RWMutex{},
        Counter:          counter,
    }

    return server, nil
}

func RunServer(server *Server) {
    defer func(tcpListener net.Listener) {
        err := tcpListener.Close()
        if err != nil {
            panic(err)
        }
    }(server.TCPListener)

    go runListener(server)
    go calcSpeed(server.Counter)
    go updatePing(server.Counter)

    window := gocv.NewWindow("Streaming")
    for {
        server.FrameMutex.RLock()
        img, err := gocv.NewMatFromBytes(common.Height, common.Width, gocv.MatTypeCV8UC3, server.Frame)
        if err != nil {
            fmt.Println(common.Red(err.Error()))
            server.FrameMutex.RUnlock()
            continue
        }
        server.Counter.Mutex.RLock()
        curSpeed := server.Counter.SpeedPerSecond
        curPing := server.Counter.PingPerSecond
        server.Counter.Mutex.RUnlock()
        gocv.PutText(
            &img,
            fmt.Sprintf("Speed: %s", getHumanReadableSpeed(curSpeed)),
            image.Point{X: 50, Y: 50},
            gocv.FontHersheySimplex,
            1,
            color.RGBA{R: 255, G: 0, B: 0, A: 0},
            3,
        )
        gocv.PutText(
            &img,
            fmt.Sprintf("Ping: %s", getHumanReadablePing(curPing)),
            image.Point{X: 50, Y: 100},
            gocv.FontHersheySimplex,
            1,
            color.RGBA{R: 255, G: 0, B: 0, A: 0},
            3,
        )
        window.IMShow(img)
        window.WaitKey(1)
        _ = img.Close()
        server.FrameMutex.RUnlock()
    }
}
