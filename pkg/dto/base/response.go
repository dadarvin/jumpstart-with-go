package base

// Error str
type Error struct {
	IsError bool   `json:"isError"` //JWT
	Message string `json:"message"`
}

// JsonResponse struct for JSON response
type JsonResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// JSONErrorResponse struct for JSON error response
type JsonErrorResponse struct {
	Error *ApiResponse `json:"error"`
}

// ApiResponse response for Api request
type ApiResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
