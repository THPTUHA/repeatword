package vocab

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/THPTUHA/repeatword/db"
)

// map db model
type Mean struct {
	db.Mean
	Examples []*db.Example
}
type VobPart struct {
	db.VobPart
	Means      []*Mean
	Pronounces []*db.Pronounce
}
type Vocabulary struct {
	db.Vob
	Correct bool
	Parts   []*VobPart
}

type PronouncesJson struct {
	AudioSrc  string `json:"audio_src"`
	LocalFile string `json:"local_file"`
	Region    string `json:"region"`
	Pro       string `json:"pro"`
}

type MeanJson struct {
	Level    string   `json:"level"`
	Meaning  string   `json:"meaning"`
	Examples []string `json:"examples"`
}

type VobPartJson struct {
	Means      []*MeanJson       `json:"means"`
	Pronounces []*PronouncesJson `json:"pronounces"`
	Type       string            `json:"type"`
	Title      string            `json:"title"`
}

type VobJson struct {
	ID    int            `json:"id"`
	Parts []*VobPartJson `json:"parts"`
	Word  string         `json:"word"`
}

func (v *Vocabulary) MarshalJSON() ([]byte, error) {

	parts := make([]*VobPartJson, len(v.Parts))

	for i, part := range v.Parts {
		parts[i] = &VobPartJson{
			Type:  part.Type.String,
			Title: part.Title.String,
		}
		parts[i].Means = make([]*MeanJson, len(part.Means))
		for j, mean := range part.Means {
			parts[i].Means[j] = &MeanJson{
				Level:   mean.Level.String,
				Meaning: mean.Meaning.String,
			}
			parts[i].Means[j].Examples = make([]string, len(part.Means[j].Examples))
			for k, example := range part.Means[j].Examples {
				parts[i].Means[j].Examples[k] = example.Example.String
			}
		}
		parts[i].Pronounces = make([]*PronouncesJson, len(part.Pronounces))
		for j, pro := range part.Pronounces {
			parts[i].Pronounces[j] = &PronouncesJson{
				AudioSrc:  pro.AudioSrc.String,
				LocalFile: pro.LocalFile.String,
				Region:    pro.Region.String,
				Pro:       pro.Pro.String,
			}
		}
	}

	data := struct {
		Word  string         `json:"word"`
		Parts []*VobPartJson `json:"parts"`
	}{
		Word:  v.Word.String,
		Parts: parts,
	}

	return json.Marshal(&data)
}

func (v *Vocabulary) UnmarshalJSON(data []byte) error {
	var vob VobJson
	if err := json.Unmarshal(data, &vob); err != nil {
		return err
	}
	v.Parts = make([]*VobPart, len(vob.Parts))
	v.Word = sql.NullString{String: vob.Word, Valid: true}
	for i, part := range vob.Parts {
		v.Parts[i] = &VobPart{
			VobPart: db.VobPart{
				Type:  sql.NullString{String: part.Type, Valid: true},
				Title: sql.NullString{String: part.Title, Valid: true},
			},
		}
		v.Parts[i].Means = make([]*Mean, len(part.Means))
		for j, mean := range part.Means {
			v.Parts[i].Means[j] = &Mean{
				Mean: db.Mean{
					Meaning: sql.NullString{String: mean.Meaning, Valid: true},
					Level:   sql.NullString{String: mean.Level, Valid: true},
				},
			}
			v.Parts[i].Means[j].Examples = make([]*db.Example, len(mean.Examples))
			for k, example := range mean.Examples {
				v.Parts[i].Means[j].Examples[k] = &db.Example{
					Example: sql.NullString{String: example, Valid: true},
				}
			}
		}

		v.Parts[i].Pronounces = make([]*db.Pronounce, len(part.Pronounces))
		for j, pro := range part.Pronounces {
			v.Parts[i].Pronounces[j] = &db.Pronounce{
				AudioSrc:  sql.NullString{String: pro.AudioSrc, Valid: true},
				LocalFile: sql.NullString{String: pro.LocalFile, Valid: true},
				Region:    sql.NullString{String: pro.Region, Valid: true},
				Pro:       sql.NullString{String: pro.Pro, Valid: true},
			}
		}
	}
	return nil
}

func (v *Vocabulary) String() string {
	result := fmt.Sprintf("Word: %s\n", v.Word.String)

	for _, part := range v.Parts {
		result += fmt.Sprintf("  Title: %s\n", part.Title.String)
		result += fmt.Sprintf("  Type: %s\n", part.Type.String)

		for _, pro := range part.Pronounces {
			result += fmt.Sprintf("    Region: %s\n", pro.Region.String)
			result += fmt.Sprintf("    AudioSrc: %s\n", pro.AudioSrc.String)
			result += fmt.Sprintf("    Pro: %s\n", pro.Pro.String)
			result += fmt.Sprintf("    LocalFile: %s\n", pro.LocalFile.String)
		}

		for _, mean := range part.Means {
			result += fmt.Sprintf("    Meaning: %s\n", mean.Mean.Meaning.String)
			for _, example := range mean.Examples {
				result += fmt.Sprintf("    Example: %s\n", example.Example.String)
			}
		}
	}

	return result
}
