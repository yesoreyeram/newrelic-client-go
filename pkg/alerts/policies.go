package alerts

import (
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
	update(accountID int, policy Policy) (*Policy, error)
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

// QueryPolicy is similar to a Policy, but the resulting NerdGraph objects are
// string IDs in the JSON response.
type QueryPolicy struct {
	ID                 int                    `json:"id,string"`
	IncidentPreference IncidentPreferenceType `json:"incidentPreference"`
	Name               string                 `json:"name"`
}

type QueryPolicyInput struct {
	IncidentPreference IncidentPreferenceType `json:"incidentPreference"`
	Name               string                 `json:"name"`
}

type QueryPolicyCreateInput struct {
	QueryPolicyInput
}

type QueryPolicyUpdateInput struct {
	QueryPolicyInput
}

// ListPoliciesParams represents a set of filters to be used when querying New
// Relic alert policies.
type ListPoliciesParams struct {
	Name string `url:"filter[name],omitempty"`
}

// ListPolicies returns a list of Alert Policies for a given account.
func (a *Alerts) ListPolicies(params *ListPoliciesParams) ([]Policy, error) {
	var method PolicyCaller

	restCallerv := policyREST{
		parent: a,
	}

	method = &restCallerv

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
	return method.update(accountID, policy)
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

func (a *Alerts) CreatePolicyMutation(accountID int, policy QueryPolicyCreateInput) (*QueryPolicy, error) {
	vars := map[string]interface{}{
		"accountID": accountID,
		"policy":    policy,
	}

	resp := alertQueryPolicyCreateResponse{}

	if err := a.client.Query(alertsPolicyCreatePolicy, vars, &resp); err != nil {
		return nil, err
	}

	return &resp.QueryPolicy, nil
}

func (a *Alerts) UpdatePolicyMutation(accountID int, policyID int, policy QueryPolicyUpdateInput) (*QueryPolicy, error) {
	vars := map[string]interface{}{
		"accountID": accountID,
		"policyID":  policyID,
		"policy":    policy,
	}

	resp := alertQueryPolicyUpdateResponse{}

	if err := a.client.Query(alertsPolicyUpdatePolicy, vars, &resp); err != nil {
		return nil, err
	}

	return &resp.QueryPolicy, nil
}

// QueryPolicy queries NerdGraph for a policy matching the given account ID and
// policy ID.
func (a *Alerts) QueryPolicy(accountID, id int) (*QueryPolicy, error) {
	resp := alertQueryPolicyResponse{}
	vars := map[string]interface{}{
		"accountID": accountID,
		"policyID":  id,
	}

	if err := a.client.Query(alertPolicyQueryPolicy, vars, &resp); err != nil {
		return nil, err
	}

	return &resp.Actor.Account.Alerts.Policy, nil
}

// DeletePolicyMutation is the NerdGraph mutation to delete a policy given the
// account ID and the policy ID.
func (a *Alerts) DeletePolicyMutation(accountID, id int) (*QueryPolicy, error) {
	policy := &QueryPolicy{}

	resp := alertQueryPolicyDeleteRespose{}
	vars := map[string]interface{}{
		"accountID": accountID,
		"policyID":  id,
	}

	if err := a.client.Query(alertPolicyDeletePolicy, vars, &resp); err != nil {
		return nil, err
	}

	return policy, nil
}

type alertPoliciesResponse struct {
	Policies []Policy `json:"policies,omitempty"`
}

type alertPolicyResponse struct {
	Policy Policy `json:"policy,omitempty"`
}

type alertPolicyRequestBody struct {
	Policy Policy `json:"policy"`
}

type alertQueryPolicyCreateResponse struct {
	QueryPolicy QueryPolicy `json:"alertsPolicyCreate"`
}

type alertQueryPolicyUpdateResponse struct {
	QueryPolicy QueryPolicy `json:"alertsPolicyCreate"`
}

type alertQueryPolicyResponse struct {
	Actor struct {
		Account struct {
			Alerts struct {
				Policy QueryPolicy `json:"policy"`
			} `json:"alerts"`
		} `json:"account"`
	} `json:"actor"`
}

type alertQueryPolicyDeleteRespose struct {
	AlertsPolicyDelete struct {
		ID int `json:"id,string"`
	} `json:"alertsPolicyDelete"`
}

const (
	graphqlAlertPolicyFields = `
						id
						name
						incidentPreference
	`
	alertPolicyQueryPolicy = `query($accountID: Int!, $policyID: ID!) {
		actor {
			account(id: $accountID) {
				alerts {
					policy(id: $policyID) {` + graphqlAlertPolicyFields + `
					}
				}
			}
		}
	}`

	alertsPolicyCreatePolicy = `mutation CreatePolicy($accountID: Int!, $policy: AlertsPolicyInput!){
		alertsPolicyCreate(accountId: $accountID, policy: $policy) {` + graphqlAlertPolicyFields + `
		} }`

	alertsPolicyUpdatePolicy = `mutation UpdatePolicy($accountID: Int!, $policyID: ID!, $policy: AlertsPolicyUpdateInput!){
			alertsPolicyUpdate(accountId: $accountID, id: $policyID, policy: $policy) {` + graphqlAlertPolicyFields + `
			}
		}`

	alertPolicyDeletePolicy = `mutation DeletePolicy($accountID: Int!, $policyID: ID!){
		alertsPolicyDelete(accountId: $accountID, id: $policyID) {
			id
		} }`
)
