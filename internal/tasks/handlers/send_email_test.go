package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendEmail(t *testing.T) {
	clientName := "Client"
	clientEmail := "bohomolyulian3022003@gmail.com"

	recipientName := "TestUser"
	recipientEmail := "elliotaldersonhome@gmail.com"

	err := SendEmail(clientName, clientEmail, recipientName, recipientEmail, "Test", "<html><body><p>Test</p></body></html>")
	assert.NoError(t, err)
}
