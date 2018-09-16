# go neo4j cypher http api client
`Golang` HTTP api client to access [neo4j](https://neo4j.com/) [transactional cypher endpoint](https://neo4j.com/docs/developer-manual/3.4/http-api/#http-api-transactional).

## Tested neo4j versions

|Version|Success?|
|-|-|
|3.4|yes|

## Installation 
```bash
go get github.com/windler/go-neo4j-cypher
```

## Usage
```go
import neo4j "github.com/windler/go-neo4j-cypher"

func query() {
   client := neo4j.NewHTTPCypherClient("http://", "neofj-host", 7474, "myuser", "secret")
   result, err := client.Execute(&neo4j.CypherStatement{
        Statement:  `MATCH (s) 
        WHERE s.name = {name}
        return s.name as nodeName`,
		Parameters: map[string]interface{}{
            "name": "Alfred",
        },
    })
    
    if err != nil {
        panic(err.Error())
    }

    if len(result.Errors) > 0 {
        //handle Errors
    }

    // optional: map result to string slice
    strResult := result.Map("nodeName", func(rowValue interface{}, meta neo4j.CypherQueryResultValueMeta) interface{} {
		return rowValue.(string)
    })

    //... handle strResult ... 
}

```