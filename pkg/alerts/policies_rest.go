package alerts

import (
	"fmt"

	"github.com/newrelic/newrelic-client-go/pkg/errors"
)

type policyREST struct {
	parent *Alerts
}

func (p *policyREST) list(accountID int, params *ListPoliciesParams) ([]Policy, error) {
	alertPolicies := []Policy{}

	nextURL := "/alerts_policies.json"

	for nextURL != "" {
		response := alertPoliciesResponse{}
		resp, err := p.parent.client.Get(nextURL, &params, &response)
		if err != nil {
			return nil, err
		}

		alertPolicies = append(alertPolicies, response.Policies...)

		paging := p.parent.pager.Parse(resp)
		nextURL = paging.Next
	}

	return alertPolicies, nil
}

// GetPolicy returns a specific alert policy by ID for a given account.
func (p *policyREST) get(accountID, id int) (*Policy, error) {
	policies, err := p.list(accountID, &ListPoliciesParams{})

	if err != nil {
		return nil, err
	}

	for _, policy := range policies {
		if policy.ID == id {
			return &policy, nil
		}
	}

	return nil, errors.NewNotFoundf("no alert policy found for id %d", id)
}

// CreatePolicy creates a new alert policy for a given account.
func (p *policyREST) create(accountID int, policy Policy) (*Policy, error) {
	reqBody := alertPolicyRequestBody{
		Policy: policy,
	}
	resp := alertPolicyResponse{}

	_, err := p.parent.client.Post("/alerts_policies.json", nil, &reqBody, &resp)

	if err != nil {
		return nil, err
	}

	return &resp.Policy, nil
}

// update an alert policy for a given account.
func (p *policyREST) update(accountID int, policyID int, policy Policy) (*Policy, error) {

	reqBody := alertPolicyRequestBody{
		Policy: policy,
	}
	resp := alertPolicyResponse{}
	url := fmt.Sprintf("/alerts_policies/%d.json", policy.ID)

	_, err := p.parent.client.Put(url, nil, &reqBody, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.Policy, nil
}

// remove an existing alert policy for a given account.
func (p *policyREST) remove(accountID, id int) (*Policy, error) {
	resp := alertPolicyResponse{}
	url := fmt.Sprintf("/alerts_policies/%d.json", id)

	_, err := p.parent.client.Delete(url, nil, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.Policy, nil
}
