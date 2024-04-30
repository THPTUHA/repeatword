package crawl

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/PuerkitoBio/goquery"
	"github.com/THPTUHA/repeatword/audio"
	"github.com/THPTUHA/repeatword/config"
	"github.com/THPTUHA/repeatword/db"
	"github.com/THPTUHA/repeatword/vocab"
)

type CamCrawler struct {
	BaseCrawler
}

func NewCamCrawler() *CamCrawler {
	return &CamCrawler{BaseCrawler: BaseCrawler{
		BaseUrl:      "https://dictionary.cambridge.org/dictionary/english",
		AudioBaseUrl: "https://dictionary.cambridge.org",
	}}
}

func (cc *CamCrawler) Crawl(word string) error {
	cfg, err := config.Get()
	if err != nil {
		log.Fatal(err)
	}
	d, err := db.ConnectMysql()
	if err != nil {
		log.Fatal(err)
	}
	queries := db.New(d)
	ctx := context.Background()

	w, err := queries.GetWord(ctx, sql.NullString{String: word, Valid: true})
	if err != nil && err != sql.ErrNoRows {
		log.Fatal(err)
	}
	if w.ID != 0 {
		log.Fatalf("word %s exist", word)
	}
	doc, err := cc.fetchData(fmt.Sprintf("%s/%s", cc.BaseUrl, word))
	if err != nil {
		return err
	}

	vb := &vocab.Vocabulary{}
	vb.Word = sql.NullString{String: word, Valid: true}
	// parts
	doc.Find(".dictionary").Each(func(i int, s *goquery.Selection) {
		var part vocab.VobPart
		// header part
		headerEle := s.Find(".di-head").First()
		if headerEle != nil {
			part.Title.String = headerEle.Text()
		}

		part.Type.String = s.Find(".pos").Nodes[0].FirstChild.Data

		s.Find(".pos-header").Each(func(i int, s *goquery.Selection) {
			s.Find(".dpron-i").Each(func(i int, s *goquery.Selection) {
				var pros db.Pronounce
				pros.Region.String = s.Find(".region").Text()
				pros.AudioSrc.String, _ = s.Find("source").First().Attr("src")
				pros.Pro.String = s.Find(".pron").Text()
				part.Pronounces = append(part.Pronounces, &pros)
			})
		})

		// illus part
		s.Find(".def-block").Each(func(i int, s *goquery.Selection) {
			var mean vocab.Mean
			mean.Level.String = s.Find(".def-info").Text()
			mean.Meaning.String = s.Find(".ddef_d").Text()
			examples := make([]*db.Example, 0)
			s.Find(".examp").Each(func(i int, s *goquery.Selection) {
				examples = append(examples, &db.Example{Example: sql.NullString{String: s.Text(), Valid: true}})
			})
			mean.Examples = examples
			part.Means = append(part.Means, &mean)
		})

		vb.Parts = append(vb.Parts, &part)
	})

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	if len(vb.Parts) == 0 {
		log.Fatalf(fmt.Sprintf("Not found word %s", word))
	}
	for _, part := range vb.Parts {
		for idx, pros := range part.Pronounces {
			pros.LocalFile.String = fmt.Sprintf("%s_00%d.mp3", word, idx+1)

			err = audio.Download(
				fmt.Sprintf("%s%s", cc.AudioBaseUrl, pros.AudioSrc.String),
				path.Join(pwd, cfg.DataDir, pros.LocalFile.String),
			)
			if err != nil {
				panic(err)
			}
		}
	}

	jvb, err := json.Marshal(vb)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(jvb))

	err = queries.SetWord(ctx, db.SetWordParams{1, jvb})

	return nil
}
