package main

import (
	"flag"
	"fmt"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/auth"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"image"
	"image/png"
	"os"
)

var (
	ip   string
	port string
)

// ty TwistedAsylum in the gophertunnel discord
func SkinToRGBA(s protocol.Skin) *image.RGBA {
	t := image.NewRGBA(image.Rect(0, 0, int(s.SkinImageWidth), int(s.SkinImageHeight)))
	t.Pix = s.SkinData
	return t
}

func CapeToRGBA(s protocol.Skin) *image.RGBA {
	t := image.NewRGBA(image.Rect(0, 0, int(s.CapeImageWidth), int(s.CapeImageHeight)))
	t.Pix = s.CapeData
	return t
}

func main() {
	flag.StringVar(&ip, "ip", "127.0.0.1", "Servers IP Address")
	flag.StringVar(&port, "port", "19132", "Servers Port")
	flag.Parse()
	_ = os.Mkdir("stolen", 0755)

	dialer := minecraft.Dialer {
		TokenSource: auth.TokenSource,
	}

	address := ip + ":" + port
	fmt.Println("Connecting to " + address)
	conn, err := dialer.Dial("raknet", address)
	if err != nil {
		panic(err)
	}

	defer conn.Close()
	_ = conn.DoSpawn()

	for {
		pk, err := conn.ReadPacket()
		if err != nil {
			break
		}

		switch p := pk.(type) {
		case *packet.PlayerList:
			go func() {
				for _, player := range p.Entries {
					name := player.Username
					skin := SkinToRGBA(player.Skin)
					cape := CapeToRGBA(player.Skin)
					path, _ := os.Getwd()
					_ = os.Mkdir(fmt.Sprintf("%s\\stolen\\%s", path, name), 0755)
					fileSkin, _ := os.Create(fmt.Sprintf("%s\\stolen\\%s\\skin.png", path, name))
					fileCape, _ := os.Create(fmt.Sprintf("%s\\stolen\\%s\\cape.png", path, name))
					_ = png.Encode(fileSkin, skin)
					_ = png.Encode(fileCape, cape)
					fileSkin.Close()
					fileCape.Close()
					fmt.Println("Stolen " + name)
				}
			}()
		}
	}
}