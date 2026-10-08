package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"task01/bff/config"
	"task01/bff/schema"
	"task01/bff/services"
	"time"
)

type NotificationHandler struct {
	BaseHandler
}

func NewNotificationHandler() *NotificationHandler {
	return &NotificationHandler{}
}

func (handler *NotificationHandler) Mount(mux *http.ServeMux) {
	mux.HandleFunc("/notification/health/", handler.CheckHealth)
	mux.HandleFunc("/notification/send/", handler.SendNotification)
}

func (handler *NotificationHandler) CheckHealth(resp http.ResponseWriter, req *http.Request) {
	resp.WriteHeader(http.StatusOK)
	resSchema := schema.NewHealthCheck()
	handler.MakeResponse(resp, resSchema)
}

func (handler *NotificationHandler) SendNotification(resp http.ResponseWriter, req *http.Request) {

	if req.Method != http.MethodPost {
		resp.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	sendNotificationReq := schema.SendNotificationRequest{}
	json.NewDecoder(req.Body).Decode(&sendNotificationReq)

	c := config.GetConfig(sendNotificationReq.Kind)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	err := services.SendNotificationService(ctx, sendNotificationReq.Kind, sendNotificationReq.Message, sendNotificationReq.Receiver, c)
	if err != nil {
		handler.MakeResponse(resp, schema.SendNotificationResponse{Success: false, Detail: err.Error()})
		return
	}

	handler.MakeResponse(resp, schema.SendNotificationResponse{Success: true, Detail: "message sent"})
}

func (handler *NotificationHandler) GetNotifications(resp http.ResponseWriter, req *http.Request) {}

func (handler *NotificationHandler) DeleteNotification(resp http.ResponseWriter, req *http.Request) {}
