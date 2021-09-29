package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	dapr "github.com/dapr/go-sdk/client"
)

// QueryResponse is the response object for querying state.
type QueryResponse struct {
	Results []QueryResult `json:"results"`
	Token   string        `json:"token,omitempty"`
}

// QueryResult is an object representing a single entry in query result.
type QueryResult struct {
	Key   string  `json:"key"`
	Data  []byte  `json:"data"`
	ETag  *string `json:"etag,omitempty"`
	Error string  `json:"error,omitempty"`
}

func main() {
	var (
		port, store, op, param string
	)
	flag.StringVar(&store, "s", "statestore", "statestore name")
	flag.StringVar(&port, "p", "3500", "dapr port")
	flag.StringVar(&op, "o", "", "operation (set/get/query)")
	flag.StringVar(&param, "i", "", "operation input parameter")
	flag.Parse()

	// create the client
	client, err := dapr.NewClientWithPort(port)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	if err = run(client, port, store, op, param); err != nil {
		panic(err)
	}
}

func run(client dapr.Client, port, store, op, param string) error {
	switch op {

	case "set":
		fmt.Printf("Set JSON array of objects from %s to %s\n", param, store)

		content, err := os.ReadFile(param)
		if err != nil {
			return fmt.Errorf("failed to read input file: %v", err)
		}

		resp, err := http.Post(fmt.Sprintf("http://localhost:%s/v1.0/state/%s", port, store), "application/json", bytes.NewBuffer(content))
		if err != nil {
			return fmt.Errorf("failed in http.Post: %v", err)
		}
		defer resp.Body.Close()

		_, err = io.Copy(os.Stdout, resp.Body)
		if err != nil {
			return fmt.Errorf("failed to copy: %v", err)
		}
		fmt.Println("")

	case "get":
		fmt.Printf("Get object with key %s from %s\n", param, store)

		resp, err := http.Get(fmt.Sprintf("http://localhost:%s/v1.0/state/%s/%s", port, store, param))
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		val, err := bytes2json(data)
		if err != nil {
			return err
		}
		fmt.Println(val)

	case "query":
		fmt.Printf("Query objects in %s\n", store)

		content, err := os.ReadFile(param)
		if err != nil {
			return err
		}

		resp, err := http.Post(fmt.Sprintf("http://localhost:%s/v1.0/state/%s/query", port, store), "application/json", bytes.NewBuffer(content))
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		var qr QueryResponse
		err = json.Unmarshal(body, &qr)
		if err != nil {
			return err
		}

		fmt.Println("Result:")
		for _, item := range qr.Results {
			val, err := bytes2json(item.Data)
			if err != nil {
				return err
			}
			fmt.Println("KEY:", item.Key)
			fmt.Println(val)
		}
		fmt.Println("Token:", qr.Token)

	default:
		return fmt.Errorf("unsupported operation %q", op)
	}
	fmt.Println("Done")
	return nil
}

func bytes2json(data []byte) (string, error) {
	var v interface{}
	err := json.Unmarshal(data, &v)
	if err != nil {
		return "", err
	}
	data, err = json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}
