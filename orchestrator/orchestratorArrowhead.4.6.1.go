package orchestrator

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/MrDweller/orchestrator-connection/models"
)

const ORCHESTRATION_ARROWHEAD_4_6_1 OrchestratorImplementationType = "orchestration-arrowhead-4.6.1"

type OrchestrationArrowhead_4_6_1 struct {
	Orchestrator
	models.CertificateInfo
}

type OrchestrationDTO struct {
	RequesterSystem    models.SystemDefinition `json:"requesterSystem"`
	RequestedService   RequestedService        `json:"requestedService"`
	OrchestrationFlags map[string]bool         `json:"orchestrationFlags"`
}

type RequestedService struct {
	InterfaceRequirements        []string `json:"interfaceRequirements"`
	ServiceDefinitionRequirement string   `json:"serviceDefinitionRequirement"`
}

func (orchestrator OrchestrationArrowhead_4_6_1) Connect() error {
	result, err := orchestrator.echoOrchestrator()
	if err != nil {
		return err
	}

	if string(result) != "Got it!" {
		return errors.New("can't establish a connection with the orchestrator")
	}

	return nil
}

func (orchestrator OrchestrationArrowhead_4_6_1) Orchestration(requestedService models.ServiceDefinition, requesterSystem models.SystemDefinition) (*models.OrchestrationResponse, error) {
	orchestrationDTO := OrchestrationDTO{
		RequesterSystem: requesterSystem,
		RequestedService: RequestedService{
			InterfaceRequirements: []string{
				"HTTP-SECURE-JSON",
			},
			ServiceDefinitionRequirement: requestedService.ServiceDefinition,
		},
	}
	payload, err := json.Marshal(orchestrationDTO)

	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", "https://"+orchestrator.Address+":"+strconv.Itoa(orchestrator.Port)+"/orchestrator/orchestration", bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client, err := orchestrator.getClient()
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		errorString := fmt.Sprintf("status: %s, body: %s", resp.Status, string(body))
		return nil, errors.New(errorString)
	}

	var orchestrationResponse models.OrchestrationResponse
	json.Unmarshal(body, &orchestrationResponse)
	return &orchestrationResponse, nil
}

func (orchestrator OrchestrationArrowhead_4_6_1) echoOrchestrator() ([]byte, error) {
	req, err := http.NewRequest("GET", "https://"+orchestrator.Address+":"+strconv.Itoa(orchestrator.Port)+"/orchestrator/echo", nil)
	if err != nil {
		return nil, err
	}

	client, err := orchestrator.getClient()
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (orchestrator OrchestrationArrowhead_4_6_1) getClient() (*http.Client, error) {
	cert, err := tls.LoadX509KeyPair(orchestrator.CertFilePath, orchestrator.KeyFilePath)
	if err != nil {
		return nil, err
	}

	// Load truststore.p12
	truststoreData, err := os.ReadFile(orchestrator.Truststore)
	if err != nil {
		return nil, err

	}

	// Extract the root certificate(s) from the truststore
	pool := x509.NewCertPool()
	if ok := pool.AppendCertsFromPEM(truststoreData); !ok {
		return nil, err
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates:       []tls.Certificate{cert},
				RootCAs:            pool,
				InsecureSkipVerify: false,
			},
		},
	}
	return client, nil
}
