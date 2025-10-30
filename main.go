package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

// Bot parameters
var (
	GuildID        = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	BotToken       = flag.String("token", "", "Bot access token")
	RemoveCommands = flag.Bool("rmcmd", false, "Remove all commands after shutdowning or not")
)

var s *discordgo.Session

func init() {
	flag.Parse()

	var err error
	s, err = discordgo.New("Bot " + *BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "ping",
			Description: "Check that the bot is online",
		},
		{
			Name:        "schedule",
			Description: "Generate a schedule based on the type of day it is",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "day-type",
					Description: "The type of day: normal, lateStart, halfDay, assembly",
					Required:    true,
				},
			},
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"ping": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Pong! The bot is, in fact, online",
				},
			})
		},
		"schedule": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
				// Ignore type for now, they will be discussed in "responses"
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       fmt.Sprintf("%s Schedule", dayTypeFormal),
							Description: schedule,
						},
					},
				},
			})
		},
	}
)

func init() {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func main() {
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, *GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	if *RemoveCommands {
		log.Println("Removing commands...")

		for _, v := range registeredCommands {
			err := s.ApplicationCommandDelete(s.State.User.ID, *GuildID, v.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}

	log.Println("Gracefully shutting down.")
}
