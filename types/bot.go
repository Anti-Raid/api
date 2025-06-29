package types

import (
	"github.com/bwmarrin/discordgo"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

type BotState struct {
	Commands           []discordgo.ApplicationCommand `json:"commands"`
	Settings           []CanonicalConfigOption        `json:"settings"`
	CommandPermissions map[string][]string            `json:"command_permissions"`
}

type CanonicalCommandArgument struct {
	Name        string   `json:"name"`
	Description *string  `json:"description"`
	Required    bool     `json:"required"`
	Choices     []string `json:"choices"`
}

type CanonicalCommand struct {
	Name               string                     `json:"name"`
	QualifiedName      string                     `json:"qualified_name"`
	Description        *string                    `json:"description"`
	NSFW               bool                       `json:"nsfw"`
	Subcommands        []CanonicalCommand         `json:"subcommands"`
	SubcommandRequired bool                       `json:"subcommand_required"`
	Arguments          []CanonicalCommandArgument `json:"arguments"`
}

type CanonicalColumnType struct {
	Scalar *struct {
		Inner CanonicalInnerColumnType `json:"inner"`
	} `json:"Scalar,omitempty"`
	Array *struct {
		Inner CanonicalInnerColumnType `json:"inner"`
	} `json:"Array,omitempty"`
}

type CanonicalInnerColumnType struct {
	Uuid   *struct{} `json:"Uuid,omitempty"`
	String *struct {
		MinLength     *int     `json:"min_length,omitempty"`
		MaxLength     *int     `json:"max_length,omitempty"`
		AllowedValues []string `json:"allowed_values,omitempty"`
		Kind          string   `json:"kind,omitempty"`
	} `json:"String,omitempty"`
	Timestamp   *struct{} `json:"Timestamp,omitempty"`
	TimestampTz *struct{} `json:"TimestampTz,omitempty"`
	Interval    *struct{} `json:"Interval,omitempty"`
	Integer     *struct{} `json:"Integer,omitempty"`
	Float       *struct{} `json:"Float,omitempty"`
	BitFlag     *struct {
		Values orderedmap.OrderedMap[string, int64] `json:"values"`
	} `json:"BitFlag,omitempty"`
	Boolean *struct{} `json:"Boolean,omitempty"`
	Json    *struct {
		MaxBytes *int `json:"max_bytes"`
	} `json:"Json,omitempty"`
}

type CanonicalColumnSuggestion struct {
	Static *struct {
		Suggestions []string `json:"suggestions"`
	} `json:"Static,omitempty"`
	None *struct{} `json:",omitempty"`
}

type CanonicalColumn struct {
	ID          string                    `json:"id"`
	Name        string                    `json:"name"`
	Description string                    `json:"description"`
	ColumnType  CanonicalColumnType       `json:"column_type"`
	Nullable    bool                      `json:"nullable"`
	Suggestions CanonicalColumnSuggestion `json:"suggestions"`
	Secret      bool                      `json:"secret"`
	IgnoredFor  []CanonicalOperationType  `json:"ignored_for"`
}

type CanonicalOperationType string

const (
	View   CanonicalOperationType = "View"
	Create CanonicalOperationType = "Create"
	Update CanonicalOperationType = "Update"
	Delete CanonicalOperationType = "Delete"
)

func (c CanonicalOperationType) List() []string {
	return []string{
		"View",
		"Create",
		"Update",
		"Delete",
	}
}

func (c CanonicalOperationType) Parse() bool {
	for _, v := range c.List() {
		if v == string(c) {
			return true
		}
	}
	return false
}

type CanonicalConfigOption struct {
	ID            string                   `json:"id"`
	Name          string                   `json:"name"`
	Description   string                   `json:"description"`
	PrimaryKey    string                   `json:"primary_key"`
	TitleTemplate string                   `json:"title_template"`
	Columns       []CanonicalColumn        `json:"columns"`
	Operations    []CanonicalOperationType `json:"operations"`
}
