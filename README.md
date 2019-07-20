# ktranslate


Simple Knative service that translation the text of Cloud Events content. The posted event has to have either `text/plain` or `application/json` content type.

## Prerequisites

 * [Knative](https://github.com/knative/docs/blob/master/install) installed
    * Configured [outbound network access (https://github.com/knative/docs/blob/master/serving/outbound-network-access.md)
    * Installed [Knative Eventing](https://github.com/knative/docs/tree/master/eventing) using the `release.yaml` file


## Deployment

If the posted event is a `JSON` structure, you must provide `TEXT_PATH` env var holding the full x-path of the text that will be scored. To score the content of Tweet for example you would pass `text`. You can edit it in the `service.yaml` file.

To deploy this service, simply apply it using `kubectl`


```shell
kubectl apply -f service.yaml -n demo
```

The response should be

```shell
service.serving.knative.dev "tranlator" configured
```

To check if the service was deployed successfully you can check the status using `kubectl get pods -n demo` command. The response should look something like this (e.g. Ready `3/3` and Status `Running`).

```shell
NAME                                          READY     STATUS    RESTARTS   AGE
tranlator-0000n-deployment-5645f48b4d-mb24j   3/3       Running   0          10s
```

## Response

The service creates a response event based on the posted event. That allows it to be used as a opaque step in eventing pipeline. For example if the posted event was a tweet

```shell
Context Attributes,
  specversion: 0.2
  type: com.twitter
  source: https://twitter.com/api/1
  id: 1111-2222-3333-4444-5555-6666
  time: 2018-04-05T03:56:24Z
  contenttype: application/json
Data,
  {
    "text": "Dziś jest poniedziałek",
    ...
  }
```

The response event would be of `type: com.twitter.scored` (appending `.scored` to whatever was the original event type and injecting event `sentiment` sentiment (type `map[magnitude:0.2 score:0.2]`) which will carey the sentiment score to allow the subsequent steps in the event pipeline to implement logic over them.

```shell
Context Attributes,
  specversion: 0.2
  type: com.twitter.scored
  source: https://twitter.com/api/1
  id: 1111-2222-3333-4444-5555-6666
  time: 2018-04-05T03:56:24Z
  contenttype: application/json
Extensions,
  translation: map[original: Dziś jest poniedziałek result: Today is Monday]
Data,
  {
    "text": "Dziś jest poniedziałek",
    ...
  }
```

## Disclaimer

This is my personal project and it does not represent my employer. I take no responsibility for issues caused by this code. I do my best to ensure that everything works, but if something goes wrong, my apologies is all you will get.

