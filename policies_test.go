package agree

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPolicy_RequireOracle(t *testing.T) {
	d := "some-domain"
	k := "some-key"
	sys := &System{
		Policies: []Policy{
			RequireOracle{},
		},
	}

	// no oracle
	_, err := sys.Get(d, k)
	assert.EqualError(t, err, "require-oracle: no oracles configured")

	// one oracle
	sys.Oracles = append(sys.Oracles, MockOracle{})
	require.Len(t, sys.Oracles, 1, "bad oracle count")

	_, err = sys.Get(d, k)
	assert.NoError(t, err)
}

func TestPolicy_UnanimousConsent(t *testing.T) {
	d := "some-domain"
	k := "some-key"
	sys := &System{
		Oracles: []Oracle{
			MockOracle{"some-key": "some-value"},
		},
		Policies: []Policy{
			UnanimousConsent{},
		},
	}

	// single oracle
	v, err := sys.Get(d, k)
	if assert.NoError(t, err) {
		assert.Equal(t, []byte("some-value"), v)
	}

	// multiple agreeing oracles
	sys.Oracles = append(sys.Oracles, sys.Oracles[0])
	require.Len(t, sys.Oracles, 2, "bad oracle count")

	v, err = sys.Get(d, k)
	if assert.NoError(t, err) {
		assert.Equal(t, []byte("some-value"), v)
	}

	//  disagreement
	sys.Oracles = append(sys.Oracles, MockOracle{"some-key": "some-other-value"})
	require.Len(t, sys.Oracles, 3, "bad oracle count")

	_, err = sys.Get(d, k)
	assert.EqualError(t, err, "unanimous-consent: oracles disagree")

	//  empty oracles
	sys.Oracles = nil
	require.Len(t, sys.Oracles, 0, "bad oracle count")

	v, err = sys.Get(d, k)
	assert.NoError(t, err)
	assert.Nil(t, v)
}
