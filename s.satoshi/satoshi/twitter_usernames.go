package satoshi

import (
	discordproto "swallowtail/s.discord/proto"
	"sync"
)

var (
	twitterUsernameMu       sync.Mutex
	usernameMetadataMapping = map[string]*TwitterUserMetaData{
		/// --- Traders --- ///

		"AltcoinPsycho": {
			Id:             "942999039192186882",
			DiscordChannel: discordproto.DiscordSatoshiTradersChannel,
		},
		"CiNoTpyrC": {
			Name:           "Rego",
			Bio:            "Supply/Demand trader; excellent with BTC calls.",
			Id:             "1134440744939397122", // rego
			DiscordChannel: discordproto.DiscordSatoshiTradersChannel,
		},
		"zhusu": {
			Name:           "Zhu Su",
			Bio:            "Trader/Invester at Three Arrows Capital, one of the best crypto investment firms in the world.",
			Id:             "79714172",
			DiscordChannel: discordproto.DiscordSatoshiTradersChannel,
		},
		"kyled116": {
			Name:           "Kyle Davies",
			Bio:            "Trader/Invester at Three Arrows Capital, one of the best crypto investment firms in the world.",
			Id:             "1140429573978378241",
			DiscordChannel: discordproto.DiscordSatoshiTradersChannel,
		},
		"e101y7": {
			Id:             "81550741",
			Name:           "Kieran",
			Bio:            "Trend follow; part of WWG",
			DiscordChannel: discordproto.DiscordSatoshiTradersChannel,
		},
		"jclcapital": {
			Id:             "931642527110742016",
			Name:           "Jordan Lindsey",
			Bio:            "Macro trader, following BTC. Nailed it the past year or so.",
			DiscordChannel: discordproto.DiscordSatoshiTradersChannel,
			Youtube:        "https://www.youtube.com/channel/UCN2WmKUchJpIcS1MupY-BuA",
			Twitter:        "twitter.com/jclcapital",
		},
		"MrktMeditations": {
			Id:             "1293878938540933122",
			Name:           "Market Mediatations",
			Bio:            "A mathematician who has made millions; this is his newschannel",
			DiscordChannel: discordproto.DiscordSatoshiNewsChannel,
			Youtube:        "https://www.youtube.com/channel/UCQEBsgNV0RGm1O2iOMDbZjA",
		},
		"ColdBloodedShill": {
			Id:             "987343085251317760",
			Name:           "ColdBloodedShiller",
			Bio:            "Monster trader & streamer.",
			DiscordChannel: discordproto.DiscordSatoshiTradersChannel,
			Twitter:        "https://twitter.com/ColdBloodShill",
		},
		"mickyMafiaTrade": {
			Id:             "1007325467001618434",
			Name:           "Michele",
			DiscordChannel: discordproto.DiscordSatoshiTradersChannel,
			Bio:            "The Oracle; serious this guy makes some unreal calls.",
		},
		"SmartContracter": {
			Id:             "939058273487003648",
			Name:           "Smart contracter",
			DiscordChannel: discordproto.DiscordSatoshiTradersChannel,
			Bio:            "Cryptocurrency swing trader.",
			Twitter:        "https://twitter.com/SmartContracter/status/1369734816539668483",
		},
		"CryptoKaleo": {
			Id:             "906234475604037637",
			Name:           "Crypto Kaleo",
			DiscordChannel: discordproto.DiscordSatoshiTradersChannel,
			Bio:            "Unbelievable Alt coin trader; ignore at your own peril.",
			Twitter:        "https://twitter.com/CryptoKaleo",
		},
		"ShardiB2": {
			Id:             "1091099554605395968",
			Name:           "ShardiB2",
			DiscordChannel: discordproto.DiscordSatoshiTradersChannel,
			Bio:            "Similar to Kaleo, great alt coin trader, she also shares solid algos",
			Twitter:        "https://twitter.com/ShardiB2",
		},
		"Pentosh1": {
			Id:             "1138993163706753029",
			Name:           "Pentosh1",
			DiscordChannel: discordproto.DiscordSatoshiTradersChannel,
			Twitter:        "https://twitter.com/Pentosh1",
		},

		/// --- Twitter Channels --- ///

		"michael_saylor": {
			Bio:            "Single-handedly is trying to take BTC to 100k, CEO microstrategy.",
			Id:             "244647486",
			DiscordChannel: discordproto.DiscordSatoshiTwitterChannel,
		},
		"elonmusk": {
			Bio:            "Full time shitposter extraordinaire, part time pump and dump fiend.",
			Id:             "44196397",
			DiscordChannel: discordproto.DiscordSatoshiTwitterChannel,
		},

		/// --- Exchanges --- ///

		"CoinbasePro": {
			Id:             "720487892670410753",
			Name:           "Coinbase Pro",
			DiscordChannel: discordproto.DiscordSatoshiExchangesChannel,
			Bio:            "Crypto Exchange based in the US.",
			Twitter:        "CoinbasePro",
		},

		/// --- Projects --- ///

		"Syntropynet": {
			Id:             "959803528222003200",
			Name:           "Syntropy (NOIA)",
			DiscordChannel: discordproto.DiscordSatoshiProjectsChannel,
			Bio:            "Syntropy Project",
			Twitter:        "https://twitter.com/Syntropynet",
		},
		"chain_swap": {
			Id:             "1367247307117391872",
			Name:           "Chain Swap",
			DiscordChannel: discordproto.DiscordSatoshiProjectsChannel,
			Bio:            "Chain Swap (airdrop)",
			Twitter:        "https://twitter.com/chain_swap",
			Tags:           []string{"airdrop"},
		},
		"BosonProtocol": {
			Id:             "1193497953933152256",
			Name:           "Boson Protocol",
			DiscordChannel: discordproto.DiscordSatoshiProjectsChannel,
			Bio:            "dCommerce & NFTs",
			Twitter:        "https://twitter.com/BosonProtocol",
		},

		// -- Solana -- //

		"mangomarkets": {
			Id:             "1344200569016246272",
			Name:           "Mango Markets",
			DiscordChannel: discordproto.DiscordSatoshiProjectsChannel,
			Bio:            "DEX with cross margin leverage trading on Solana",
			Twitter:        "https://twitter.com/mangomarkets",
		},

		"Sola_System": {
			Id:             "1376809462673993728",
			Name:           "Sola System",
			DiscordChannel: discordproto.DiscordSatoshiProjectsChannel,
			Bio:            "Solana Project TBD.",
			Twitter:        "https://twitter.com/Sola_System",
		},
		"solstarterorg": {
			Id:             "1374187417003945990",
			Name:           "Solstarter",
			DiscordChannel: discordproto.DiscordSatoshiProjectsChannel,
			Bio:            "IDO / DeFI on Solana",
			Twitter:        "https://twitter.com/solstarterorg",
		},
		"RampDefi": {
			Id:             "1219119953728557056",
			Name:           "Ramp",
			DiscordChannel: discordproto.DiscordSatoshiProjectsChannel,
			Bio:            "Liquidity on ramp solution on Solana",
			Twitter:        "https://twitter.com/RampDefi",
		},

		/// --- On-Chain Analytics --- ///

		"Whale_Sniper": {
			Id:             "1122965827702149120",
			Name:           "Whale Sniper",
			DiscordChannel: discordproto.DiscordSatoshiWhaleChannel,
			Bio:            "On chain analysis for unusual activity for given assets.",
			Twitter:        "https://twitter.com/Whale_Sniper",
		},

		/// --- News --- ///

		"Fxhedgers": {
			Id:             "31064165",
			Bio:            "Market News Channel",
			DiscordChannel: discordproto.DiscordSatoshiNewsChannel,
		},
		"HalvingTracker": {
			Id:             "1350187523759251456",
			Bio:            "An observation the current 4 year market cycle relative to the two prior.",
			DiscordChannel: discordproto.DiscordSatoshiNewsChannel,
		},
	}
)

func getMetadataMapping(username string) (*TwitterUserMetaData, bool) {
	twitterUsernameMu.Lock()
	defer twitterUsernameMu.Unlock()
	user, ok := usernameMetadataMapping[username]
	return user, ok
}
