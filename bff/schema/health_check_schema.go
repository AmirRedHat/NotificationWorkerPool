package schema

type HealthCheck struct {
	Status string `json:"status"`
}

func NewHealthCheck() *HealthCheck {
	return &HealthCheck{Status: "OK"}
}
