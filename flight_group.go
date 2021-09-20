package icache

// FlightGroupIf flight group
type FlightGroupIf interface {
	Do(string, func() (interface{}, error)) (interface{}, error)
}
