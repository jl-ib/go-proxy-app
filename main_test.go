package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"sync"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/jl-ib/proxy-app/api/middleware"	
	handlers "github.com/jl-ib/proxy-app/api/handlers"
	utils "github.com/jl-ib/proxy-app/api/utils"
	server "github.com/jl-ib/proxy-app/api/server"
)

func init() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		utils.LoadEnv()
		app := server.SetUp()
		handlers.HandlerRedirection(app)
		wg.Done()
		server.RunServer(app)
	}(wg)
	wg.Wait()
	fmt.Println("Server running...")
}

type Response struct {
	Status int `json:"status,omitempty"`
	Response string `json:"result,omitempty"`
	ResponseText []ResponseText `json:"res,omitempty"`
}

type ResponseText struct {
	Domain string
}

func TestAlgorithm(t *testing.T) {
	cases := []struct{
		Domain string
		Output string
	}{
		{Domain: "alpha", Output: "[\"alpha\"]"},
		{Domain: "beta", Output: "[\"alpha\",\"beta\"]"},
		{Domain: "omega", Output: "[\"alpha\",\"beta\",\"omega\"]"},
		{Domain: "beta", Output: "[\"alpha\",\"beta\",\"omega\",\"beta\"]"},
		{Domain: "", Output: "error"},
	}
	
	valuesToCompare := &Response{}
	client := http.Client{}

	for _, singleCase := range cases {
		req, err := http.NewRequest("GET", "http://localhost:8080/ping", nil)
		req.Header.Add("domain", singleCase.Domain)
	
		response, err := client.Do(req)
	
		bytes, err := ioutil.ReadAll(response.Body)

		err = json.Unmarshal(bytes, valuesToCompare)

		fmt.Println("response", valuesToCompare.Response, valuesToCompare.ResponseText, valuesToCompare.Status, singleCase.Output)
	
		assert.Nil(t, err)
		assert.Equal(t, valuesToCompare.Response, singleCase.Output)
	}
}

func TestSorting(t *testing.T) {
	unSortedCases := []*middleware.Queue {
		{Domain: "tetha", Weight: 3, Priority:4},
		{Domain: "beta", Weight: 5, Priority:1},
		{Domain: "omega", Weight: 1, Priority:5},
		{Domain: "phi", Weight: 2, Priority:1},
		{Domain: "alpha", Weight: 5, Priority:5},
	}

	sortedCases := []*middleware.Queue {
		{Domain: "alpha", Weight: 5, Priority:5},
		{Domain: "omega", Weight: 1, Priority:5},
		{Domain: "tetha", Weight: 3, Priority:4},
		{Domain: "beta", Weight: 5, Priority:1},
		{Domain: "phi", Weight: 2, Priority:1},
	}

	unSortedCases = middleware.PrioritizeQueue(unSortedCases)

	assert.Equal(t, sortedCases, unSortedCases)
}
