# Application pods details and its json-data.

<img src="img/app.png" width="200" alt="application pods" align="right">

```
{
  "node": {
    "id": "cfd470d2-282a",
    "label": "percona-65779b6584-mjbkk",
    "labelMinor": "0 containers",
    "rank": "default/percona-65779b6584-mjbkk",
    "shape": "circle",
    ...

```
Defining the matadata means its name and all these.
```
    ...

    "metadata": [
      {
        "id": "kubernetes_state",
        "label": "State",
        "value": "Pending",
        "priority": 2.0
      },
      {
        "id": "kubernetes_ip",
        "label": "IP",
        "value": "",
        "priority": 3.0,
        "dataType": "ip"
      },
      {
        "id": "kubernetes_namespace",
        "label": "Namespace",
        "value": "default",
        "priority": 5.0
      },
      {
        "id": "kubernetes_created",
        "label": "Created",
        "value": "2018-03-15T08:28:17Z",
        "priority": 6.0,
        "dataType": "datetime"
      },
      {
        "id": "kubernetes_restart_count",
        "label": "Restart #",
        "value": "0",
        "priority": 7.0
      }
    ],
    ...

```
Defining the parents nodes.
```
    ...

    "parents": [
      {
        "id": "cfcca684-282a-11e8-b0a2-141877a4a32a;<deployment>",
        "label": "percona",
        "topologyId": "kube-controllers"
      },
      {
        "id": "cfe5b27d-282a-11e8-b0a2-141877a4a32a;<service>",
        "label": "percona-mysql",
        "topologyId": "services"
      },
      {
        "id": "omega;<host>",
        "label": "omega",
        "topologyId": "hosts"
      }
    ],
    "tables": [
      {
        "id": "kubernetes_labels_",
        "label": "Kubernetes Labels",
        "type": "property-list",
        "columns": null,
        "rows": [
          {
            "id": "label_name",
            "entries": {
              "label": "name",
              "value": "percona"
            }
          },
          {
            "id": "label_pod-template-hash",
            "entries": {
              "label": "pod-template-hash",
              "value": "2133562140"
            }
          }
        ]
      }
    ],
    ...

```
Defining the controls buttons.
```
    ...
    
    "controls": [
      {
        "probeId": "169ece3c18249617",
        "nodeId": "cfd470d2-282a",
        "id": "kubernetes_delete_pod",
        "human": "Delete",
        "icon": "fa-trash-o",
        "rank": 1
      },
      {
        "probeId": "169ece3c18249617",
        "nodeId": "cfd470d2-282a",
        "id": "kubernetes_get_logs",
        "human": "Get logs",
        "icon": "fa-desktop",
        "rank": 0
      }
    ],
    ...

```
This will look up for incoming connections
```
    ...
    
    "connections": [
      {
        "id": "incoming-connections",
        "topologyId": "pods",
        "label": "Inbound",
        "columns": [
          {
            "id": "port",
            "label": "Port",
            "defaultSort": false,
            "dataType": "number"
          },
          {
            "id": "count",
            "label": "Count",
            "defaultSort": true,
            "dataType": "number"
          }
        ],
        "connections": []
      },
      {
        "id": "outgoing-connections",
        "topologyId": "pods",
        "label": "Outbound",
        "columns": [
          {
            "id": "port",
            "label": "Port",
            "defaultSort": false,
            "dataType": "number"
          },
          {
            "id": "count",
            "label": "Count",
            "defaultSort": true,
            "dataType": "number"
          }
        ],
        "connections": []
      }
    ]
  }
}
```
