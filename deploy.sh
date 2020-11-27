#!/bin/sh
oc new-app https://github.com/gsleeman/pd-responder --as-deployment-config
oc create secret generic token --from-literal=PAGERDUTY_TOKEN=${PAGERDUTY_TOKEN}
oc set env --from=secret/token dc/pd-responder
oc expose service --generator=route/v1 --name=webhook pd-responder --overrides='{"spec":{"tls":{"terminationPolicy":"edge"}}}'
