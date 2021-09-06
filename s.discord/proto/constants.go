package discordproto

const (
	DiscordRolesUpdateActorPaymentsSystem = "actor-payments-system"
)

const (
	// Guilds
	DiscordSatoshiGuildID = "814144801977008228"
)

const (
	// Satoshi Discord Channel IDs
	DiscordSatoshiTwitterChannel   = "816794087868465163"
	DiscordSatoshiAlertsChannel    = "816794120851816479"
	DiscordSatoshiTestingChannel   = "817513133274824715"
	DiscordSatoshiWhaleChannel     = "817789196319195166"
	DiscordSatoshiTradersChannel   = "817789261415448606"
	DiscordSatoshiNewsChannel      = "817789219656826970"
	DiscordSatoshiExchangesChannel = "818909423530541116"
	DiscordSatoshiProjectsChannel  = "826528849374216192"
	DiscordSatoshiPriceBotChannel  = "831234720943702066"
	DiscordSatoshiGeneralChannel   = "814144802458828852"

	// Futures channels
	DiscordSatoshiModMessagesChannel = "847954019758112808"
	DiscordSatoshiSwingsChannel      = "847953925575671848"
	DiscordSatoshiFuturesChannel     = "814146537088221284"
	DiscordSatoshiModTradesChannel   = "883692707611050024"
	DiscordSatoshiChallengesChannel  = "884524239410061323"

	// Satoshi Discord Pulse Channels
	DiscordSatoshiAccountsPulseChannel = "883709489101033602"
	DiscordSatoshiPaymentsPulseChannel = "883306360798859274"
)

const (
	// Role IDs
	DiscordSatoshiFuturesRoleID = "828590713440043019"
	DiscordSatoshiAdminRoleID   = "816722849599455234"

	// Role Names
	DiscordSatoshiFuturesRole = "satoshi-futures-role"
	DiscordSatoshiAdminRole   = "satoshi-admin-role"
)

var (
	discordSatoshiRoleIDToName = map[string]string{
		DiscordSatoshiFuturesRoleID: DiscordSatoshiFuturesRole,
		DiscordSatoshiAdminRoleID:   DiscordSatoshiAdminRole,
	}

	discordSatoshiRoleNameToID = map[string]string{
		DiscordSatoshiFuturesRole: DiscordSatoshiFuturesRoleID,
		DiscordSatoshiAdminRole:   DiscordSatoshiAdminRoleID,
	}
)

// ConvertRoleIDToName
//
// Converts the discord role id to a name.
func ConvertRoleIDToName(roleID string) (string, bool) {
	v, ok := discordSatoshiRoleIDToName[roleID]
	return v, ok
}

// ConvertRoleNameToID
//
// Converts the discord role name to an id.
func ConvertRoleNameToID(name string) (string, bool) {
	v, ok := discordSatoshiRoleNameToID[name]
	return v, ok
}

const (
	// Other Discord Channel IDs
	DiscordMoonModMessagesChannel = "813362955516903484"
	DiscordMoonSwingGroupChannel  = "814141004508561419"
)
