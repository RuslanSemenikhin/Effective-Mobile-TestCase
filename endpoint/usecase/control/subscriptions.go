package control

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func ListSubscriptions(
	c context.Context,
) {
	reqID := uuid.New().String()
	fmt.Printf("request ID - '%s'", reqID)
}
