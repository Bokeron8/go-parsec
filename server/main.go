package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"unsafe"

	"github.com/quic-go/quic-go"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	WIDTH  = 800
	HEIGHT = 600
    PIXEL_SIZE = 4 // RGB24 = 3 bytes por píxel
)

func main() {
	// quic-server-init-start
	listener, err := quic.ListenAddr("localhost:4242", generateTLSConfig(), nil)
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

			if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
                log.Fatal(err)
            }
            defer sdl.Quit()
        
            window, err := sdl.CreateWindow("Stream SDL2",
                sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
                WIDTH, HEIGHT, sdl.WINDOW_SHOWN)
            if err != nil {
                log.Fatal(err)
            }
            defer window.Destroy()
        
            renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
            if err != nil {
                log.Fatal(err)
            }
            defer renderer.Destroy()
        
            texture, _ := renderer.CreateTexture(
                sdl.PIXELFORMAT_BGRX8888, // O BGRX8888 si SDL lo requiere
                sdl.TEXTUREACCESS_STREAMING,
                WIDTH, HEIGHT)
            
            if err != nil {
                log.Fatal(err)
            }
            defer texture.Destroy()
        
            frameSize := WIDTH * HEIGHT * PIXEL_SIZE
            frame := make([]byte, frameSize)
        
            for {
                _, err := io.ReadFull(stream, frame)
                if err != nil {
                    log.Println("Error al leer frame:", err)
                    break
                }
        
                // Cargar los pixeles en la textura
                err = texture.Update(nil, unsafe.Pointer(&frame[0]), WIDTH*PIXEL_SIZE)
                if err != nil {
                    log.Println("Error al actualizar textura:", err)
                    break
                }

                for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
                    switch e := event.(type) {
                    case *sdl.QuitEvent:
                        log.Println("Cerrando ventana")
                        return
                    case *sdl.KeyboardEvent:
                        if e.Keysym.Sym == sdl.K_ESCAPE && e.State == sdl.PRESSED {
                            log.Println("ESC presionado. Saliendo.")
                            return
                        }
                    }
                }
                
        
                // Renderizar
                renderer.Clear()
                renderer.Copy(texture, nil, nil)
                renderer.Present()
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
