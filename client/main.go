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
	/*cmd := exec.Command("ffmpeg",
		"-f", "x11grab",         // Captura pantalla (Linux)
		"-video_size", "1366x768",
        "-framerate", "30",
		"-i", ":0.0",            // Pantalla 0
        "-vf", "scale=200:100",
		"-vcodec", "rawvideo",   // Sin compresión
		"-pix_fmt", "rgb24",     // 3 bytes por píxel
		"-f", "rawvideo",        // Salida raw
		"-")*/

        cmd := exec.Command("ffmpeg",
		"-f", "kmsgrab",
		"-device", "/dev/dri/card1",
		"-framerate", "30",
		"-i", "-",
		"-vf", "hwmap=derive_device=vaapi,hwdownload,format=bgr0,scale=800x600",
		"-f", "rawvideo",
		"-",
	)
        
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
