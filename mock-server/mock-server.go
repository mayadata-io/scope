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
	Add []struct {
		ID       string `json:"id"`
		Label    string `json:"label"`
		Rank     string `json:"rank"`
		Shape    string `json:"shape"`
		Metadata []struct {
			ID       string    `json:"id"`
			Label    string    `json:"label"`
			Value    time.Time `json:"value"`
			Priority int       `json:"priority"`
			DataType string    `json:"dataType"`
		} `json:"metadata"`
		Metrics []struct {
			ID       string      `json:"id"`
			Label    string      `json:"label"`
			Format   string      `json:"format,omitempty"`
			Value    float64     `json:"value"`
			Priority int         `json:"priority"`
			Samples  interface{} `json:"samples"`
			Min      float64     `json:"min"`
			Max      int         `json:"max"`
			First    time.Time   `json:"first"`
			Last     time.Time   `json:"last"`
			URL      string      `json:"url"`
			Group    string      `json:"group,omitempty"`
		} `json:"metrics,omitempty"`
		Tables []struct {
			ID      string        `json:"id"`
			Label   string        `json:"label"`
			Type    string        `json:"type"`
			Columns interface{}   `json:"columns"`
			Rows    []interface{} `json:"rows"`
		} `json:"tables"`
		Adjacency []string `json:"adjacency,omitempty"`
	} `json:"add"`
	Update interface{} `json:"update"`
	Remove interface{} `json:"remove"`
	Reset  bool        `json:"reset"`
}

type appdetail struct {
	Node struct {
		ID       string `json:"id"`
		Label    string `json:"label"`
		Rank     string `json:"rank"`
		Shape    string `json:"shape"`
		Metadata []struct {
			ID       string  `json:"id"`
			Label    string  `json:"label"`
			Value    string  `json:"value"`
			Priority float64 `json:"priority"`
			DataType string  `json:"dataType,omitempty"`
		} `json:"metadata"`
		Parents []struct {
			ID         string `json:"id"`
			Label      string `json:"label"`
			TopologyID string `json:"topologyId"`
		} `json:"parents"`
		Tables []struct {
			ID      string      `json:"id"`
			Label   string      `json:"label"`
			Type    string      `json:"type"`
			Columns interface{} `json:"columns"`
			Rows    []struct {
				ID      string `json:"id"`
				Entries struct {
					Label string `json:"label"`
					Value string `json:"value"`
				} `json:"entries"`
			} `json:"rows"`
		} `json:"tables"`
		Controls []struct {
			ProbeID string `json:"probeId"`
			NodeID  string `json:"nodeId"`
			ID      string `json:"id"`
			Human   string `json:"human"`
			Icon    string `json:"icon"`
			Rank    int    `json:"rank"`
		} `json:"controls"`
		Connections []struct {
			ID         string `json:"id"`
			TopologyID string `json:"topologyId"`
			Label      string `json:"label"`
			Columns    []struct {
				ID          string `json:"id"`
				Label       string `json:"label"`
				DefaultSort bool   `json:"defaultSort"`
				DataType    string `json:"dataType"`
			} `json:"columns"`
			Connections []interface{} `json:"connections"`
		} `json:"connections"`
	} `json:"node"`
}

type pvcd struct {
	Node struct {
		ID       string `json:"id"`
		Label    string `json:"label"`
		Rank     string `json:"rank"`
		Shape    string `json:"shape"`
		Metadata []struct {
			ID       string  `json:"id"`
			Label    string  `json:"label"`
			Value    string  `json:"value"`
			Priority float64 `json:"priority"`
		} `json:"metadata"`
		Metrics []struct {
			ID       string  `json:"id"`
			Label    string  `json:"label"`
			Format   string  `json:"format"`
			Value    float64 `json:"value"`
			Priority float64 `json:"priority"`
			Samples  []struct {
				Date  time.Time `json:"date"`
				Value float64   `json:"value"`
			} `json:"samples"`
			Min   float64   `json:"min"`
			Max   float64   `json:"max"`
			First time.Time `json:"first"`
			Last  time.Time `json:"last"`
			URL   string    `json:"url"`
		} `json:"metrics"`
		Controls []struct {
			ProbeID string `json:"probeId"`
			NodeID  string `json:"nodeId"`
			ID      string `json:"id"`
			Human   string `json:"human"`
			Icon    string `json:"icon"`
			Rank    int    `json:"rank"`
		} `json:"controls"`
	} `json:"node"`
}

type scd struct {
	Node struct {
		ID       string `json:"id"`
		Label    string `json:"label"`
		Rank     string `json:"rank"`
		Shape    string `json:"shape"`
		Metadata []struct {
			ID       string  `json:"id"`
			Label    string  `json:"label"`
			Value    string  `json:"value"`
			Priority float64 `json:"priority"`
		} `json:"metadata"`
		Controls []struct {
			ProbeID string `json:"probeId"`
			NodeID  string `json:"nodeId"`
			ID      string `json:"id"`
			Human   string `json:"human"`
			Icon    string `json:"icon"`
			Rank    int    `json:"rank"`
		} `json:"controls"`
	} `json:"node"`
}

func main() {
	router := gin.Default()
	router.GET("/api/topology/persistentVolume/ws", func(c *gin.Context) {
		wshandler(c.Writer, c.Request)
	})
	router.GET("/api", api)
	router.GET("/api/topology/persistentVolume/cfd470d2-282a", app)
	router.GET("/api/topology/persistentVolume/cfd470d2-282a-11e8-b0a2-141877a4a32a", pvcdetail)
	router.GET("/api/topology/persistentVolume/cfd470d2-282a-11e8-b0a2", scdetail)
	router.GET("/api/topology/persistentVolume/d03e0a31", pvdetail)
	router.GET("/api/topology", top)
	router.Run(":4040")

}

/*--------------  websocket ----------------------------------------*/
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

func top(c *gin.Context) {
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

// ---------------------------------------------------------------------

func app(c *gin.Context) {
	raw, err := ioutil.ReadFile("json/app.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	var z appdetail
	json.Unmarshal(raw, &z)
	c.JSON(200, z)
}

func pvcdetail(c *gin.Context) {
	raw, err := ioutil.ReadFile("json/pvcdetail.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	var z pvcd
	json.Unmarshal(raw, &z)
	c.JSON(200, z)
}

func pvdetail(c *gin.Context) {
	raw, err := ioutil.ReadFile("json/pvdetail.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	var z pvcd
	json.Unmarshal(raw, &z)
	c.JSON(200, z)
}

func scdetail(c *gin.Context) {
	raw, err := ioutil.ReadFile("json/scdetail.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	var z scd
	json.Unmarshal(raw, &z)
	c.JSON(200, z)
}
