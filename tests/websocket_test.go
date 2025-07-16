package tests

import (
	"Complaingo/internal/domain/models"
	"Complaingo/internal/middleware"
	websockets "Complaingo/internal/websockets"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

// Mock Kafka Producer
type MockKafkaProducer struct{}

func (m *MockKafkaProducer) SendMessage(msg string) {}

// Mock Message Repo
type MockMessageRepo struct {
	Saved []*models.MessageEntity
}

func (m *MockMessageRepo) SaveMessage(ctx context.Context, msg *models.MessageEntity) error {
	m.Saved = append(m.Saved, msg)
	return nil
}

func TestWebSocketEndToEnd(t *testing.T) {
	mockRepo := &MockMessageRepo{}
	mockKafka := &MockKafkaProducer{}
	handler := websockets.NewwebsocketHandler(mockRepo, mockKafka)

	// 2: Create a test server
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, middleware.ContextUserID, 123)
		ctx = context.WithValue(ctx, middleware.ContextRole, "user")
		r = r.WithContext(ctx)
		handler.HandleWebsocket(w, r)
	}))
	defer srv.Close()

	// 3: Connect WebSocket client
	wsURL := "ws" + srv.URL[len("http"):] // e.g., ws://127.0.0.1:12345
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	assert.NoError(t, err)
	defer conn.Close()

	// 4: Subscribe
	subMsg := websockets.Message{
		Type:    "subscribe",
		Channel: "test-channel",
	}
	err = conn.WriteJSON(subMsg)
	assert.NoError(t, err)

	// 5: Publish to channel
	pubMsg := websockets.Message{
		Type:    "publish",
		Channel: "test-channel",
		From:    "123",
		Message: "Hello Test Channel!",
	}
	err = conn.WriteJSON(pubMsg)
	assert.NoError(t, err)

	// Set read deadline to avoid hanging
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))

	var received websockets.Message
	err = conn.ReadJSON(&received)
	assert.NoError(t, err)
	assert.Equal(t, pubMsg.Message, received.Message)

	// 6: Direct message to self
	directMsg := websockets.Message{
		Type:    "direct",
		From:    "123",
		To:      "123",
		Message: "Hello Me!",
	}
	err = conn.WriteJSON(directMsg)
	assert.NoError(t, err)

	conn.SetReadDeadline(time.Now().Add(2 * time.Second))

	var directReceived websockets.Message
	err = conn.ReadJSON(&directReceived)
	assert.NoError(t, err)
	assert.Equal(t, directMsg.Message, directReceived.Message)

	//7: Verify message was saved
	time.Sleep(500 * time.Millisecond)
	assert.Equal(t, 1, len(mockRepo.Saved))
	assert.Equal(t, "Hello Me!", mockRepo.Saved[0].Message)
}
