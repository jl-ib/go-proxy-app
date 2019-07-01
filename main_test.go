package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"sync"
	"testing"
	"github.com/stretchr/testify/assert"
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
	fmt.Print("Server running...")
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
		// {Domain: "omega", Output: "[\"alpha\"]"},
		{Domain: "", Output: "domain error"},
	}
	
	valuesToCompare := &Response{}
	client := http.Client{}

	for _, singleCase := range cases {
		req, err := http.NewRequest("GET", "http://localhost:8080/ping", nil)
		req.Header.Add("domain", singleCase.Domain)
	
		response, err := client.Do(req)
	
		bytes, err := ioutil.ReadAll(response.Body)

		err = json.Unmarshal(bytes, valuesToCompare)

		fmt.Println(valuesToCompare.Response)
	
		assert.Nil(t, err)
		assert.Equal(t, valuesToCompare.Response, singleCase.Output)
	}
}
