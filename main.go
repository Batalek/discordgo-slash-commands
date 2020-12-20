package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

var roles = map[string]string{
	"Role1": "420295797134065666",
	"Role2": "420295813064032288",
	"Role3": "420295828205469696",
	"Role4": "420295846127468547",
}
var guild = "344845336726077450"
var token = "<TOKEN>"

func main() {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}
	//Get AppID
	user, err := session.User("@me")
	if err != nil {
		panic(err)
	}
	appID := user.ID
	//Register event handler
	session.AddHandler(handler)
	err = session.Open()
	if err != nil {
		panic(err)
	}
	//Create choices from map
	choices := make([]CommandOptionChoice, len(roles))
	i := 0
	for k := range roles {
		choices[i] = CommandOptionChoice{
			Name:  k,
			Value: k,
		}
		i++
	}
	//Create command
	command := Command{
		Name:        "role",
		Description: "Assigns you a role",
		Options:     []CommandOption{
			{
				Type:        OptionTypeString,
				Name:        "role",
				Description: "This is the role you want to choose",
				Required:    true,
				Choices:     choices,
				Options:     nil,
			},
		},
	}
	//Register command for guild
	err = GuildCommandCreate(session, appID, guild, command)
	if err != nil {
		panic(err)
	}
	<- make(chan int)
}

func handler(session *discordgo.Session, event *discordgo.Event){
	if event.Type == "INTERACTION_CREATE" {
		//Parse event
		interaction, err := InteractionFromRaw(event.RawData)
		if err != nil {
			panic(err)
		}
		if interaction.Data.Name == "role" {
			selectedRole := interaction.Data.Options[0].Value.(string)
			roleID := roles[selectedRole]
			err = session.GuildMemberRoleAdd(interaction.GuildID, interaction.Member.User.ID, roleID)
			if err != nil {
				panic(err)
			}
			response := InteractionResponse{
				Type: InteractionResponseTypeChannelMessageWithSource,
				Data: InteractionResponseData{
					TTS:             false,
					Content:         fmt.Sprintf("You now have the %s role!", selectedRole),
				},
			}
			err = interaction.Respond(session, response)
			if err != nil {
				panic(err)
			}
		}
	}
}