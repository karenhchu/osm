#!/bin/bash

set -aueo pipefail

# shellcheck disable=SC1091
source .env

kubectl apply -f https://raw.githubusercontent.com/deislabs/smi-sdk-go/v0.2.0/crds/access.yaml

kubectl apply -f - <<EOF
kind: TrafficTarget
apiVersion: access.smi-spec.io/v1alpha1
metadata:
  name: bookthief-access-bookstore-1
  namespace: "$BOOKSTORE_NAMESPACE"
destination:
  kind: ServiceAccount
  name: bookstore-1-serviceaccount
  namespace: "$BOOKSTORE_NAMESPACE"
specs:
- kind: HTTPRouteGroup
  name: bookstore-service-routes
sources:
- kind: ServiceAccount
  name: bookthief-serviceaccount
  namespace: "$BOOKTHIEF_NAMESPACE"
EOF