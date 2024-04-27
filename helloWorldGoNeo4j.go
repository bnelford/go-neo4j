package main

import (
	"fmt"
	"time"

	"github.com/neo4j/neo4j-go-driver/neo4j"
)

func mergeGraph(uri, username, password, name string) (string, error) {
	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		return "", err
	}
	defer driver.Close()

	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	update, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MERGE (a:Person {name: $name}) ON CREATE SET a.name = $name RETURN a.name + ', from node ' + id(a)",
			map[string]interface{}{"name": name})
		if err != nil {
			return nil, err
		}

		if result.Next() {
			return result.Record().Values[0], nil
		}

		return nil, result.Err()
	})
	if err != nil {
		return "", err
	}

	return update.(string), nil
}

func main() {
	start := time.Now()
	for i := 0; i < 1000; i++ {
		mergeGraph("bolt://localhost:7687", "neo4j", "password", (fmt.Sprintf("Beau%d", i)))
	}
	fmt.Println(time.Since(start))
}
