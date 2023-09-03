package service

import "gin-jwt-boilerplate/internal/domain"

// GetApplicationHealth returns the health of the application
// By default, it returns UP
func GetApplicationHealth() domain.Health {
	return domain.Health{Status: "UP"}
}
