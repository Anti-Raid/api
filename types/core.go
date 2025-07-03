package types

import (
	"time"

	"github.com/anti-raid/eureka/dovewing/dovetypes"
	"github.com/bwmarrin/discordgo"
)

// API configuration data
type ApiConfig struct {
	MainServer          string `json:"main_server" description:"The ID of the main Anti-Raid Discord Server" validate:"required"`
	SupportServerInvite string `json:"support_server_invite" comment:"Discord Support Server Link" default:"https://discord.gg/u78NFAXm" validate:"required"`
	ClientID            string `json:"client_id" description:"The ID of the Anti-Raid bot client" validate:"required"`
}

// This represents a IBL Popplio API Error
type ApiError struct {
	Context map[string]string `json:"context,omitempty" description:"Context of the error. Usually used for validation error contexts"`
	Message string            `json:"message" description:"Message of the error"`
}

type BotState struct {
	Commands []discordgo.ApplicationCommand `json:"commands"`
}

type DashboardGuild struct {
	ID          string `json:"id" description:"The ID of the guild"`
	Name        string `json:"name" description:"The name of the guild"`
	Avatar      string `json:"avatar" description:"The avatar url of the guild"`
	Permissions int64  `json:"permissions" description:"The permissions the user has in the guild"`
}

type DashboardGuildData struct {
	Guilds        []*DashboardGuild `json:"guilds" description:"The guilds the user is in"`
	BotInGuilds   []string          `json:"has_bot" description:"A list of guild IDs that the user has the bot in"`
	UnknownGuilds []string          `json:"unknown_guilds" description:"A list of guild IDs where the bot is in an outage etc. in"`
}

// Represents a user on Antiraid
type User struct {
	User       *dovetypes.PlatformUser `json:"user" description:"The user object of the user"`
	State      string                  `db:"state" json:"state" description:"The state of the user"`
	VoteBanned bool                    `db:"vote_banned" json:"vote_banned" description:"Whether or not the user is banned from recieving rewards from voting"`
	CreatedAt  time.Time               `db:"created_at" json:"created_at" description:"The time the user was created"`
	UpdatedAt  time.Time               `db:"updated_at" json:"updated_at" description:"The time the user was last updated"`
}

type UserGuildBaseData struct {
	OwnerID   string                        `json:"owner_id" description:"The ID of the guild owner"`
	Name      string                        `json:"name" description:"The name of the guild"`
	Icon      *string                       `json:"icon" description:"The icon of the guild"`
	Roles     []SerenityRole                `json:"roles" description:"The roles of the guild"`
	UserRoles []string                      `json:"user_roles" description:"The role IDs the user has in the guild"`
	BotRoles  []string                      `json:"bot_roles" description:"The role IDs the bot has in the guild"`
	Channels  []GuildChannelWithPermissions `json:"channels" description:"The channels of the guild with permission info"`
}

// SettingsExecute allows execution of a settings operation
type SettingsExecute struct {
	Operation string         `json:"operation" description:"The operation type to execute"`
	Setting   string         `json:"setting" description:"The name of the setting"`
	Fields    map[string]any `json:"fields" description:"The fields to execute the operation with"`
}

type DispatchResult struct {
	Type string `json:"type" description:"The type of the dispatch result [Ok or Err]"`
	Data any    `json:"data" description:"The data of the dispatch result"`
}
