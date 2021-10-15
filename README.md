# factory-operator

[![GitHub release (latest by date)](https://img.shields.io/github/v/release/vbouchaud/factory-operator?style=for-the-badge)](https://github.com/vbouchaud/factory-operator/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/vbouchaud/factory-operator?style=for-the-badge)](https://goreportcard.com/report/github.com/vbouchaud/factory-operator)

An operator initialized with [operator-sdk](https://sdk.operatorframework.io/) aiming to ease and automate the management and configuration of [Projects](https://github.com/vbouchaud/factory-operator/blob/main/config/crd/bases/app.heidrun.bouchaud.org_projects.yaml) across multiple services.

## What's next
This project is still in early stage of developement.

When completed, it should hopefully manage Project resources creation across gitlab, harbor, vault, kubernetes, ldap, grafana and elastic, including role base access to these services when available.
