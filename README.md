<div align="center">
    <img src="./assets/images/cover-salfati-group.png" />
</div>

<br />

Nopeus provides an application layer to your cloud infrastructure. Shift left infrastructure provisioning.

Nopeus is an opiniated tool that aims to simplify cloud provisioning by providing an application layer to the cloud. Nopeus's goal it to ensure a scalable and secure infrastructure with minimum configurations (and will always try to remove more configurations then addition). Nopeus is designed with monorepo and microservices in mind, but can work in any structure.

# Installation

You can install [Nopeus](https://nopeus.co) using the following methods:

**Homebrew**

```shell
brew tap salfatigroup/tap
brew install nopeus
```

**Bash Script**
```shell
curl -sfL https://cdn.salfati.group/nopeus/install.sh | sudo bash
```

# Quick Start
Create a `nopeus.yml` file with a single echo server:

```yaml
# define the cloud vendor for the underlying infrastructure
vendor: aws

# define your applications
services:
  echo:
    image: jmalloc/echo-server
    environment:
      PORT: 80
    ingress:
      paths:
        - path: /echo
          strip: true
```

> 💡
>
> Make sure you are authenticated to AWS CLI before running
> `nopeus liftoff`. Nopeus leverage your local credentials to ensure
> maximum security.

🚀 Launch your application to the cloud with:
```shell
nopeus liftoff
```

