package status

// Conf is a configuration for each Observable.
type Conf struct {
	// Session is an Observable for a SessionID status.
	Session *Observable
}

// Observable represents a Status service configuration.
type Observable struct {
	// ServerName is used to resolve server SNI during a TLS handshake.
	ServerName string

	// AuthTTL is a time in seconds for user to authenticate.
	AuthTTL int64

	// ChoiceTTL is a time in seconds for user to make a choice.
	ChoiceTTL int64

	// VoteTTL is a time in seconds for user to confirm a choice.
	VoteTTL int64

	// VerifyTTL is a time in seconds for user to verify a choice.
	VerifyTTL int64
}
