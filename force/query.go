package force

import (
	"fmt"
	"net/url"
	"strings"
)

const (
	baseQueryString = "SELECT %v FROM %v"
)

// BuildQuery is a func that is used to build a query
func BuildQuery(fields, table string, constraints []string) string {
	query := fmt.Sprintf(baseQueryString, fields, table)
	if len(constraints) > 0 {
		query += fmt.Sprintf(" WHERE %v", strings.Join(constraints, " AND "))
	}

	return query
}

// Query Use the Query resource to execute a SOQL query that returns all the results in a single response,
// or if needed, returns part of the results and an identifier used to retrieve the remaining results.
func (forceAPI *API) Query(query string, out interface{}) (err error) {
	uri := forceAPI.apiResources[queryKey]

	params := url.Values{
		"q": {query},
	}

	err = forceAPI.Get(uri, params, out)

	return
}

// QueryAll Use the QueryAll resource to execute a SOQL query that includes information about records that have
// been deleted because of a merge or delete. Use QueryAll rather than Query, because the Query resource
// will automatically filter out items that have been deleted.
func (forceAPI *API) QueryAll(query string, out interface{}) (err error) {
	uri := forceAPI.apiResources[queryAllKey]

	params := url.Values{
		"q": {query},
	}

	err = forceAPI.Get(uri, params, out)

	return
}

// QueryNext is a func that does a get
func (forceAPI *API) QueryNext(uri string, out interface{}) (err error) {
	err = forceAPI.Get(uri, nil, out)

	return
}
