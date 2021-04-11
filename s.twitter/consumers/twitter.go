package consumers

import (
	"context"
	"swallowtail/s.twitter/clients"
	"sync"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/monzo/slog"
)

var (
	twitterMtx sync.Mutex
	consumers  = []Consumer{}

	rolesMtx sync.RWMutex
	rolesID  = map[string]string{
		"futures": "<@&828590713440043019>",
	}

	usernameMetadataMapping = map[string]*TwitterUserMetaData{
		/// --- Traders --- ///

		"AltcoinPsycho": {
			Id:             "942999039192186882",
			DiscordChannel: discordTradersChannel,
		},
		"CiNoTpyrC": {
			Name:           "Rego",
			Bio:            "Supply/Demand trader; excellent with BTC calls.",
			Id:             "1134440744939397122", // rego
			DiscordChannel: discordTradersChannel,
		},
		"zhusu": {
			Name:           "Zhu Su",
			Bio:            "Trader/Invester at Three Arrows Capital, one of the best crypto investment firms in the world.",
			Id:             "79714172",
			DiscordChannel: discordTradersChannel,
		},
		"kyled116": {
			Name:           "Kyle Davies",
			Bio:            "Trader/Invester at Three Arrows Capital, one of the best crypto investment firms in the world.",
			Id:             "1140429573978378241",
			DiscordChannel: discordTradersChannel,
		},
		"Fxhedgers": {
			Id:             "31064165",
			Bio:            "Market News Channel",
			DiscordChannel: discordNewsChannel,
		},
		"e101y7": {
			Id:             "81550741",
			Name:           "Kieran",
			Bio:            "Trend follow; part of WWG",
			DiscordChannel: discordTradersChannel,
		},
		"jclcapital": {
			Id:             "931642527110742016",
			Name:           "Jordan Lindsey",
			Bio:            "Macro trader, following BTC. Nailed it the past year or so.",
			DiscordChannel: discordTradersChannel,
			Youtube:        "https://www.youtube.com/channel/UCN2WmKUchJpIcS1MupY-BuA",
			Twitter:        "twitter.com/jclcapital",
		},
		"MrktMeditations": {
			Id:             "1293878938540933122",
			Name:           "Market Mediatations",
			Bio:            "A mathematician who has made millions; this is his news channel",
			DiscordChannel: discordNewsChannel,
			Youtube:        "https://www.youtube.com/channel/UCQEBsgNV0RGm1O2iOMDbZjA",
		},
		"ChocolateMastr": {
			Id:             "947839753717731328",
			Name:           "Willy Wonka",
			Bio:            "Great calls, extensive research for spot buys.",
			DiscordChannel: discordTradersChannel,
		},
		"TraderKoz": {
			Id:             "1019660554472837120",
			Name:           "TraderKoz",
			Bio:            "Respectable trader in the twittersphere.",
			DiscordChannel: discordTradersChannel,
			Twitter:        "https://twitter.com/TraderKoz",
		},
		"ColdBloodedShill": {
			Id:             "987343085251317760",
			Name:           "ColdBloodedShiller",
			Bio:            "Monster trader & streamer.",
			DiscordChannel: discordTradersChannel,
			Twitter:        "https://twitter.com/ColdBloodShill",
		},
		"mickyMafiaTrade": {
			Id:             "1007325467001618434",
			Name:           "Michele",
			DiscordChannel: discordTradersChannel,
			Bio:            "The Oracle; serious this guy makes some unreal calls.",
		},
		"eliz883": {
			Id:             "993962483332329472",
			Name:           "Eli",
			DiscordChannel: discordTradersChannel,
			Bio:            "Part of WWG; catches bottoms",
			Twitter:        "https://twitter.com/eliz883",
		},
		"SmartContracter": {
			Id:             "939058273487003648",
			Name:           "Smart contracter",
			DiscordChannel: discordTradersChannel,
			Bio:            "Cryptocurrency swing trader.",
			Twitter:        "https://twitter.com/SmartContracter/status/1369734816539668483",
		},
		"CryptoKaleo": {
			Id:             "906234475604037637",
			Name:           "Crypto Kaleo",
			DiscordChannel: discordTradersChannel,
			Bio:            "Unbelievable Alt coin trader; ignore at your own peril.",
			Twitter:        "https://twitter.com/CryptoKaleo",
			Tags:           []string{"futures"},
		},
		"ShardiB2": {
			Id:             "1091099554605395968",
			Name:           "ShardiB2",
			DiscordChannel: discordTradersChannel,
			Bio:            "Similar to Kaleo, great alt coin trader, she also shares solid algos",
			Twitter:        "https://twitter.com/ShardiB2",
			Tags:           []string{"futures"},
		},

		/// --- Twitter Channels --- ///

		"michael_saylor": {
			Bio:            "Single-handedly is trying to take BTC to 100k, CEO microstrategy.",
			Id:             "244647486",
			DiscordChannel: discordTwitterChannel,
		},
		"elonmusk": {
			Bio:            "Full time shitposter extraordinaire, part time pump and dump fiend.",
			Id:             "44196397",
			DiscordChannel: discordTwitterChannel,
		},

		/// --- Exchanges --- ///

		"CoinbasePro": {
			Id:             "720487892670410753",
			Name:           "Coinbase Pro",
			DiscordChannel: discordExchangesChannel,
			Bio:            "Crypto Exchange based in the US.",
			Twitter:        "CoinbasePro",
		},

		/// --- Projects --- ///

		"Syntropynet": {
			Id:             "959803528222003200",
			Name:           "Syntropy (NOIA)",
			DiscordChannel: discordProjectsChannel,
			Bio:            "Syntropy Project",
			Twitter:        "https://twitter.com/Syntropynet",
		},
		"chain_swap": {
			Id:             "1367247307117391872",
			Name:           "Chain Swap",
			DiscordChannel: discordProjectsChannel,
			Bio:            "Chain Swap (airdrop)",
			Twitter:        "https://twitter.com/chain_swap",
			Tags:           []string{"airdrop"},
		},
		"BosonProtocol": {
			Id:             "1193497953933152256",
			Name:           "Boson Protocol",
			DiscordChannel: discordProjectsChannel,
			Bio:            "dCommerce & NFTs",
			Twitter:        "https://twitter.com/BosonProtocol",
		},

		// -- Solana -- //

		"mangomarkets": {
			Id:             "1344200569016246272",
			Name:           "Mango Markets",
			DiscordChannel: discordProjectsChannel,
			Bio:            "DEX with cross margin leverage trading on Solana",
			Twitter:        "https://twitter.com/mangomarkets",
		},

		"Sola_System": {
			Id:             "1376809462673993728",
			Name:           "Sola System",
			DiscordChannel: discordProjectsChannel,
			Bio:            "Solana Project TBD.",
			Twitter:        "https://twitter.com/Sola_System",
		},
		"solstarterorg": {
			Id:             "1374187417003945990",
			Name:           "Solstarter",
			DiscordChannel: discordProjectsChannel,
			Bio:            "IDO / DeFI on Solana",
			Twitter:        "https://twitter.com/solstarterorg",
		},
		"RampDefi": {
			Id:             "1219119953728557056",
			Name:           "Ramp",
			DiscordChannel: discordProjectsChannel,
			Bio:            "Liquidity on ramp solution on Solana",
			Twitter:        "https://twitter.com/RampDefi",
		},

		/// --- On-Chain Analytics --- ///

		"Whale_Sniper": {
			Id:             "1122965827702149120",
			Name:           "Whale Sniper",
			DiscordChannel: discordWhaleChannel,
			Bio:            "On chain analysis for unusual activity for given assets.",
			Twitter:        "https://twitter.com/Whale_Sniper",
		},
		"HalvingTracker": {
			Id:             "1350187523759251456",
			Bio:            "An observation the current 4 year market cycle relative to the two prior.",
			DiscordChannel: discordNewsChannel,
		},
		"whale_alert": {
			Id:             "1039833297751302144",
			Bio:            "Big money movement alerts",
			DiscordChannel: discordWhaleChannel,
		},
	}
)

type TwitterUserMetaData struct {
	Bio            string
	Name           string
	Id             string
	DiscordChannel string
	Emoji          string
	Twitter        string
	Twitch         string
	Youtube        string
	Tags           []string
	Filter         func(string) bool
}

type Consumer func(tweet *twitter.Tweet)

func register(consumer Consumer) {
	twitterMtx.Lock()
	defer twitterMtx.Unlock()
	consumers = append(consumers, consumer)
}

func New() *TwitterConsumer {
	return &TwitterConsumer{
		cli:  clients.New(),
		done: make(chan struct{}, 1),
	}
}

type TwitterConsumer struct {
	cli  *clients.TwitterClient
	done chan struct{}
}

func (tw *TwitterConsumer) Run(ctx context.Context) (func(), error) {
	twitterMtx.Lock()
	defer twitterMtx.Unlock()

	usersToConsumer := []string{}
	for _, user := range usernameMetadataMapping {
		usersToConsumer = append(usersToConsumer, user.Id)
	}

	demux := twitter.NewSwitchDemux()
	demux.Tweet = consumers[0]

	filterParams := &twitter.StreamFilterParams{
		Follow:        usersToConsumer,
		StallWarnings: twitter.Bool(true),
	}

	stream, err := tw.cli.Client.Streams.Filter(filterParams)
	if err != nil {
		return nil, err
	}

	slog.Info(ctx, "Stream created; waiting for tweets...")
	go demux.HandleChan(stream.Messages)

	return func() {
		defer slog.Info(ctx, "Closing twitter stream")
		defer stream.Stop()
	}, nil
}

func (tw *TwitterConsumer) Done(ctx context.Context) {
	go func() {
		defer slog.Info(ctx, "Cancelling twitter consumer")
		tw.done <- struct{}{}
	}()
}

func GetMetadataMapping(username string) (*TwitterUserMetaData, bool) {
	twitterMtx.Lock()
	defer twitterMtx.Unlock()
	user, ok := usernameMetadataMapping[username]
	return user, ok
}

func GetRoleID(roleName string) (string, bool) {
	rolesMtx.RLock()
	defer rolesMtx.RUnlock()
	roleID, ok := rolesID[roleName]
	return roleID, ok
}
