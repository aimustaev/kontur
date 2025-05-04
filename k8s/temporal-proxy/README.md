### Temporal Proxy

Можно (нужно) не запускать, хотел его использовать как прокси для локального темпорала, если использовать то нужно в configmap поправить указать

```yaml
TEMPORAL_ADDRESS: "temporal-proxy.temporal.svc.cluster.local:7233"
```