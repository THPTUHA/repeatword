package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/THPTUHA/repeatword/crawl"
	"github.com/THPTUHA/repeatword/game"
	"github.com/THPTUHA/repeatword/logger"
	"github.com/sirupsen/logrus"
)

func main() {
	cid := flag.Uint64("cid", 1, "collection id")
	limit := flag.Uint64("lm", 10, "limit question per")
	action := flag.String("action", "", "action to perform: play or crawl")
	flag.Parse()

	if *action == "" {
		fmt.Println("action is required")
		return
	}

	switch *action {
	case "play":
		game := game.Init(&game.Config{
			CollectionID: *cid,
			Limit:        *limit,
			Logger:       logger.InitLogger(logrus.DebugLevel.String()),
		})
		game.Play()
	case "crawl":
		if flag.NArg() < 1 {
			fmt.Println("need word")
			return
		}
		crawler := crawl.NewCamCrawler()
		err := crawler.Crawl(flag.Arg(0))
		if err != nil {
			log.Fatalln(err)
		}
	default:
		fmt.Println("unknown action")
	}
}
