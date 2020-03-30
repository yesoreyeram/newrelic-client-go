//go:generate ../../bin/typegen -v

// Package nerdstorage provides access to the API for New Relic NerdStorage.
package nerdstorage

import (
	"github.com/newrelic/newrelic-client-go/internal/http"
	"github.com/newrelic/newrelic-client-go/internal/logging"
	"github.com/newrelic/newrelic-client-go/pkg/config"
)

// NerdStorage is used to communicate with the New Relic Nerdstorage product.
type NerdStorage struct {
	client http.Client
	logger logging.Logger
}

// New returns a new client for interacting with New Relic Nerdstorage.
func New(config config.Config) NerdStorage {
	return NerdStorage{
		client: http.NewClient(config),
		logger: config.GetLogger(),
	}
}

func (n *NerdStorage) QueryCollection(accountID int, collection string) {
	resp := nerdStorageQueryCollectionResponse{}

	// NewRelic-Package-ID

	vars := map[string]interface{}{
		"accountID":  accountID,
		"collection": collection,
	}

	if err := n.client.Query(collectionQuery, vars, &resp); err != nil {
		return nil, err
	}

	return &resp.Actor.Account.Alerts.Policy, nil
}

const (
	nerdStorageWriteDocument = `
		mutation($collection: String!, $document: NerdStorageDocument!, $documentId: String!, $scope: NerdStorageScopeInput!) {
			writeDocument(collection: $colletction, document: $document, documentId: $documentId, scope: $scope) { }
		}
		`

	nerdStorageDeleteDocument = `
		mutation($collection: String!, $documentId: String!, $scope: NerdStorageScopeInput!) {
			deleteDocument(collection: $colletction, document: $document, documentId: $documentId, scope: $scope) {
				deleted
		} }
		`

	nerdStorageDeleteCollection = `
		mutation($collection: String!, $documentId: String!, $scope: NerdStorageScopeInput!) {
			deleteCollection(collection: $colletction, scope: $scope) {
				deleted
		} }
		`

	collectionQuery = `
		query($accountID: Int!, $collection: String!) {
			actor {
				account(id: $accountID) {
					nerdStorage {
						collection(collection: $collection) {
							document
							id
						}
					}
				}
			}
		}
		`

	documentQuery = `
		query($collection: String!, $documentId: String!) {
			actor {
				nerdStorage {
					document(collection: $collection, documentId: $documentId)
				}
			}
		}
		`
)
