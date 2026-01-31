//go:build !tcp
// +build !tcp

package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ninjatech/tesselbox-go/pkg/render"
)

func main() {
	game := render.NewGame()

	ebiten.SetWindowSize(render.ScreenWidth, render.ScreenHeight)
	ebiten.SetWindowTitle("Tesselbox - Go/Ebiten Version")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}