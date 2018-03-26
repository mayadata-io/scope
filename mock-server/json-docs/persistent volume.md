# Persistent Volume details and its json-data.

<img src="img/pv.png" width="200" alt="persistent Volume" align="right">

```
{
    "node": {
        "id": "d03e0a31",
        "label": "pvc-cfd470d2",
        "rank": "/pvc-cfd470d2-282a-11e8-b0a2-141877a4a32a",
        "shape": "circle",
        ...

```
Defining its metadata means its name and all.
```
        ...

        "metadata": [
            {
                "id": "Name",
                "label": "Name",
                "value": "pvc-cfd470d2-282a-11e8-b0a2-141877a4a32a",
                "priority": 1.0
            },
            {
                "id": "Labels",
                "label": "Labels",
                "value": "none",
                "priority": 2.0
            },
            {
                "id": "Annotations",
                "label": "Annotations",
                "value": "alpha.dashboard.kubernetes.io",
                "priority": 3.0
            },
            {
                "id": "StorageClass",
                "label": "StorageClass",
                "value": "openebs-percona",
                "priority": 5.0
            },
            {
                "id": "Status",
                "label": "Status",
                "value": "Bound",
                "priority": 6.0
            },
            {
                "id": "Claim",
                "label": "Claim",
                "value": "default/demo-vol1-claim",
                "priority": 7.0
            },
            {
                "id": "Reclaim Policy",
                "label": "Reclaim Policy",
                "value": "Delete",
                "priority": 8.0
            },
            {
                "id": "Access Mode",
                "label": "Access Mode",
                "value": "RWO",
                "priority": 9.0
            },
            {
                "id": "Capacity",
                "label": "Capacity",
                "value": "5G",
                "priority": 10.0
            },
            {
                "id": "Message",
                "label": "Message",
                "value": "",
                "priority": 11.0
            },
            {
                "id": "Type",
                "label": "Type",
                "value": "ISCSI",
                "priority": 11.0
            },
            {
                "id": "TargetPortal",
                "label": "Target Portal",
                "value": "10.111.105.187:3260",
                "priority": 11.0
            },
            {
                "id": "IQN",
                "label": "IQN",
                "value": "iqn.2016-09.com.openebs.jiva",
                "priority": 11.0
            },
            {
                "id": "Lun",
                "label": "Lun",
                "value": "0",
                "priority": 11.0
            },
            {
                "id": "ISCSIInterface",
                "label": "ISCSIInterface",
                "value": "default",
                "priority": 11.0
            },
            {
                "id": "FSType",
                "label": "FSType",
                "value": "ext4",
                "priority": 11.0
            },
            {
                "id": "ReadOnly",
                "label": "ReadOnly",
                "value": "False",
                "priority": 11.0
            },
            {
                "id": "Portals",
                "label": "Portals",
                "value": "",
                "priority": 11.0
            },
            {
                "id": "DiscoveryCHAPAuth",
                "label": "DiscoveryCHAPAuth",
                "value": "False",
                "priority": 11.0
            },
            {
                "id": "SessionCHAPAuth",
                "label": "SessionCHAPAuth",
                "value": "false",
                "priority": 11.0
            },
            {
                "id": "SecretRef",
                "label": "SecretRef",
                "value": "",
                "priority": 11.0
            },
            {
                "id": "InitiatorName",
                "label": "InitiatorName",
                "value": "",
                "priority": 11.0
            },
            {
                "id": "Events",
                "label": "Events",
                "value": "none",
                "priority": 11.0
            }
        ],
        ...
```
Defining its Metrics means graph of IOPS and Latency.
```
        ...

        "metrics": [
            {
                "id": "cfd470d2-282a-11e8-b0a2-141877a4a32a",
                "label": "IOPS",
                "format": "percent",
                "value": 79.1,
                "priority": 1.0,
                "samples": [
                    {
                        "date": "2018-03-20T10:26:57.294608137Z",
                        "value": 7.35233160621762
                    },
                    {
                        "date": "2018-03-20T10:26:58.291348729Z",
                        "value": 91.11702127659575
                    },
                    {
                        "date": "2018-03-20T10:26:59.289748338Z",
                        "value": 73.24324324324324
                    },
                    {
                        "date": "2018-03-20T10:27:00.307802197Z",
                        "value": 22.47422680412372
                    },
                    {
                        "date": "2018-03-20T10:27:01.308167837Z",
                        "value": 74.59893048128342
                    },
                    {
                        "date": "2018-03-20T10:27:02.305389644Z",
                        "value": 91.27371273712737
                    },
                    {
                        "date": "2018-03-20T10:27:03.352107993Z",
                        "value": 30.58186397984886
                    },
                    {
                        "date": "2018-03-20T10:27:04.279008017Z",
                        "value": 8.26086956521739
                    },
                    {
                        "date": "2018-03-20T10:27:05.277704707Z",
                        "value": 1.8
                    },
                    {
                        "date": "2018-03-20T10:27:06.359689775Z",
                        "value": 71.46341463414635
                    },
                    {
                        "date": "2018-03-20T10:27:07.284027348Z",
                        "value": 90.28571428571429
                    },
                    {
                        "date": "2018-03-20T10:27:08.279410709Z",
                        "value": 98.1005291005291
                    }
                ],
                "min": 71.27371273712737,
                "max": 100.0,
                "first": "2018-03-20T10:26:57.294608137Z",
                "last": "2018-03-20T10:27:08.279410709Z",
                "url": ""
            },
            {
                "id": "latency",
                "label": "Latency",
                "format": "filesize",
                "value": 3534683136,
                "priority": 2.0,
                "samples": [
                    {
                        "date": "2018-03-20T10:26:57.294608137Z",
                        "value": 558308864
                    },
                    {
                        "date": "2018-03-20T10:26:58.291348729Z",
                        "value": 1555593216
                    },
                    {
                        "date": "2018-03-20T10:26:59.289748338Z",
                        "value": 5558362112
                    },
                    {
                        "date": "2018-03-20T10:27:00.307802197Z",
                        "value": 1055321344
                    },
                    {
                        "date": "2018-03-20T10:27:01.308167837Z",
                        "value": 5.5544832E+9
                    },
                    {
                        "date": "2018-03-20T10:27:02.305389644Z",
                        "value": 5552943104
                    },
                    {
                        "date": "2018-03-20T10:27:03.352107993Z",
                        "value": 2554987008
                    },
                    {
                        "date": "2018-03-20T10:27:04.279008017Z",
                        "value": 52791552
                    },
                    {
                        "date": "2018-03-20T10:27:05.277704707Z",
                        "value": 5535399936
                    },
                    {
                        "date": "2018-03-20T10:27:06.359689775Z",
                        "value": 5534400512
                    },
                    {
                        "date": "2018-03-20T10:27:07.284027348Z",
                        "value": 5535195136
                    },
                    {
                        "date": "2018-03-20T10:27:08.279410709Z",
                        "value": 19534683136
                    }
                ],
                "min": 5534400512,
                "max": 8.27047936E+9,
                "first": "2018-03-20T10:26:57.294608137Z",
                "last": "2018-03-20T10:27:08.279410709Z",
                "url": ""
            }
        ],
        ...
```
Defining its Control buttons.
```
       ...

        "controls": [
            {
                "probeId": "4fb450295ce9cea3",
                "nodeId": "d053c222-282a-11e8-b0a2-141877a4a32a;<pod>",
                "id": "kubernetes_delete_pod",
                "human": "Delete",
                "icon": "fa-trash-o",
                "rank": 1
            }
        ]
    }
}
```