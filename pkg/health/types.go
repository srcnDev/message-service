package health

// Status represents health status
type Status struct {
	Status string `json:"status"`
	Uptime string `json:"uptime"`
}
