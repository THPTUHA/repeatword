package crawl

import (
	"os"
	"path"

	"github.com/THPTUHA/repeatword/audio"
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
	// doc, err := cc.fetchData(fmt.Sprintf("%s/%s", cc.BaseUrl, word))
	// if err != nil {
	// 	return err
	// }

	// var vb vocab.Vocabulary
	// vb.Word = word
	// // parts
	// doc.Find(".pr .dictionary").Each(func(i int, s *goquery.Selection) {
	// 	fmt.Println("Run here")
	// 	var part vocab.VocabPart
	// 	// header part
	// 	headerEle := s.Find(".di-head").First()
	// 	if headerEle != nil {
	// 		part.Header = headerEle.Text()
	// 	}

	// 	part.Type = s.Find(".posgram").Text()

	// 	s.Find(".dpron-i").Each(func(i int, s *goquery.Selection) {
	// 		var pros vocab.Pronounce
	// 		pros.Region = s.Find(".region").Text()
	// 		pros.AudioSrc, _ = s.Find("source").First().Attr("src")
	// 		pros.Pro = s.Find(".pron").Text()
	// 		part.Pronounces = append(part.Pronounces, &pros)
	// 	})

	// 	// illus part
	// 	s.Find(".dsense").Each(func(i int, s *goquery.Selection) {
	// 		var ill vocab.VocabIllus
	// 		ill.Mean = s.Find(".ddef_d").Text()
	// 		examples := make([]string, 0)
	// 		s.Find(".examp").Each(func(i int, s *goquery.Selection) {
	// 			examples = append(examples, s.Text())
	// 		})
	// 		ill.Examples = examples
	// 		part.Illustration = append(part.Illustration, &ill)
	// 	})

	// 	vb.Parts = append(vb.Parts, &part)
	// })

	// fmt.Println(vb.String())
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	err = audio.Download(
		"https:/dictionary.cambridge.org/media/english/uk_pron/u/uks/uksub/uksubsp004.mp3",
		path.Join(pwd, "data"),
	)
	if err != nil {
		panic(err)
	}
	return nil
}
