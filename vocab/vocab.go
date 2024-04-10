package vocab

import "fmt"

type Pronounce struct {
	Region   string
	AudioSrc string
	Pro      string
}

type VocabIllus struct {
	Mean     string
	Examples []string
}

type VocabPart struct {
	Header       string
	Type         string
	Pronounces   []*Pronounce
	Illustration []*VocabIllus
}

type Vocabulary struct {
	Word  string
	Parts []*VocabPart
}

func (v *Vocabulary) String() string {
	result := fmt.Sprintf("Word: %s\n", v.Word)

	for _, part := range v.Parts {
		result += fmt.Sprintf("  Header: %s\n", part.Header)
		result += fmt.Sprintf("  Type: %s\n", part.Type)

		for _, pronounce := range part.Pronounces {
			result += fmt.Sprintf("    Region: %s\n", pronounce.Region)
			result += fmt.Sprintf("    AudioSrc: %s\n", pronounce.AudioSrc)
			result += fmt.Sprintf("    Pro: %s\n", pronounce.Pro)
		}

		for _, illus := range part.Illustration {
			result += fmt.Sprintf("    Mean: %s\n", illus.Mean)
			for _, example := range illus.Examples {
				result += fmt.Sprintf("    Example: %s\n", example)
			}
		}
	}

	return result
}
