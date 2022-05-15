package app

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	EventID  uuid.UUID
	UserID   uuid.UUID
	Title    string
	DateTime time.Time
}

func (n Notification) String() string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("New notification: %s at %s", n.Title, n.DateTime.Format(time.RFC3339)))
	return builder.String()
}
