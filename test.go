package agree

// MockOracle is a test helper
type MockOracle map[string]string

var _ Oracle = MockOracle{}

// GetOracleView implement `Oracle`
func (o MockOracle) GetOracleView(key string) ([]byte, error) {
	if o == nil {
		return nil, nil
	}

	return []byte(o[key]), nil
}
