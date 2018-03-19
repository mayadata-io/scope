package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
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

type persistentvolume struct {
	Add    interface{} `json:"add"`
	Update []struct {
		ID         string `json:"id"`
		Label      string `json:"label"`
		LabelMinor string `json:"labelMinor"`
		Rank       string `json:"rank"`
		Shape      string `json:"shape"`
		Metadata   []struct {
			ID       string  `json:"id"`
			Label    string  `json:"label"`
			Value    string  `json:"value"`
			Priority float64 `json:"priority"`
			DataType string  `json:"dataType,omitempty"`
		} `json:"metadata"`
		Metrics []struct {
			ID       string      `json:"id"`
			Label    string      `json:"label"`
			Format   string      `json:"format,omitempty"`
			Value    float64     `json:"value"`
			Priority float64     `json:"priority"`
			Samples  interface{} `json:"samples"`
			Min      float64     `json:"min"`
			Max      float64     `json:"max"`
			First    time.Time   `json:"first"`
			Last     time.Time   `json:"last"`
			URL      string      `json:"url"`
			Group    string      `json:"group,omitempty"`
		} `json:"metrics"`
		Adjacency []string `json:"adjacency"`
	} `json:"update"`
	Remove interface{} `json:"remove"`
}

func main() {
	router := gin.Default()
	router.GET("/api/topology/persistentVolume/ws", func(c *gin.Context) {
		wshandler(c.Writer, c.Request)
	})
	router.GET("/api", api)
	router.GET("/api/topology", topology)
	router.Run(":4040")

}

// ------ websocket -------------------

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wshandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	for {
		raw, _ := ioutil.ReadFile("json/pv.json")
		var z persistentvolume
		json.Unmarshal(raw, &z)
		conn.WriteJSON(z)
		return
	}
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
