package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jroimartin/gocui"
)

// Todo: move start time to a saved file, so we can resume existing games
// Todo: move game board creation to init
// Todo: enable the runtime clock on the bottom

func main() {
	logfile, err := os.OpenFile("runtime.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		panic(err)
	}
	defer logfile.Close()
	log.SetOutput(logfile)
	log.Println("Starting Game")

	game := Game{}
	game.Start()

	g := gocui.NewGui()
	if err := g.Init(); err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlQ, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	g.SetLayout(introLayout)

	//runtimeChannel := make(chan int64)
	//go updateRuntimeChannel(game, runtimeChannel)
	//go showRuntime(g, runtimeChannel)
	go removeSplash(g)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

type Game struct {
	StartTime time.Time
}

func currentUnixTime() int64 {
	return time.Now().Unix()
}

func (g *Game) Start() {
	g.StartTime = time.Now()
}

func (g Game) Runtime() int64 {
	return time.Now().Unix() - g.StartTime.Unix()
}
func removeSplash(g *gocui.Gui) {
	time.Sleep(1000 * time.Millisecond)
	g.Execute(func(g *gocui.Gui) error {
		log.Println("Removing the Splash Screen")
		if err := g.DeleteView("hello"); err != nil {
			log.Println(err)
		}
		return nil
	})
}

//func updateRuntimeChannel(game Game, r chan int64) {
//  for {
//    log.Println("Game Runtime:", game.Runtime())
//    r <- game.Runtime()
//    time.Sleep(1000 * time.Millisecond)
//  }
//}

//func showRuntime(g *gocui.Gui, r chan int64) {
//  select {
//  case runtime := <-r:
//    g.Execute(func(g *gocui.Gui) error {
//      log.Println("I was executed")
//      v, err := g.View("hello")
//      if err != nil {
//        log.Panicln(err)
//      }
//      greeting := "Game has been running for"
//      fmt.Fprintln(v, greeting, runtime)
//      return nil
//    })
//  default:
//    // don't block
//  }
//}

func introLayout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	greeting := "Welcome to MLG Idler"
	greetingLength := len(greeting) / 2
	if v, err := g.SetView("hello", maxX/2-greetingLength-1, maxY/2, maxX/2+greetingLength, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, greeting)
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	log.Println("Normal Quit")
	return gocui.ErrQuit
}
