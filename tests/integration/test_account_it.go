package integration

import (
	"net/http"
	"testing"
)
//TODO complete it
func TestGetAccountIntegration(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/api/v1/accounts/12")
	if err != nil {
		t.Fatalf("Failed to send GET request: %v", err)
	}
	defer resp.Body.Close()

	t.Logf("Status Code: %d", resp.StatusCode)
}
