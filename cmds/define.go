package cmds

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bwmarrin/discordgo"
)

type Definition struct {
	Word      string `json:"word"`
	Phonetic  string `json:"phonetic"`
	Phonetics []struct {
		Text      string `json:"text"`
		Audio     string `json:"audio"`
		SourceURL string `json:"sourceUrl,omitempty"`
		License   struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"license"`
	} `json:"phonetics"`
	Meanings []struct {
		PartOfSpeech string `json:"partOfSpeech"`
		Definitions  []struct {
			Definition string `json:"definition"`
			Synonyms   []any  `json:"synonyms"`
			Antonyms   []any  `json:"antonyms"`
			Example    string `json:"example,omitempty"`
		} `json:"definitions"`
		Synonyms []string `json:"synonyms"`
		Antonyms []any    `json:"antonyms"`
	} `json:"meanings"`
	License struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"license"`
	SourceUrls []string `json:"sourceUrls"`
}

func Define(s *discordgo.Session, i *discordgo.InteractionCreate) {
	word := i.ApplicationCommandData().Options[0].StringValue()

	res, err := http.Get("https://api.dictionaryapi.dev/api/v2/entries/en/" + word)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "Failed to define the word :(",
						Description: err.Error(),
						Color:       0xff0000,
					},
				},
			},
		})

		return
	}

	definitions := []Definition{}
	defer res.Body.Close()
	json.NewDecoder(res.Body).Decode(&definitions)
	definition := definitions[0]

	prettyDefinitions := []*discordgo.MessageEmbedField{}

	for _, meaning := range definition.Meanings {
		prettyDefinitionsForPartOfSpeech := ""
		for _, def := range meaning.Definitions {
			prettyDefinitionsForPartOfSpeech += fmt.Sprintf("  %s\n", def.Definition)
		}
		prettyDefinitions = append(
			prettyDefinitions,
			&discordgo.MessageEmbedField{
				Name:  meaning.PartOfSpeech,
				Value: prettyDefinitionsForPartOfSpeech,
			},
		)

	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:  fmt.Sprintf("%s - %s", definition.Word, definition.Phonetic),
					Fields: prettyDefinitions,
					Color:  0x1abc9c,
				},
			},
		},
	})

}
