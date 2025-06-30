package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"os/exec"

	"github.com/quic-go/quic-go"
)

func main() {
	// quic-server-init-start
	listener, err := quic.ListenAddr("192.168.0.41:4242", generateTLSConfig(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("QUIC server listening on localhost:4242")

	for {
		conn, err := listener.Accept(context.Background())
		if err != nil {
			log.Println("Accept error:", err)
			continue
		}
		go handleConn(conn)
	}
	// quic-server-init-end
}

func handleConn(conn quic.Connection) {
	// quic-server-handle-start
	defer conn.CloseWithError(0, "bye")

	for {
		stream, err := conn.AcceptStream(context.Background())
		if err != nil {
			return
		}

		go func(s quic.Stream) {
			defer s.Close()
            //ffplay -fflags nobuffer -flags low_delay -strict experimental -i -

			cmd := exec.Command("ffmpeg",
            "-fflags", "nobuffer",
            "-flags", "low_delay",
            "-i", "-",
            "-f", "sdl", "Video")


            stdin, _ := cmd.StdinPipe()
            if err := cmd.Start(); err != nil {
                log.Fatal(err)
            }

            _, err = io.Copy(stdin, s)
            if err != nil {
                log.Fatal(err)
            }
		}(stream)
	}
	// quic-server-handle-end
}

func generateTLSConfig() *tls.Config {
    cert, err := tls.LoadX509KeyPair("server.crt", "server.key")

	if err != nil {
		log.Fatal(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{"quic-0rtt-example"},
	}
}
