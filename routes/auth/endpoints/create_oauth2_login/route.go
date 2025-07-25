package create_oauth2_login

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Anti-Raid/api/state"
	"github.com/Anti-Raid/api/types"
	"github.com/bwmarrin/discordgo"
	"github.com/redis/rueidis"

	"github.com/anti-raid/eureka/jsonimpl"
	"github.com/anti-raid/eureka/ratelimit"

	docs "github.com/anti-raid/eureka/doclib"
	"github.com/anti-raid/eureka/uapi"

	"github.com/anti-raid/eureka/crypto"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"golang.org/x/exp/slices"
)

var (
	compiledMessages = uapi.CompileValidationErrors(types.AuthorizeRequest{})
)

func Docs() *docs.Doc {
	return &docs.Doc{
		Summary:     "Create Oauth2 Login",
		Description: "Takes in a ``code`` query parameter and returns a temporary session token`. **Cannot be used outside of the site for security reasons**",
		Req:         types.AuthorizeRequest{},
		Resp:        types.CreateUserSessionResponse{},
	}
}

func Route(d uapi.RouteData, r *http.Request) uapi.HttpResponse {
	limit, err := ratelimit.Ratelimit{
		Expiry:      5 * time.Minute,
		MaxRequests: 15,
		Bucket:      "oauth2_login",
	}.Limit(d.Context, r)

	if err != nil {
		state.Logger.Error("Error while ratelimiting", zap.Error(err), zap.String("bucket", "login"))
		return uapi.DefaultResponse(http.StatusInternalServerError)
	}

	if limit.Exceeded {
		return uapi.HttpResponse{
			Json: types.ApiError{
				Message: "You are being ratelimited. Please try again in " + limit.TimeToReset.String(),
			},
			Headers: limit.Headers(),
			Status:  http.StatusTooManyRequests,
		}
	}

	var req types.AuthorizeRequest

	hresp, ok := uapi.MarshalReqWithHeaders(r, &req, limit.Headers())

	if !ok {
		return hresp
	}

	// Validate the payload
	err = state.Validator.Struct(req)

	if err != nil {
		errors := err.(validator.ValidationErrors)
		return uapi.ValidatorErrorResponse(compiledMessages, errors)
	}

	if req.Protocol != "a1" {
		return uapi.HttpResponse{
			Json: types.ApiError{
				Message: "Your client is outdated and is not supported. Please contact the developers of this client.",
			},
			Status:  http.StatusBadRequest,
			Headers: limit.Headers(),
		}
	}

	if !slices.Contains(state.Config.DiscordAuth.AllowedRedirects, req.RedirectURI) {
		return uapi.HttpResponse{
			Json: types.ApiError{
				Message: "Malformed redirect_uri",
			},
			Status:  http.StatusBadRequest,
			Headers: limit.Headers(),
		}
	}

	codeused, _ := state.Rueidis.Do(d.Context, state.Rueidis.B().Exists().Key("codecache:"+req.Code).Build()).ToInt64()

	if codeused == 1 {
		return uapi.HttpResponse{
			Json: types.ApiError{
				Message: "Code has been clearly used before and is as such invalid",
			},
			Status:  http.StatusBadRequest,
			Headers: limit.Headers(),
		}
	}

	err = state.Rueidis.Do(d.Context, state.Rueidis.B().Set().Key("codecache:"+req.Code).Value("0").Ex(5*time.Minute).Build()).Error()

	if err != nil && !errors.Is(err, rueidis.Nil) {
		state.Logger.Error("Failed to set code cache", zap.Error(err))
		return uapi.HttpResponse{
			Json: types.ApiError{
				Message: "Failed to set code cache",
			},
			Status:  http.StatusInternalServerError,
			Headers: limit.Headers(),
		}
	}

	httpResp, err := http.PostForm("https://discord.com/api/v10/oauth2/token", url.Values{
		"client_id":     {state.Config.DiscordAuth.ClientID},
		"client_secret": {state.Config.DiscordAuth.ClientSecret},
		"grant_type":    {"authorization_code"},
		"code":          {req.Code},
		"redirect_uri":  {req.RedirectURI},
		"scope":         {"identify"},
	})

	if err != nil {
		state.Logger.Error("Failed to send oauth2 token request to discord", zap.Error(err))
		return uapi.HttpResponse{
			Json: types.ApiError{
				Message: "Failed to send token request to Discord",
			},
			Status:  http.StatusInternalServerError,
			Headers: limit.Headers(),
		}
	}

	defer httpResp.Body.Close()

	body, err := io.ReadAll(httpResp.Body)

	if err != nil {
		state.Logger.Error("Failed to read oauth2 token response from discord", zap.Error(err))
		return uapi.HttpResponse{
			Json: types.ApiError{
				Message: "Failed to read token response from Discord",
			},
			Status:  http.StatusInternalServerError,
			Headers: limit.Headers(),
		}
	}

	var token struct {
		AccessToken string `json:"access_token"`
		Scope       string `json:"scope"`
	}

	err = jsonimpl.Unmarshal(body, &token)

	if err != nil {
		state.Logger.Error("Failed to parse oauth2 token response from discord", zap.Error(err))
		return uapi.HttpResponse{
			Json: types.ApiError{
				Message: "Failed to parse token response from Discord",
			},
			Status:  http.StatusBadRequest,
			Headers: limit.Headers(),
		}
	}

	// Validate scope
	scopesStr := strings.Split(strings.ReplaceAll(token.Scope, "%20", ""), " ")

	if !slices.Contains(scopesStr, "identify") {
		return uapi.HttpResponse{
			Json: types.ApiError{
				Message: "Required identify scope not granted",
			},
			Status:  http.StatusBadRequest,
			Headers: limit.Headers(),
		}
	}

	if !slices.Contains(scopesStr, "guilds") {
		return uapi.HttpResponse{
			Json: types.ApiError{
				Message: "Required guilds scope not granted",
			},
			Status:  http.StatusBadRequest,
			Headers: limit.Headers(),
		}
	}

	if token.AccessToken == "" {
		state.Logger.Error("No access token provided by discord")
		return uapi.HttpResponse{
			Json: types.ApiError{
				Message: "No access token provided by Discord",
			},
			Status:  http.StatusBadRequest,
			Headers: limit.Headers(),
		}
	}

	cli := &http.Client{}

	var httpReq *http.Request
	httpReq, err = http.NewRequestWithContext(d.Context, "GET", "https://discord.com/api/v10/users/@me", nil)

	if err != nil {
		state.Logger.Error("Failed to create oauth2 request to discord", zap.Error(err))
		return uapi.HttpResponse{
			Json: types.ApiError{
				Message: "Failed to create request to Discord to fetch user info",
			},
			Status:  http.StatusInternalServerError,
			Headers: limit.Headers(),
		}
	}

	httpReq.Header.Set("Authorization", "Bearer "+token.AccessToken)

	httpResp, err = cli.Do(httpReq)

	if err != nil {
		state.Logger.Error("Failed to send oauth2 request to discord", zap.Error(err))
		return uapi.HttpResponse{
			Json: types.ApiError{
				Message: "Failed to send oauth2 request to Discord",
			},
			Status:  http.StatusInternalServerError,
			Headers: limit.Headers(),
		}
	}

	defer httpResp.Body.Close()

	body, err = io.ReadAll(httpResp.Body)

	if err != nil {
		state.Logger.Error("Failed to read oauth2 response from discord", zap.Error(err))
		return uapi.HttpResponse{
			Json: types.ApiError{
				Message: "Failed to read oauth2 response from Discord",
			},
			Status:  http.StatusInternalServerError,
			Headers: limit.Headers(),
		}
	}

	var user discordgo.User

	err = jsonimpl.Unmarshal(body, &user)

	if err != nil {
		state.Logger.Error("Failed to parse oauth2 response from discord", zap.Error(err))
		return uapi.HttpResponse{
			Json: types.ApiError{
				Message: "Failed to parse oauth2 response from Discord",
			},
			Status:  http.StatusInternalServerError,
			Headers: limit.Headers(),
		}
	}

	if user.ID == "" {
		state.Logger.Error("No user ID provided by discord. Invalid code/access token?")
		return uapi.HttpResponse{
			Json: types.ApiError{
				Message: "No user ID provided by Discord. Invalid code/access token?",
			},
			Status:  http.StatusBadRequest,
			Headers: limit.Headers(),
		}
	}

	// Check if user exists on database
	var exists bool

	err = state.Pool.QueryRow(d.Context, "SELECT EXISTS(SELECT 1 FROM users WHERE user_id = $1)", user.ID).Scan(&exists)

	if err != nil {
		state.Logger.Error("Failed to check if user exists on database", zap.Error(err), zap.String("userID", user.ID))
		return uapi.HttpResponse{
			Json: types.ApiError{
				Message: "Failed to check if user exists on database",
			},
			Status:  http.StatusInternalServerError,
			Headers: limit.Headers(),
		}
	}

	if !exists {
		// Create user
		_, err = state.Pool.Exec(
			d.Context,
			"INSERT INTO users (user_id, access_token) VALUES ($1, $2)",
			user.ID,
			token.AccessToken,
		)

		if err != nil {
			state.Logger.Error("Failed to create user on database", zap.Error(err), zap.String("userID", user.ID))
			return uapi.HttpResponse{
				Json: types.ApiError{
					Message: "Failed to create user on database",
				},
				Status:  http.StatusInternalServerError,
				Headers: limit.Headers(),
			}
		}
	} else {
		// Get state
		var s string

		err = state.Pool.QueryRow(d.Context, "SELECT state FROM users WHERE user_id = $1", user.ID).Scan(&s)

		if err != nil {
			state.Logger.Error("Failed to get user state from database", zap.Error(err), zap.String("userID", user.ID))
			return uapi.HttpResponse{
				Json: types.ApiError{
					Message: "Failed to get user state from database",
				},
				Status:  http.StatusInternalServerError,
				Headers: limit.Headers(),
			}
		}

		banned := s == "banned" || s == "api_banned"

		if banned && req.Scope != "ban_exempt" {
			return uapi.HttpResponse{
				Json: types.ApiError{
					Message: "You are banned from Anti-Raid. If you think this is a mistake, please contact support.",
				},
				Status:  http.StatusForbidden,
				Headers: limit.Headers(),
			}
		}

		if !banned && req.Scope == "ban_exempt" {
			return uapi.HttpResponse{
				Json: types.ApiError{
					Message: "The selected scope is not allowed for unbanned users [ban_exempt].",
				},
				Status:  http.StatusForbidden,
				Headers: limit.Headers(),
			}
		}

		// Update access token
		err = state.Pool.QueryRow(d.Context, "UPDATE users SET access_token = $1 WHERE user_id = $2 RETURNING access_token", token.AccessToken, user.ID).Scan(&token.AccessToken)

		if err != nil {
			state.Logger.Error("Failed to update access token on database", zap.Error(err), zap.String("userID", user.ID))
			return uapi.HttpResponse{
				Json: types.ApiError{
					Message: "Failed to update access token on database",
				},
				Status:  http.StatusInternalServerError,
				Headers: limit.Headers(),
			}
		}
	}

	var sessionToken = crypto.RandString(128)
	var sessionId string
	var expiry time.Time

	err = state.Pool.QueryRow(d.Context, "INSERT INTO web_api_tokens (user_id, type, token, expiry) VALUES ($1, 'login', $2, NOW() + INTERVAL '1 hour') RETURNING id, expiry", user.ID, sessionToken).Scan(&sessionId, &expiry)

	if err != nil {
		state.Logger.Error("Failed to create session token", zap.Error(err), zap.String("userID", user.ID))
		return uapi.HttpResponse{
			Json: types.ApiError{
				Message: "Failed to create session token",
			},
			Status:  http.StatusInternalServerError,
			Headers: limit.Headers(),
		}
	}

	// Create authUser and send
	var authUser = types.CreateUserSessionResponse{
		UserID:    user.ID,
		SessionID: sessionId,
		Token:     sessionToken,
		Expiry:    expiry,
		User: &types.PartialUser{
			ID:         user.ID,
			Username:   user.Username,
			GlobalName: user.GlobalName,
			Avatar:     user.Avatar,
		},
	}

	return uapi.HttpResponse{
		Json:    authUser,
		Headers: limit.Headers(),
	}
}
