package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/THPTUHA/repeatword/config"
	"github.com/THPTUHA/repeatword/crawl"
	"github.com/THPTUHA/repeatword/game"
	"github.com/THPTUHA/repeatword/logger"
	"github.com/THPTUHA/repeatword/setup"
	"github.com/sirupsen/logrus"
)

func main() {
	env := os.Getenv("ENV")
	var cfg *config.Configs

	if env == "dev" {
		ex, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		cfg, err = config.Set(path.Join(ex, "config.dev.yaml"))
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		ex, err := os.Executable()
		if err != nil {
			panic(err)
		}
		exPath := filepath.Dir(ex)

		cfg, err = config.Set(path.Join(exPath, "config.yaml"))
		if err != nil {
			log.Fatalln(err)
		}
	}

	cid := flag.Uint64("cid", 1, "collection id")
	limit := flag.Uint64("lm", 10, "limit question per")
	mode := flag.Uint("mode", 0, "Mode game")
	action := flag.String("action", "", "action to perform: play or crawl")
	recentDay := flag.Int("rd", -1, "recent day")
	flag.Parse()

	if *action == "" {
		fmt.Println("action is required")
		return
	}

	switch *action {
	case "setup":
		if err := setup.Setup(); err != nil {
			log.Fatalln(err)
		}
	case "play":
		game := game.Init(&game.Config{
			CollectionID: *cid,
			Limit:        *limit,
			Logger:       logger.InitLogger(logrus.DebugLevel.String()),
			RecentDayNum: *recentDay,
			Root:         cfg,
			Mode:         *mode,
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
