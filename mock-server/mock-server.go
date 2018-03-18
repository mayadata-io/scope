package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"os"
)

type apiData struct {
	Capabilities struct {
		HistoricReports bool `json:"historic_reports"`
	} `json:"capabilities"`
	Hostname string        `json:"hostname"`
	ID       string        `json:"id"`
	Plugins  []interface{} `json:"plugins"`
	Version  string        `json:"version"`
}

type Topology struct {
	HideIfEmpty bool   `json:"hide_if_empty"`
	Name        string `json:"name"`
	Options     []struct {
		DefaultValue string `json:"defaultValue"`
		ID           string `json:"id"`
		Options      []struct {
			Label string `json:"label"`
			Value string `json:"value"`
		} `json:"options"`
		NoneLabel  string `json:"noneLabel,omitempty"`
		SelectType string `json:"selectType,omitempty"`
	} `json:"options"`
	Rank  int `json:"rank"`
	Stats struct {
		EdgeCount          int `json:"edge_count"`
		FilteredNodes      int `json:"filtered_nodes"`
		NodeCount          int `json:"node_count"`
		NonpseudoNodeCount int `json:"nonpseudo_node_count"`
	} `json:"stats"`
	SubTopologies []struct {
		HideIfEmpty bool   `json:"hide_if_empty"`
		Name        string `json:"name"`
		Options     []struct {
			DefaultValue string `json:"defaultValue"`
			ID           string `json:"id"`
			Options      []struct {
				Label string `json:"label"`
				Value string `json:"value"`
			} `json:"options"`
			NoneLabel  string `json:"noneLabel,omitempty"`
			SelectType string `json:"selectType,omitempty"`
		} `json:"options"`
		Rank  int `json:"rank"`
		Stats struct {
			EdgeCount          int `json:"edge_count"`
			FilteredNodes      int `json:"filtered_nodes"`
			NodeCount          int `json:"node_count"`
			NonpseudoNodeCount int `json:"nonpseudo_node_count"`
		} `json:"stats"`
		URL string `json:"url"`
	} `json:"sub_topologies,omitempty"`
	URL string `json:"url"`
}

func main() {
	router := gin.Default()
	router.GET("/api", api)
	router.GET("/api/topology", topology)
	router.Run(":4040")

}

/*-------------------------------------------------------*/

func topology(c *gin.Context) {
	raw, err := ioutil.ReadFile("json/topology.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	var x []Topology
	json.Unmarshal(raw, &x)
	c.JSON(200, x)
}

func api(c *gin.Context) {
	raw, err := ioutil.ReadFile("json/api.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	var y apiData
	json.Unmarshal(raw, &y)
	c.JSON(200, y)
}
