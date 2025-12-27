package dto

import "fmt"

type ConnectionMessageDTO struct {
	Message string `json:"message"`
}

func NewConnectionMessageDTO(message string, args ...any) ConnectionMessageDTO {
	return ConnectionMessageDTO{
		Message: fmt.Sprintf(message, args...),
	}
}
