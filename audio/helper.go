package audio

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
)

const (
	userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36"
)

func Download(url, destination string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Host", "dictionary.cambridge.org")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Accept", `text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7`)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download audio, status code: %d", resp.StatusCode)
	}

	err = os.MkdirAll(filepath.Dir(destination), 0755)
	if err != nil {
		return err
	}

	audioData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(destination, audioData, 0644)
	if err != nil {
		return err
	}

	fmt.Printf("Audio downloaded successfully to: %s\n", destination)
	return nil
}

var context *oto.Context

func PlayAudio(dataDir, filename string) {
	file, err := os.Open(path.Join(dataDir, filename))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	mp3Decoder, err := mp3.NewDecoder(file)
	if err != nil {
		log.Fatal(err)
	}

	if context == nil {
		context, err = oto.NewContext(mp3Decoder.SampleRate(), 2, 2, 8192)

		if err != nil {
			log.Fatal(err)
		}
	}

	player := context.NewPlayer()

	if _, err := io.Copy(player, mp3Decoder); err != nil {
		log.Fatal(err)
	}

	player.Close()
}
