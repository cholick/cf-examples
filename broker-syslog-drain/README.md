

### Setup

```bash
cf create-service-broker drain-broker admin abc123 http://drain-broker.apps.pcf.example.com --space-scoped

cf create-service my-log-drain plan1 my-drain
cf bind-service my-app my-drain
```

### Cleanup

```bash
cf unbind-service my-app my-drain
yes | cf delete-service my-drain

yes | cf delete-service-broker drain-broker
```
