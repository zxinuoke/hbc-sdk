package utils

// Msg - Transactions messages must fulfill the Msg
type Msg interface {
	// Return the message type.
	// Must be alphanumeric or empty.
	Route() string

	// Returns a human-readable string for the message, intended for utilization
	// within tags
	Type() string

	// ValidateBasic does a simple validation check that
	// doesn't require access to any other information.
	ValidateBasic() error

	// Get the canonical byte representation of the Msg.
	GetSignBytes() []byte

	// Signers returns the addrs of signers that must sign.
	// CONTRACT: All signatures must be present to be valid.
	// CONTRACT: Returns addrs in some deterministic order.
	GetSigners() []CUAddress

	// Get involved addresses of this msg so that we can publish account balance change
	GetInvolvedAddresses() []CUAddress
}
