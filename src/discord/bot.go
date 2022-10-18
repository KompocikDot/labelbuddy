package discord

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/kompocikdot/labelbuddy/src/restocks"
)


func RunBot(token string) {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("error creating Discord session,", err)
	}


	command := &discordgo.ApplicationCommand{
		Name: "get-restocks-labels",
		Description: "Scrapes your restocks labels",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "email",
				Description: "Restocks account email",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "password",
				Description: "Restocks account password",
				Required:    true,
			},
		},
	}

	restocksLabelFunc := func(s *discordgo.Session, i *discordgo.InteractionCreate){
		options := i.ApplicationCommandData().Options
		optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
		for _, opt := range options {
			optionMap[opt.Name] = opt
		}

		items, err := restocks.RetrieveSalesLinks(
			optionMap["email"].StringValue(), optionMap["password"].StringValue(),
		)

		fields := []*discordgo.MessageEmbedField{}
		if err != nil {
			fields = append(fields, &discordgo.MessageEmbedField{
				Name: "Error",
				Value: err.Error(),
			})
		} else {
			if len(items) > 0 {
				for _, item := range items {
					singleField := &discordgo.MessageEmbedField{
						Name: fmt.Sprintf("Ship to: %s", item.ShipTo),
						Value: item.Link,
					}
					fields = append(fields, singleField)
				}
			} else {
				field := &discordgo.MessageEmbedField{
					Name: "No labels found",
					Value: "None",
				}
				fields = append(fields, field)
			}
		}



		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags: 1 << 6,
				Embeds: []*discordgo.MessageEmbed{
					{
						Description: "Make sure you're logged in before downloading labels.",
						Fields: fields,
					},
				},
			},
		})
		if err != nil {
			log.Panic(err.Error())
		}
	}



	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		restocksLabelFunc(s, i)
	})
	
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection:", err)
		return
	}

	_, err = dg.ApplicationCommandCreate(dg.State.User.ID, "", command)
	if err != nil {
		log.Panicf("Cannot create '%v' command: %v", command.Name, err)
	}

	dg.Identify.Intents = discordgo.IntentsGuildMessages

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()
}
