package agree

// Policy represents some abstract policy which can influence the outcome of an
// agreement calculation.
type Policy interface {
	ApplyAgreementPolicy(sys *System, query *Query) error
}

// Oracle represents some external oracle, which can provide the view of some
// key at any point in time.
type Oracle interface {
	GetOracleView(domain string, key string) ([]byte, error)
}

// Query represents an in-progress query against an agreement system.
type Query struct {
	Domain string
	Key    string
	Result []byte
	Err    error
}

// System represents an agreement system, comprised of a set of oracles and a
// set of agreement policies.
type System struct {
	Oracles  []Oracle
	Policies []Policy
}

// Get resolves the value of `key` using `sys`.  Any policies configured on
// `sys` will be used to augment the agreement.
func (sys *System) Get(domain, key string) ([]byte, error) {
	q := &Query{Domain: domain, Key: key}

	for _, p := range sys.Policies {
		err := p.ApplyAgreementPolicy(sys, q)
		if err != nil {
			return nil, err
		}
	}

	return q.Result, q.Err
}
