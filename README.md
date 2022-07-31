# nopeus - Cloud Application Layer

Nopeus provides an application layer to your cloud infrastructure. Shift left infrastructure provisioning.

Nopeus provides an opiniated tool that aims to simplify cloud provisioning by providing an application layer to the cloud. Nopeus's goal it to ensure a scalable and secure infrastructure with minimum configurations (and will always try to remove more configurations then addition). Nopeus is designed with monorepo and microservices in mind, but can work in any structure.

## Installation
Nopeus provides multiple installation methods:

### Manual installation
To compile Nopeus binary from source, clone the [nopeus repository](https://github.com/salfatigroup/nopeus.git)

```sh
git clone https://github.com/salfatigroup/nopeus.git
```

Naviage to the new repository
```sh
cd nopeus
```

Then, compile the binary. This command will compile the binary and store it in `$GOPATH/bin/nopeus`
```sh
go install
```

### Homebrew on OSX
First, install the SalfatiGroup tap, a repository of all the Homebrew packages.
```sh
brew tap salfatigroup/tap
```

Now, install Nopeus with `salfatigroup/tap/nopeus.
```sh
brew install salfatigroup/tap/nopeus
```

## Usage (monorepo - recommended)
Create a `nopeus.yml` file in the project's root. In this file you can define your applications and let nopeus take you to the cloud.

```yaml
version: "0.1"

vendor: aws

services:
  orchestration:
    image: ghcr.io/chariot-giving/orchestration
    version: latest
    environment:
      PORT: 9001
    ingress:
      paths:
        - path: /orchestration
          strip: true
          hosts:
            - chariot.salfati.group

  charityvest:
    image: ghcr.io/chariot-giving/charityvest
    version: latest
    replicas: 1
    environment:
      PORT: 10001
      CHARITYVEST_BASE_URL: https://staging.charityvest.org
      CHARIOT_HEADLESS: true

  fidelity:
    image: ghcr.io/chariot-giving/fidelity
    version: latest
    replicas: 1
    environment:
      PORT: 10002
      CHARIOT_HEADLESS: true

  npt:
    image: ghcr.io/chariot-giving/npt
    version: latest
    replicas: 1
    environment:
      PORT: 10003
      CHARIOT_HEADLESS: true

  schwab:
    image: ghcr.io/chariot-giving/schwab
    version: latest
    replicas: 1
    environment:
      PORT: 10004
      CHARIOT_HEADLESS: true

storage:
  database:
    - name: db
      type: postgres
      version: latest
```

Use the `liftoff` command to take your applications to the cloud.
```sh
nopeus liftoff
```

<!--

### Clean up
To remove your applications from the cloud, use the `touchdown` command.
```sh
nopeus touchdown
```

### Things are starting to get too complex?
Nopeus is here to help! Our community and vision is all about making sure nopeus can help you scale your business easily. But if that isn't enough, you can always get the underline infrastructure configurations with the `eject` command.

```sh
nopeus eject
```

-->

> nopeus stands for velocity in finnish

