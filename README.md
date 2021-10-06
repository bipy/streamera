# streamera

**Term Project of Computer Networking**

**streamera** is a Stream Camera based on TCP, which contains **client mode** and **server mode**.

## Features

- **Client Mode & Server Mode**

    - **Client Mode**

        Capture camera video and send stream to server side.

    - **Server Mode** 

        Receive tcp stream and display the video to screen, along with network speed and ping.

- **Connection-Error Handle**

    - **Client Mode**

        streamera client will try to reconnect to the server every 3 second when TCP connection is lost. And the maximum bufferd size is 180 frames, which is probably up to 6 seconds of the video.

    - **Server Mode**

        streamera server will continuously serve the incoming tcp connection, but the limit is set up to 1 client at the same time, which means **only one** video window will be shown during streamera server running. 

        To be simple, When current TCP connnection is lost, the server will wait for the next  incoming connection and start over the video streaming.

- **Speed & Ping Monitoring**

    Draw current **network speed** and **ping** directly to the video frame.

- **Thread Safe**

    streamera uses golang as the main programming language: goroutine, channel and mutex are widely used in streamera. Thus streamera has high profermance and stability while keeping thread safe.

## How to Install

streamera requires [OpenCV](https://opencv.org/) as computer vision library and [GoCV](https://github.com/hybridgroup/gocv) as golang bindings.

### MacOS

Install **OpenCV**, **pkg-config**

```bash
brew install opencv pkg-config
```

Install **GoCV**

```bash
go get -u -d gocv.io/x/gocv
```

Done!

### Linux or Windows

Please check [Getting Start :: GoCV](https://gocv.io/getting-started/) for more information.



## Usage

Clone this repo:

```bash
git clone git@github.com:bipy/streamera.git
```

Build streamera:

```bash
cd streamera
go build
```

To run as server mode:

```bash
./streamera -m server
```

To run as client mode:

```bash
./streamera -m client
```




