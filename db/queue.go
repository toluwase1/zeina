package db

import "zeina/models"

var queue []models.TransactionRequest

func PushRequestToQueue(queueRequest models.TransactionRequest) {
	queue = append(queue, queueRequest)
}

func GetAllRequestsFromQueue() []models.TransactionRequest {
	requests := make([]models.TransactionRequest, len(queue))
	copy(requests, queue)
	queue = []models.TransactionRequest{}
	return requests
}
