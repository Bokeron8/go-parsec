package main

import (
	"context"
	"crypto/tls"
	"io"
	"log"
	"os/exec"

	"github.com/quic-go/quic-go"
)

func generateTLSConfig() *tls.Config {
    return &tls.Config{
        InsecureSkipVerify: true,
        ServerName: "localhost", // For testing purposes only; do not use in production
        NextProtos:         []string{"quic-0rtt-example"},
    }
}

func main() {

    session, err := quic.DialAddr(context.Background(), "localhost:4242", generateTLSConfig(), nil)
    if err != nil {
        log.Fatal(err)
    }
    
    stream, err := session.OpenStreamSync(context.Background())
    if err != nil {
        log.Fatal(err)
    }

    // 🔁 Captura pantalla con ffmpeg
	cmd := exec.Command("ffmpeg",
    "-f", "x11grab",                 // Linux
    "-video_size", "1366x768",
    "-i", ":0.0",                    // pantalla 0
    "-f", "mpegts",
    "-codec:v", "libx264",
    "-preset", "ultrafast",
    "-tune", "zerolatency",
    "-")

    stdout, _ := cmd.StdoutPipe()
    if err := cmd.Start(); err != nil {
        log.Fatal(err)
    }

    // 🚀 Enviar datos por QUIC
    _, err = io.Copy(stream, stdout)
    if err != nil {
        log.Fatal(err)
    }
    

    /*
    n, err := stream.Write([]byte("tu mmama trolo desde QUIC!"))
    if err != nil {
        log.Fatal(err)
    }

    

    println("ok")
    println(n)*/

    if err = stream.Close(); err != nil {
        log.Fatal(err)
    }
    

}
