package simulation

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/commitHub/commitBlockchain/baseapp"
	"github.com/commitHub/commitBlockchain/x/mock/simulation"
)

// AllInvariants tests all governance invariants
func AllInvariants() simulation.Invariant {
	return func(t *testing.T, app *baseapp.BaseApp, log string) {
		// TODO Add some invariants!
		// Checking proposal queues, no passed-but-unexecuted proposals, etc.
		require.Nil(t, nil)
	}
}
