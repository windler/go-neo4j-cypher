package neo4j

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type httpCypherClient struct {
	auth    string
	url     string
	verbose bool
}

//RequestExecuter executes a http.Request
type RequestExecuter interface {
	Do(req http.Request) (*http.Response, error)
}

//NewHTTPCypherClient creates a CypherClient using the buildin http client
func NewHTTPCypherClient(scheme, host string, port int64, user, password string) CypherClient {
	return &httpCypherClient{
		auth:    base64.StdEncoding.EncodeToString([]byte(user + ":" + password)),
		url:     scheme + host + ":" + strconv.FormatInt(port, 10) + "/db/data/transaction/commit",
		verbose: false,
	}
}

//Execute implements CypherClient
func (c *httpCypherClient) Execute(statement *CypherStatement) (ExecuteResult, error) {
	res, err := c.exec([]*CypherStatement{
		statement,
	})
	return res.Convert(), err
}

func (c *httpCypherClient) Verbose() {
	c.verbose = true
}

//ExecuteBatch implements CypherClient
func (c *httpCypherClient) ExecuteBatch(statements []*CypherStatement) (ExecuteBatchResult, error) {
	res, err := c.exec(statements)
	return res.ConvertBatch(), err
}

func (c *httpCypherClient) exec(data []*CypherStatement) (CypherQueryResult, error) {
	result := &CypherQueryResult{}

	payload := &CypherQueryPaylod{
		Statements: data,
	}
	reqData, err := json.Marshal(payload)
	if c.verbose {
		fmt.Println("Executing cypher:")
		fmt.Println(string(reqData[:]))
	}

	if err != nil {
		return *result, err
	}

	client := &http.Client{}
	req, _ := http.NewRequest("POST", c.url, bytes.NewBuffer(reqData))
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", "Basic "+c.auth)
	resp, err := client.Do(req)

	if err != nil {
		return *result, err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if c.verbose {
		fmt.Println("Got response:")
		fmt.Println(fmt.Println(string(body[:])))
	}

	if err != nil {
		return *result, err
	}

	json.Unmarshal(body, result)

	return *result, nil
}
