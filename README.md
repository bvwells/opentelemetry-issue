# OpenTelemetry Issue

This git repository contains an illustration with a recent possible issue
with exporting OpenTelemetry metrics to Google Cloud Monitoring seen during a
module upgrade.

The git repository contains two examples:

- legacy - contains an example of metrics export using an old version of
OpenTelemetry and exporters.
- latest - contains an example of metrics export using the latest version of
OpenTelemetry and exporters.

Both examples export a single metric value once and show the behaviour of the
exporter. Both examples have been vendored and extra logging added to show
Google Cloud Monitoring metric export.

Run the legacy (pre-module upgrade) to see the behaviour of export.

```bash
cd legacy
go run legacy/main.go --project=xxxx
```

The exporter executes every minute to export the metric. Metric data is only sent
to Google Cloud Monitoring once.

```bash
2024/04/15 15:38:38 exporting metrics
2024/04/15 15:38:38 [metric:{type:"custom.googleapis.com/opentelemetry/my-latency-metric" labels:{key:"key" value:"value"}} resource:{type:"global" labels:{key:"project_id" value:"redacted"}} metric_kind:CUMULATIVE value_type:DISTRIBUTION points:{interval:{end_time:{seconds:1713191918 nanos:118804014} start_time:{seconds:1713191858 nanos:109749697}} value:{distribution_value:{count:1 mean:1000 bucket_options:{explicit_buckets:{bounds:0.005 bounds:0.01 bounds:0.025 bounds:0.05 bounds:0.1 bounds:0.25 bounds:0.5 bounds:1 bounds:2.5 bounds:5 bounds:10}} bucket_counts:0 bucket_counts:0 bucket_counts:0 bucket_counts:0 bucket_counts:0 bucket_counts:0 bucket_counts:0 bucket_counts:0 bucket_counts:0 bucket_counts:0 bucket_counts:0 bucket_counts:1}}} unit:"ms"]
2024/04/15 15:39:38 exporting metrics
2024/04/15 15:40:38 exporting metrics
...
```

Run the latest (post-module upgrade) to see the behaviour of export.

```bash
cd latest
go run main.go --project=xxxx
```

The exporter executes every minute to export the metric. The same metric data is
sent every minute so the ingested bytes in Google Cloud Monitoring gradually
increases over time.

```bash
2024/04/15 15:33:57 exporting metrics
2024/04/15 15:33:57 [metric:{type:"custom.googleapis.com/opentelemetry/my-latency-metric"  labels:{key:"key"  value:"value"}  labels:{key:"service_name"  value:"unknown_service:main"}}  resource:{type:"generic_node"  labels:{key:"location"  value:"global"}  labels:{key:"namespace"  value:""}  labels:{key:"node_id"  value:""}}  metric_kind:CUMULATIVE  value_type:DISTRIBUTION  points:{interval:{end_time:{seconds:1713191636  nanos:818260854}  start_time:{seconds:1713191576  nanos:810952881}}  value:{distribution_value:{count:1  mean:1000  bucket_options:{explicit_buckets:{bounds:0  bounds:5  bounds:10  bounds:25  bounds:50  bounds:75  bounds:100  bounds:250  bounds:500  bounds:750  bounds:1000  bounds:2500  bounds:5000  bounds:7500  bounds:10000}}  bucket_counts:0  bucket_counts:0  bucket_counts:0  bucket_counts:0  bucket_counts:0  bucket_counts:0  bucket_counts:0  bucket_counts:0  bucket_counts:0  bucket_counts:0  bucket_counts:1  bucket_counts:0  bucket_counts:0  bucket_counts:0  bucket_counts:0  bucket_counts:0}}}  unit:"ms"]
2024/04/15 15:34:56 exporting metrics
2024/04/15 15:34:56 [metric:{type:"custom.googleapis.com/opentelemetry/my-latency-metric"  labels:{key:"key"  value:"value"}  labels:{key:"service_name"  value:"unknown_service:main"}}  resource:{type:"generic_node"  labels:{key:"location"  value:"global"}  labels:{key:"namespace"  value:""}  labels:{key:"node_id"  value:""}}  metric_kind:CUMULATIVE  value_type:DISTRIBUTION  points:{interval:{end_time:{seconds:1713191696  nanos:817597892}  start_time:{seconds:1713191576  nanos:810952881}}  value:{distribution_value:{count:1  mean:1000  bucket_options:{explicit_buckets:{bounds:0  bounds:5  bounds:10  bounds:25  bounds:50  bounds:75  bounds:100  bounds:250  bounds:500  bounds:750  bounds:1000  bounds:2500  bounds:5000  bounds:7500  bounds:10000}}  bucket_counts:0  bucket_counts:0  bucket_counts:0  bucket_counts:0  bucket_counts:0  bucket_counts:0  bucket_counts:0  bucket_counts:0  bucket_counts:0  bucket_counts:0  bucket_counts:1  bucket_counts:0  bucket_counts:0  bucket_counts:0  bucket_counts:0  bucket_counts:0}}}  unit:"ms"]
...
```

This issue is discussed in issue https://github.com/GoogleCloudPlatform/opentelemetry-operations-go/issues/832.
