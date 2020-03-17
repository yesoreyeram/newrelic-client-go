package alerts

import "fmt"

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

// policyNerdGraph implements the PolicyCaller interface.
type policyNerdGraph struct {
	parent *Alerts
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

func policyToQueryPolicy(policy Policy) QueryPolicy {
	return QueryPolicy{
		ID:                 policy.ID,
		IncidentPreference: policy.IncidentPreference,
		Name:               policy.Name,
	}
}

func queryPolicyToPolicy(policy QueryPolicy) Policy {
	return Policy{
		ID:                 policy.ID,
		IncidentPreference: policy.IncidentPreference,
		Name:               policy.Name,
	}
}

func (p *policyNerdGraph) list(accountID int, params *ListPoliciesParams) ([]Policy, error) {
	policies := []Policy{}

	return policies, fmt.Errorf("list not implemented")
}

func (p *policyNerdGraph) create(accountID int, policy Policy) (*Policy, error) {
	vars := map[string]interface{}{
		"accountID": accountID,
		"policy":    policyToQueryPolicy(policy),
	}

	resp := alertQueryPolicyCreateResponse{}

	if err := p.parent.client.Query(alertsPolicyCreatePolicy, vars, &resp); err != nil {
		return nil, err
	}

	returnPolicy := queryPolicyToPolicy(resp.QueryPolicy)
	return &returnPolicy, nil
}

func (p *policyNerdGraph) update(accountID int, policyID int, policy Policy) (*Policy, error) {
	vars := map[string]interface{}{
		"accountID": accountID,
		"policyID":  policyID,
		"policy":    policyToQueryPolicy(policy),
	}

	resp := alertQueryPolicyUpdateResponse{}

	if err := p.parent.client.Query(alertsPolicyUpdatePolicy, vars, &resp); err != nil {
		return nil, err
	}

	returnPolicy := queryPolicyToPolicy(resp.QueryPolicy)
	return &returnPolicy, nil
}

// QueryPolicy queries NerdGraph for a policy matching the given account ID and
// policy ID.
func (p *policyNerdGraph) get(accountID, id int) (*Policy, error) {
	resp := alertQueryPolicyResponse{}
	vars := map[string]interface{}{
		"accountID": accountID,
		"policyID":  id,
	}

	if err := p.parent.client.Query(alertPolicyQueryPolicy, vars, &resp); err != nil {
		return nil, err
	}

	returnPolicy := queryPolicyToPolicy(resp.Actor.Account.Alerts.Policy)
	return &returnPolicy, nil
}

// DeletePolicyMutation is the NerdGraph mutation to delete a policy given the
// account ID and the policy ID.
func (p *policyNerdGraph) remove(accountID, id int) (*Policy, error) {
	policy := &QueryPolicy{}

	resp := alertQueryPolicyDeleteRespose{}
	vars := map[string]interface{}{
		"accountID": accountID,
		"policyID":  id,
	}

	if err := p.parent.client.Query(alertPolicyDeletePolicy, vars, &resp); err != nil {
		return nil, err
	}

	returnPolicy := queryPolicyToPolicy(*policy)
	return &returnPolicy, nil
}
