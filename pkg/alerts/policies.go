package alerts

import (
	"fmt"

	"github.com/newrelic/newrelic-client-go/internal/serialization"
)

// IncidentPreferenceType specifies rollup settings for alert policies.
type IncidentPreferenceType string

var (
	// IncidentPreferenceTypes specifies the possible incident preferenece types for an alert policy.
	IncidentPreferenceTypes = struct {
		PerPolicy             IncidentPreferenceType
		PerCondition          IncidentPreferenceType
		PerConditionAndTarget IncidentPreferenceType
	}{
		PerPolicy:             "PER_POLICY",
		PerCondition:          "PER_CONDITION",
		PerConditionAndTarget: "PER_CONDITION_AND_TARGET",
	}
)

type PolicyCaller interface {
	// find(accountID int, name string) (*Policy, error)
	list(accountID int, params *ListPoliciesParams) ([]Policy, error)

	create(accountID int, policy Policy) (*Policy, error)
	get(accountID int, policyID int) (*Policy, error)
	update(accountID int, policyID int, policy Policy) (*Policy, error)
	remove(accountID int, policyID int) (*Policy, error) // delete is a reserved word...
}

// Policy represents a New Relic alert policy.
type Policy struct {
	ID                 int                      `json:"id,omitempty"`
	IncidentPreference IncidentPreferenceType   `json:"incident_preference,omitempty"`
	Name               string                   `json:"name,omitempty"`
	CreatedAt          *serialization.EpochTime `json:"created_at,omitempty"`
	UpdatedAt          *serialization.EpochTime `json:"updated_at,omitempty"`
}

// ListPoliciesParams represents a set of filters to be used when querying New
// Relic alert policies.
type ListPoliciesParams struct {
	Name string `url:"filter[name],omitempty"`
}

// ListPolicies returns a list of Alert Policies for a given account.
func (a *Alerts) ListPolicies(params *ListPoliciesParams) ([]Policy, error) {
	var method PolicyCaller

	restCaller := policyREST{
		parent: a,
	}

	nerdgraphCaller := policyNerdGraph{
		parent: a,
	}

	fmt.Printf("NerdGraph: %+v", nerdgraphCaller)

	method = &restCaller
	method = &nerdgraphCaller

	accountID := 0
	return method.list(accountID, params)
}

// GetPolicy returns a specific alert policy by ID for a given account.
func (a *Alerts) GetPolicy(id int) (*Policy, error) {
	var method PolicyCaller

	restCaller := policyREST{
		parent: a,
	}

	method = &restCaller

	accountID := 0
	return method.get(accountID, id)
}

// CreatePolicy creates a new alert policy for a given account.
func (a *Alerts) CreatePolicy(policy Policy) (*Policy, error) {
	var method PolicyCaller

	restCaller := policyREST{
		parent: a,
	}

	method = &restCaller

	accountID := 0
	return method.create(accountID, policy)
}

// UpdatePolicy update an alert policy for a given account.
func (a *Alerts) UpdatePolicy(policy Policy) (*Policy, error) {
	var method PolicyCaller

	restCaller := policyREST{
		parent: a,
	}

	method = &restCaller

	accountID := 0
	return method.update(accountID, policy.ID, policy)
}

// DeletePolicy deletes an existing alert policy for a given account.
func (a *Alerts) DeletePolicy(id int) (*Policy, error) {
	var method PolicyCaller

	restCaller := policyREST{
		parent: a,
	}

	method = &restCaller

	accountID := 0
	return method.remove(accountID, id)
}
