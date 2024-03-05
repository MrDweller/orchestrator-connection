package orchestrator

import (
	"errors"
	"fmt"

	"github.com/MrDweller/orchestrator-connection/models"
)

type OrchestratorConnection interface {
	Connect() error
	Orchestration(requestedService string, requesterSystem models.SystemDefinition, additionalParameters any) (*models.OrchestrationResponse, error)
}

type OrchestratorImplementationType string

func NewConnection(orchestrator Orchestrator, orchestratorImplementationType OrchestratorImplementationType, certificateInfo models.CertificateInfo) (OrchestratorConnection, error) {
	var orchestratorConnection OrchestratorConnection

	switch orchestratorImplementationType {
	case ORCHESTRATION_ARROWHEAD_4_6_1:
		orchestratorConnection = OrchestrationArrowhead_4_6_1{
			Orchestrator:    orchestrator,
			CertificateInfo: certificateInfo,
		}
		break
	default:
		errorString := fmt.Sprintf("the service registry %s has no implementation", orchestratorImplementationType)
		return nil, errors.New(errorString)
	}

	err := orchestratorConnection.Connect()
	if err != nil {
		return nil, err
	}

	return orchestratorConnection, nil
}
