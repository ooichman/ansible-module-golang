package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"bytes"
	"strconv"
	"net/http"
)

type ModuleArgs struct {
	Ssl bool 
	Target_port int 
	Protocol string
	Target_fqdn string
	Spearedge_fqdn string
	Spearedge_path string
	Openshift_node string
	Name string
	
}

type SpearedgeInput struct {
    Port string `json:"port"`
    Target string `json:"target"`
    Protocol string `json:"protocol"`
    Hostname string `json:hostname`
}

type Response struct {
	busy    bool   `json:busy`
	Changed bool   `json:"changed"`
	Failed  bool   `json:"failed"`
	Msg     string `json:"msg"`
}

func prettyprint(b []byte) ([]byte, error) {
    var out bytes.Buffer
    err := json.Indent(&out, b, "", "  ")
    return out.Bytes(), err
}

func ExitJson(responseBody Response) {
	returnResponse(responseBody)
}

func FailJson(responseBody Response) {
	responseBody.Failed = true
	returnResponse(responseBody)
}

func returnResponse(responseBody Response) {
	var response []byte
	var err error
	response, err = json.Marshal(responseBody)
	if err != nil {
		response, _ = json.Marshal(Response{Msg: "Invalid response object"})
	}
	fmt.Println(string(response))
	if responseBody.Failed {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}

func main() {
	var response Response

	if len(os.Args) != 2 {
		response.Msg = "No argument file provided"
		FailJson(response)
	}

	argsFile := os.Args[1]

	text, err := ioutil.ReadFile(argsFile)
	if err != nil {
		response.Msg = "Could not read configuration file: " + argsFile
		FailJson(response)
	}

	var moduleArgs ModuleArgs
	err = json.Unmarshal(text, &moduleArgs)
	if err != nil {
		response.Msg = "Configuration file not valid JSON: " + argsFile
		FailJson(response)
	}


	var spearedge_url string

	if moduleArgs.Ssl == true {
			spearedge_url = "https://"
	} else {
		spearedge_url = "http://"
	}
	
	if moduleArgs.Spearedge_fqdn == "" {
		response.Msg = "Undefined spearedge FQDN in the values"
		FailJson(response)
	} else {
		spearedge_url = spearedge_url + moduleArgs.Spearedge_fqdn
	}

    if moduleArgs.Spearedge_path == "" {
		response.Msg = "Undefined spearedge PATH in the values"
		FailJson(response)
	} else {
		spearedge_url = spearedge_url + "/" + moduleArgs.Spearedge_path
	}

	var s_input SpearedgeInput

		s_input.Port = strconv.Itoa(moduleArgs.Target_port)
		s_input.Protocol = moduleArgs.Protocol
		s_input.Hostname = moduleArgs.Openshift_node
		s_input.Target = moduleArgs.Target_fqdn

		sendbody , _ := json.Marshal(s_input)

		url_response , error := http.Post(spearedge_url, "application/json", bytes.NewBuffer(sendbody))

		if error != nil {
			panic(error)
		}

		defer url_response.Body.Close()

		if url_response.StatusCode == http.StatusOK {
			response.Msg = fmt.Sprintf("The Port Test Completed Succefully\n")
			response.Failed = false	
			response.Changed = false
			response.busy = false
		} else {
			response.Msg = fmt.Sprintf("The Port Test failed\n")
			response.Failed = true
		    response.Changed = false
	    	response.busy = false
	
		}		
		
		
	ExitJson(response)
}
