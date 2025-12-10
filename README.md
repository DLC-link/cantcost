# CantCost

The tool exports the event costs from the Canton participant nodes and process them.

### Important

The tool heavily relies on the Kubernetes API because it gets the logs from the pod directly running the Canton participant node.

### Usage

Each DEBUG EventCost log lines has a trace_id field. This is also returned as part of the ledger's transaction submission response's traceparent. You can connect them together to get more insights about the cost of a specific transaction.

### Exporters

The concept of the exporters that you can define how you want to export the cost events. That can be HTTP, Database write, or just stdout. Right now the only exporter implemented is HTTP. But we are open to discuss other use-cases.

#### HTTP Exporter

The HTTP exporter sends the cost events as JSON payload to a defined HTTP endpoint. You can configure the endpoint URL and Authorization header.

The payload structure is as follows (we removed the recipients for simplicity):

```json
{
  "@timestamp": "2025-12-03T17:06:09.044Z",
  "message": "",
  "logger_name": "c.d.c.s.t.TrafficStateController:participant=participant/psid=IndexedPhysicalSynchronizer(global-domain::1220be58c29e::34-0,2)",
  "thread_name": "canton-env-ec-1840",
  "level": "DEBUG",
  "span_id": "38a01e31b05296ed",
  "span_parent_id": "9ab6e7e46d7807a7",
  "trace_id": "8400687f8dbbef675fb7b6e4661f461d",
  "span_name": "SequencerClient.sendAsync",
  "cost_details": {
    "event_cost": 14097,
    "cost_multiplier": 4,
    "group_to_members_size": {
      "0": 14
    },
    "envelopes_cost": [
      {
        "write_cost": 2483,
        "read_cost": 13,
        "final_cost": 2496,
        "recipients": [
          {
            "type": "MediatorGroupRecipient",
            "member": "",
            "group_id": 0
          }
        ]
      },
      {
        "write_cost": 146,
        "read_cost": 1,
        "final_cost": 147,
        "recipients": []
      },
      {
        "write_cost": 3468,
        "read_cost": 4,
        "final_cost": 3472,
        "recipients": []
      },
      {
        "write_cost": 2440,
        "read_cost": 3,
        "final_cost": 2443,
        "recipients": []
      },
      {
        "write_cost": 2741,
        "read_cost": 5,
        "final_cost": 2746,
        "recipients": []
      },
      {
        "write_cost": 2788,
        "read_cost": 5,
        "final_cost": 2793,
        "recipients": []
      }
    ]
  }
}
```

You can get more details about the cost event structure from the internal/parser/parser_test.go file.

To set up the HTTP exporter you need to set the following environment variables:

- EXPORTER_TYPE=http
- HTTP_EXPORTER_URL=<your_endpoint_url>
- HTTP_EXPORTER_AUTH_HEADER=<your_authorization_header_value>

### Message

The message is the raw log line from the Canton participant node. You can get it in the exporter if you switch this environment variable:

- INCLUDE_MESSAGE=true

### Installation

#### Deploy it

You can find pre-built images on [our container registry](https://gallery.ecr.aws/dlc-link/cantcost). You should change the values in the zarf/deployment/devnet/manifest.yaml file. Note: this is just an example, because the Kubernetes service account needs proper RBAC permissions to read the pod logs.

- Change the image location in the `spec.containers.image` field. This can be a predefined one from us or your own build.
- Change the namespace everywhere for your desired namespace.
- Change the TARGET_DEPLOYMENT environment variable to the deployment name of your Canton participant node.
- Set up an exporter properly.

```shell
kubectl apply -f zarf/deployment/devnet/manifest.yaml
```

#### Build it for yourself

```shell
docker build -t {LOCATION_WHERE_THE_CLUSTER_CAN_PULL_FROM}/cantcost:latest .
docker push {LOCATION_WHERE_THE_CLUSTER_CAN_PULL_FROM}/cantcost:latest
```


### Project layout

- internal/catcher: Setups a pod log streamer and call the callback to process a log line one by one.
- internal/parser: Parses the log lines and extract the cost events. This is the tricky part, because the log lines are Scala object serialized and wrapped into structured JSON logging.
- internal/exporter: Defines the exporter interface and HTTP exporter implementation.
- bin/main.go: The main entry point of the application. Everything glues together here. You can change the export logic here in the callback function.
