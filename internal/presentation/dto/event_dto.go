package dto

import (
	"encoding/json"
	"fmt"
	"io"
)

type ConnectionMessageDTO struct {
	Message string `json:"message"`
}

func PrintNewConnectionMessage(w io.Writer, message string, args ...any) error {
	messageDTO := ConnectionMessageDTO{
		Message: fmt.Sprintf(message, args...),
	}

	data, err := json.Marshal(messageDTO)
	if err != nil {
		return fmt.Errorf("error marshalling message: %w", err)
	}

	if _, err := fmt.Fprintf(w, "%s\n", string(data)); err != nil {
		return fmt.Errorf("error sending message: %w", err)
	}

	return nil
}
