package mrfusion

// General errors.
const (
	ErrUpstreamTimeout = Error("request to backend timed out")
	ErrUninitialized   = Error("client uninitialized. Call Open() method")
)

// Error is a domain error encountered while processing mrfusion requests
type Error string

func (e Error) Error() string {
	return string(e)
}
