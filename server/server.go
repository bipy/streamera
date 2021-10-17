package server

import (
    "bytes"
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
    Frame            *bytes.Buffer
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
    buf := new(bytes.Buffer)

    server := &Server{
        TCPListener:      tcpListener,
        TCPListenerMutex: sync.Mutex{},
        Frame:            buf,
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
    go updateLatency(server.Counter)

    window := gocv.NewWindow("Streaming")
    for {
        server.FrameMutex.RLock()
        mat, err := gocv.IMDecode(server.Frame.Bytes(), gocv.IMReadUnchanged)
        if err != nil {
            server.FrameMutex.RUnlock()
            continue
        }
        server.FrameMutex.RUnlock()
        server.Counter.Mutex.RLock()
        curSpeed := server.Counter.SpeedPerSecond
        curLatency := server.Counter.LatencyPerSecond
        server.Counter.Mutex.RUnlock()
        gocv.PutText(
            &mat,
            fmt.Sprintf("Speed: %s", getHumanReadableSpeed(curSpeed)),
            image.Point{X: 50, Y: 50},
            gocv.FontHersheySimplex,
            1,
            color.RGBA{R: 255, G: 0, B: 0, A: 0},
            3,
        )
        gocv.PutText(
            &mat,
            fmt.Sprintf("Latency: %s", getHumanReadableTime(curLatency)),
            image.Point{X: 50, Y: 100},
            gocv.FontHersheySimplex,
            1,
            color.RGBA{R: 255, G: 0, B: 0, A: 0},
            3,
        )
        window.IMShow(mat)
        window.WaitKey(1)
        _ = mat.Close()
    }
}
