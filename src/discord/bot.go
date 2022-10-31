package discord

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/kompocikdot/labelbuddy/src/restocks"
	"github.com/kompocikdot/labelbuddy/src/utils"
)

var (
	toShipcommand = &discordgo.ApplicationCommand{
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

	soldCommand = &discordgo.ApplicationCommand{
		Name: "restocks-get-sold-items",
		Description: "Scrapes sold items",
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

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"restocks-get-sold-items": func(s *discordgo.Session, i *discordgo.InteractionCreate){
			options := i.ApplicationCommandData().Options
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}
	
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags: 1 << 6,
					Embeds: []*discordgo.MessageEmbed{
						{
							Description: "Fetching results...",
						},
					},
				},
			})
	
			items, _ := restocks.RetrieveItemsPayments(
				optionMap["email"].StringValue(), optionMap["password"].StringValue(),
			)
	
			csvLikeString := utils.GenerateCSV(items)
	
			s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Embeds: &[]*discordgo.MessageEmbed{
					{
						Description: "Results fetched.",
					},
				},
				Files: []*discordgo.File{
					{
						Name: "results.csv",
						ContentType: "text/csv",
						Reader: strings.NewReader(csvLikeString),
					},
				},
			})
		},

		"get-restocks-labels": func(s *discordgo.Session, i *discordgo.InteractionCreate){
			options := i.ApplicationCommandData().Options
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}
	
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags: 1 << 6,
					Embeds: []*discordgo.MessageEmbed{
						{
							Description: "Fetching results...",
						},
					},
				},
			})

			items, fileNames, pdfFilename, _ := restocks.RetrieveSalesLinks(
				optionMap["email"].StringValue(), optionMap["password"].StringValue(),
			)
			if len(items) == 0 {
				return
			}

			pdfFile, err := os.Open("generated/" + pdfFilename)
			if err != nil {
				fmt.Println(err)
			}

			var ItemEmbeds = []*discordgo.MessageEmbedField{}
			for _, item := range items {
				ItemEmbeds = append(ItemEmbeds, &discordgo.MessageEmbedField{
					Name: item.ItemName,
					Value: fmt.Sprintf("%s, %s, %s", item.Id, item.ItemSize, item.ShipTo),
				})
			}

			_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Embeds: &[]*discordgo.MessageEmbed{
					{
						Description: "Results fetched.",
						Fields: ItemEmbeds,
					},
				},
				Files: []*discordgo.File{
					{
						Name: pdfFilename,
						ContentType: "application/pdf",
						Reader: pdfFile,
					},
				},
			})
			if err != nil {
				fmt.Println(err)
			}
			restocks.ClearFiles(fileNames)
		},
	}
)


func RunBot(token string) {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("error creating Discord session,", err)
	}

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection:", err)
		return
	}

	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	dg.ApplicationCommandCreate(dg.State.User.ID, "", toShipcommand)
	dg.ApplicationCommandCreate(dg.State.User.ID, "", soldCommand)

	dg.Identify.Intents = discordgo.IntentsGuildMessages

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()
}
