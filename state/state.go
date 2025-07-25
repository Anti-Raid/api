package state

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Anti-Raid/api/config"
	"github.com/Anti-Raid/api/state/redishotcache"
	"golang.org/x/net/http2"

	"github.com/anti-raid/eureka/dovewing"
	"github.com/anti-raid/eureka/dovewing/dovetypes"
	"github.com/anti-raid/eureka/proxy"
	"github.com/anti-raid/eureka/ratelimit"
	"github.com/anti-raid/eureka/snippets"
	"github.com/bwmarrin/discordgo"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/rueidis"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

var (
	Pool                    *pgxpool.Pool
	Rueidis                 rueidis.Client // where perf is needed
	DovewingPlatformDiscord *DiscordState
	Discord                 *discordgo.Session
	Logger                  *zap.Logger
	Context                 = context.Background()
	Validator               = validator.New()
	BotUser                 *discordgo.User
	CurrentOperationMode    string // Current mode splashtail is operating in
	Config                  *config.Config

	IpcClient       http.Client
	IpcClientHttp11 http.Client
)

func Setup() {
	cfg, err := os.ReadFile("config.yaml")

	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(cfg, &Config)

	if err != nil {
		panic(err)
	}

	err = Validator.Struct(Config)

	if err != nil {
		panic("configError: " + err.Error())
	}

	Logger = snippets.CreateZap()

	// Postgres
	Pool, err = pgxpool.New(Context, Config.Meta.PostgresURL)

	if err != nil {
		panic(err)
	}

	// Reuidis
	ruOptions, err := rueidis.ParseURL(Config.Meta.RedisURL)

	if err != nil {
		panic(err)
	}

	Rueidis, err = rueidis.NewClient(ruOptions)

	if err != nil {
		panic(err)
	}

	// Discordgo
	Discord, err = discordgo.New("Bot " + Config.DiscordAuth.Token)

	if err != nil {
		panic(err)
	}

	Discord.Client.Transport = proxy.NewHostRewriter(strings.Replace(Config.Meta.Proxy, "http://", "", 1), http.DefaultTransport, func(s string) {
		Logger.Info("[PROXY]", zap.String("note", s))
	})

	// Verify token
	bu, err := Discord.User("@me")

	if err != nil {
		panic(err)
	}

	BotUser = bu

	// Load dovewing state
	baseDovewingState := dovewing.BaseState{
		Pool:    Pool,
		Logger:  Logger,
		Context: Context,
		PlatformUserCache: redishotcache.RuedisHotCache[dovetypes.PlatformUser]{
			Redis:  Rueidis,
			Prefix: "uobj__",
			For:    "dovewing",
		},
		UserExpiryTime: 8 * time.Hour,
	}

	DovewingPlatformDiscord, err = DiscordStateConfig{
		Session:        Discord,
		PreferredGuild: Config.Servers.Main,
		BaseState:      &baseDovewingState,
	}.New()

	if err != nil {
		panic(err)
	}

	ratelimit.SetupState(&ratelimit.RLState{
		HotCache: redishotcache.RuedisHotCache[int]{
			Redis:    Rueidis,
			Prefix:   "rl:",
			For:      "ratelimit",
			Disabled: Config.Meta.WebDisableRatelimits,
		},
	})

	IpcClient.Timeout = 30 * time.Second
	IpcClient.Transport = &http2.Transport{
		AllowHTTP:      true,
		DialTLSContext: nil,
		DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
			return net.Dial(network, addr)
		},
	}

	IpcClientHttp11.Timeout = 30 * time.Second
	IpcClientHttp11.Transport = &http.Transport{}
}
