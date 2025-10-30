package cmds

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func Schedule(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Access options in the order provided by the user.
	options := i.ApplicationCommandData().Options

	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	var schedule string                            // The day's schedule as a string
	var dayTypeFormal string                       // The day type as a pretty, human friendly string
	dayType := optionMap["day-type"].StringValue() // The raw day type

	// The basic format for a schedule
	dayLayout := `P1 - %s
P2 - %s
Lunch - %s
P3 - %s
P4 - %s`

	switch dayType {
	case "normal":
		dayTypeFormal = "Normal"
		schedule = fmt.Sprintf(
			dayLayout,
			"9:00-10:20",
			"10:25-11:40",
			"11:40-12:40",
			"12:40-1:55",
			"2:00-3:15",
		)
	case "lateStart":
		dayTypeFormal = "Late Start"
		schedule = fmt.Sprintf(
			dayLayout,
			"9:55-11:05",
			"11:10-12:10",
			"12:10-1:10",
			"1:10-2:10",
			"2:15-3:15",
		)
	case "halfDay":
		dayTypeFormal = "Half Day"
		schedule = fmt.Sprintf(
			"P1 - %s\nP2 - %s\nP3 - %s\nP4 - %s",
			"9:00-9:45",
			"9:50-10:35",
			"10:40-11:25",
			"11:30-12:15",
		)
	case "assembly":
		dayTypeFormal = "Assembly"
		schedule = fmt.Sprintf(
			dayLayout,
			"9:00-11:00",
			"11:05-12:10",
			"12:10-1:10",
			"1:10-2:10",
			"2:15-3:15",
		)
	default:
		dayTypeFormal = "INVALID"
		schedule = "There is no such type of day..."
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       fmt.Sprintf("%s Schedule", dayTypeFormal),
					Description: schedule,
					Color:       0x1abc9c,
				},
			},
		},
	})
}
