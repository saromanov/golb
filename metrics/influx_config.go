package metrics

const (
	defaultInfluxAddress = "http://localhost:8086"
	defaultInfluxDB      = "golbinflux"
)

// InfluxConfig defines configuration
// for init influx
type InfluxConfig struct {
	DBName  string
	Address string
}

// GetDBName returns db name
func (c *InfluxConfig) GetDBName() string {
	if c.DBName == "" {
		return c.makeDefaultDBName()
	}
	return c.DBName
}

// GetAddress returns address
func (c *InfluxConfig) GetAddress() string {
	if c.Address == "" {
		return c.makeDefaultAddress()
	}
	return c.Address
}

func (c *InfluxConfig) makeDefaultDBName() string {
	return defaultInfluxDB
}

func (c *InfluxConfig) makeDefaultAddress() string {
	return defaultInfluxAddress
}
