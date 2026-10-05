package schema

type SendNotificationRequest struct {
	Kind     string `json:"kind"`
	Message  string `json:"message"`
	Receiver string `json:"receiver"`
}

type SendNotificationResponse struct {
	Success bool   `json:"success"`
	Detail  string `json:"detail"`
}
