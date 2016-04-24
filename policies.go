package agree

import (
	"bytes"

	"github.com/pkg/errors"
)

// RequireOracle is an agreement policy that only causes an error if the system
// doesn't have at least one oracle configured.
type RequireOracle struct{}

var _ Policy = RequireOracle{}

// ApplyAgreementPolicy implements `Policy`
func (p RequireOracle) ApplyAgreementPolicy(
	sys *System,
	query *Query,
) error {

	if len(sys.Oracles) == 0 {
		return errors.New("require-oracle: no oracles configured")
	}

	return nil
}

// UnanimousConsent is an agreement policy that only succeeds if every oracle in
// a system agrees on the same value.
type UnanimousConsent struct{}

var _ Policy = UnanimousConsent{}

// ApplyAgreementPolicy implements `Policy`
func (p UnanimousConsent) ApplyAgreementPolicy(
	sys *System,
	query *Query,
) error {

	for _, o := range sys.Oracles {
		value, err := o.GetOracleView(query.Domain, query.Key)
		if err != nil {
			return errors.Wrap(err, "unanimous-consent: oracle error")
		}

		if query.Result == nil {
			query.Result = value
		}

		if !bytes.Equal(query.Result, value) {
			return errors.New("unanimous-consent: oracles disagree")
		}
	}

	return nil
}
