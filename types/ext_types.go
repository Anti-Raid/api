package types

import "github.com/bwmarrin/discordgo"

type Permissions string

// Discordgo types are sometimes not of high quality, so we need to extend them
// with our own types [taken from serenity etc]. This is the place to do that.

type SerenityRoleTags struct {
	BotID                 *string `json:"bot_id" description:"The ID of the bot the role belongs to"`
	IntegrationID         *string `json:"integration_id" description:"The ID of the integration the role belongs to"`
	PremiumSubscriber     bool    `json:"premium_subscriber" description:"Whether this is the guild's premium subscriber role"`
	SubscriptionListingID *string `json:"subscription_listing_id" description:"The id of this role's subscription sku and listing"`
	AvailableForPurchase  bool    `json:"available_for_purchase" description:"Whether this role is available for purchase"`
	GuildConnections      bool    `json:"guild_connections" description:"Whether this role is a guild's linked role"`
}

type SerenityRole struct {
	ID           string            `json:"id" description:"The ID of the role"`
	GuildID      string            `json:"guild_id" description:"The ID of the guild"`
	Color        int               `json:"color" description:"The color of the role"`
	Name         string            `json:"name" description:"The name of the role"`
	Permissions  *Permissions      `json:"permissions" description:"The permissions of the role"`
	Position     int16             `json:"position" description:"The position of the role"`
	Tags         *SerenityRoleTags `json:"tags" description:"The tags of the role"`
	Icon         *string           `json:"icon" description:"The icon of the role"`
	UnicodeEmoji string            `json:"unicode_emoji" description:"The unicode emoji of the role"`
}

type GuildChannelWithPermissions struct {
	User    Permissions        `json:"user" description:"The permissions the user has in the channel"`
	Bot     Permissions        `json:"bot" description:"The permissions the bot has in the channel"`
	Channel *discordgo.Channel `json:"channel" description:"The channel object"`
}
