package domain

// Health is a struct that contains health information about the application
// It is used to return the health information
type Health struct {
	Status string `json:"status"`
}
