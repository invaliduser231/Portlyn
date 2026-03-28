package http

import stdhttp "net/http"

const nodeAgentSource = `package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type heartbeatPayload struct {
	Version          string  ` + "`json:\"version,omitempty\"`" + `
	Status           string  ` + "`json:\"status,omitempty\"`" + `
	Load             float64 ` + "`json:\"load,omitempty\"`" + `
	BandwidthInKbps  float64 ` + "`json:\"bandwidth_in_kbps,omitempty\"`" + `
	BandwidthOutKbps float64 ` + "`json:\"bandwidth_out_kbps,omitempty\"`" + `
}

type enrollRequest struct {
	Token       string ` + "`json:\"token\"`" + `
	Name        string ` + "`json:\"name\"`" + `
	Description string ` + "`json:\"description,omitempty\"`" + `
	Version     string ` + "`json:\"version,omitempty\"`" + `
}

type enrollResponse struct {
	Node struct {
		ID uint ` + "`json:\"id\"`" + `
	} ` + "`json:\"node\"`" + `
	HeartbeatToken string ` + "`json:\"heartbeat_token\"`" + `
	HeartbeatURL   string ` + "`json:\"heartbeat_url\"`" + `
}

func main() {
	apiBase := flag.String("api", "http://localhost:8080", "api base url")
	endpoint := flag.String("endpoint", "", "heartbeat endpoint, auto-derived after enrollment if omitted")
	token := flag.String("token", "", "enrollment token or existing node bearer token")
	nodeID := flag.Uint("node-id", 0, "existing node id when using a pre-issued heartbeat token")
	name := flag.String("name", "node-agent", "node name used during enrollment")
	description := flag.String("description", "", "node description used during enrollment")
	version := flag.String("version", "dev-agent", "node version")
	interval := flag.Duration("interval", 30*time.Second, "heartbeat interval")
	flag.Parse()

	client := &http.Client{Timeout: 10 * time.Second}
	heartbeatToken := *token
	heartbeatEndpoint := strings.TrimSpace(*endpoint)

	if heartbeatEndpoint == "" {
		if *nodeID > 0 {
			heartbeatEndpoint = strings.TrimRight(*apiBase, "/") + "/api/v1/nodes/" + fmt.Sprintf("%d", *nodeID) + "/heartbeat"
		} else if *token != "" {
			enrolledEndpoint, enrolledToken, err := enrollNode(client, strings.TrimRight(*apiBase, "/"), *token, *name, *description, *version)
			if err != nil {
				log.Fatalf("enroll node: %v", err)
			}
			heartbeatEndpoint = enrolledEndpoint
			heartbeatToken = enrolledToken
		}
	}
	if heartbeatEndpoint == "" || heartbeatToken == "" {
		log.Fatal("provide either --endpoint and --token, or --api and an enrollment --token")
	}

	for {
		payload := heartbeatPayload{
			Version:          *version,
			Status:           "online",
			Load:             0.1,
			BandwidthInKbps:  128,
			BandwidthOutKbps: 64,
		}

		body, err := json.Marshal(payload)
		if err != nil {
			log.Printf("marshal heartbeat: %v", err)
			time.Sleep(*interval)
			continue
		}

		req, err := http.NewRequest(http.MethodPost, heartbeatEndpoint, bytes.NewReader(body))
		if err != nil {
			log.Printf("build request: %v", err)
			time.Sleep(*interval)
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+heartbeatToken)

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("heartbeat failed: %v", err)
			time.Sleep(*interval)
			continue
		}
		_ = resp.Body.Close()
		log.Printf("heartbeat status=%s", resp.Status)
		time.Sleep(*interval)
	}
}

func enrollNode(client *http.Client, apiBase, token, name, description, version string) (string, string, error) {
	payload := enrollRequest{
		Token:       token,
		Name:        name,
		Description: description,
		Version:     version,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return "", "", err
	}
	req, err := http.NewRequest(http.MethodPost, apiBase+"/api/v1/nodes/enroll", bytes.NewReader(body))
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return "", "", fmt.Errorf("enroll failed: %s", resp.Status)
	}
	var result enrollResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", err
	}
	return strings.TrimRight(apiBase, "/") + result.HeartbeatURL, result.HeartbeatToken, nil
}
`

func (s *Server) handleNodeAgentSource(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Content-Disposition", `attachment; filename="nodeagent.go"`)
	_, _ = w.Write([]byte(nodeAgentSource))
}
