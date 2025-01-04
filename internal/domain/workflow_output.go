package domain

const WorkflowOutputCertificate = "certificate"

type WorkflowOutput struct {
	Meta
	WorkflowId string           `json:"workflowId" db:"workflowId"`
	NodeId     string           `json:"nodeId" db:"nodeId"`
	Node       *WorkflowNode    `json:"node" db:"node"`
	Outputs    []WorkflowNodeIO `json:"outputs" db:"outputs"`
	Succeeded  bool             `json:"succeeded"db:"succeeded"`
}
