package neo4j

//CypherClient performs cypher queries against neo4js transactional cypher HTTP endpoint
type CypherClient interface {
	Execute(statement *CypherStatement) (ExecuteResult, error)
	ExecuteBatch(statements []*CypherStatement) (ExecuteBatchResult, error)
}

//CypherQueryPaylod represents the payload for http transactional api
type CypherQueryPaylod struct {
	Statements []*CypherStatement `json:"statements"`
}

//CypherStatement represents a cypher query
type CypherStatement struct {
	Statement  string           `json:"statement"`
	Parameters CypherParameters `json:"parameters"`
}

//CypherParameters are convenient to use with statements.
type CypherParameters map[string]interface{}

//CypherQueryResult is neo4js api response
type CypherQueryResult struct {
	Results []CypherQueryResultValue `json:"results"`
	Errors  []CypherQueryResultError `json:"errors"`
}

//CypherQueryResultValue is neo4js api statement result response
type CypherQueryResultValue struct {
	Columns []string                     `json:"columns"`
	Data    []CypherQueryResultValueData `json:"data"`
}

//CypherQueryResultValueData is neo4js api data response
type CypherQueryResultValueData struct {
	Row  []interface{}                `json:"row"`
	Meta []CypherQueryResultValueMeta `json:"meta"`
}

//CypherQueryResultValueMeta is neo4js api metadata response
type CypherQueryResultValueMeta struct {
	ID      int    `json:"id"`
	Type    string `json:"type"`
	Deleted bool   `json:"deleted"`
}

//CypherQueryResultError is neo4js api error response
type CypherQueryResultError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

//ExecuteBatchResult represents a convenient result of CypherQueryResult
type ExecuteBatchResult struct {
	ColumnsResults []map[string][]ExecuteResultRow
	Errors         []CypherQueryResultError
}

//ExecuteResult represents a convenient result of CypherQueryResult assuming only one statement was queried
type ExecuteResult struct {
	ColumnsResults map[string][]ExecuteResultRow
	Errors         []CypherQueryResultError
}

//ExecuteResultRow  represents a convenient result for each statement
type ExecuteResultRow struct {
	Row  interface{}
	Meta CypherQueryResultValueMeta
}

//Walk plucks the columns index and applies fn to every columns data in CypherQueryResultValue
func (r *CypherQueryResultValue) Walk(column string, fn func(data interface{}, meta CypherQueryResultValueMeta)) {
	idx := 0
	for i := 0; i < len(r.Columns); i++ {
		if r.Columns[i] == column {
			idx = i
			break
		}
	}

	for _, d := range r.Data {
		fn(d.Row[idx], d.Meta[idx])
	}
}

//ConvertBatch converts CypherQueryResult to ExecuteBatchResult
func (r *CypherQueryResult) ConvertBatch() ExecuteBatchResult {
	execResult := []map[string][]ExecuteResultRow{}

	for _, d := range r.Results {
		curResult := map[string][]ExecuteResultRow{}
		for _, c := range d.Columns {
			data := []ExecuteResultRow{}
			d.Walk(c, func(row interface{}, meta CypherQueryResultValueMeta) {
				data = append(data, ExecuteResultRow{
					Meta: meta,
					Row:  row,
				})
			})
			curResult[c] = data
		}
		execResult = append(execResult, curResult)
	}

	return ExecuteBatchResult{
		Errors:         r.Errors,
		ColumnsResults: execResult,
	}
}

//Convert converts CypherQueryResult to ExecuteResult
func (r *CypherQueryResult) Convert() ExecuteResult {
	execResult := map[string][]ExecuteResultRow{}

	if len(r.Results) == 1 {
		d := r.Results[0]
		curResult := map[string][]ExecuteResultRow{}
		for _, c := range d.Columns {
			data := []ExecuteResultRow{}
			d.Walk(c, func(row interface{}, meta CypherQueryResultValueMeta) {
				data = append(data, ExecuteResultRow{
					Meta: meta,
					Row:  row,
				})
			})
			curResult[c] = data
		}
		execResult = curResult
	}

	return ExecuteResult{
		Errors:         r.Errors,
		ColumnsResults: execResult,
	}
}

//Map plucks all values for column and maps the interface{} using fn. Map return a new slice containing fns results.
func (r *ExecuteResult) Map(column string, fn func(rowValue interface{}, meta CypherQueryResultValueMeta) interface{}) []interface{} {
	res := []interface{}{}
	for c, v := range r.ColumnsResults {
		if c == column {
			for _, row := range v {
				res = append(res, fn(row.Row, row.Meta))
			}
			break
		}
	}
	return res
}
