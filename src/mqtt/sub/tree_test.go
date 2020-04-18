package sub

import (
	"github.com/google/uuid"
	"strings"
	"testing"
)

func TestTreeAdd(t *testing.T) {

	topic := "product1/device1/get"
	clientId := uuid.New()
	topics := strings.Split(topic, "/")

	AddTreeSub(topics, clientId)

}
