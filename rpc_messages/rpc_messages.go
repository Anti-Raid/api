package rpc_messages

import (
	"github.com/Anti-Raid/api/types"
	"github.com/bwmarrin/discordgo"
)

type BaseGuildUserInfo struct {
	OwnerID   string                              `json:"owner_id"`
	Name      string                              `json:"name"`
	Icon      *string                             `json:"icon"`
	Roles     []types.SerenityRole                `json:"roles"`
	UserRoles []string                            `json:"user_roles"`
	BotRoles  []string                            `json:"bot_roles"`
	Channels  []types.GuildChannelWithPermissions `json:"channels"`
}

type SettingsOperationRequest struct {
	Fields  any `json:"fields"`
	Op      string         `json:"op"`
	Setting string         `json:"setting"`
}

type DispatchResult struct {
	Type string `json:"type" description:"The type of the dispatch result [Ok or Err]"`
	Data any    `json:"data" description:"The data of the dispatch result"`
}

type TWState struct {
	Commands []discordgo.ApplicationCommand `json:"commands"`
}
