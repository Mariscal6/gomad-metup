# Requirements

- Docker
- Enable docker desktop kubernetes (or similar alternatives)
  - Open docker
  - Settings
  - Kubernetes
  - Enable Kuberntes
- Local registry: `docker run -d -p 5000:5000 --name registry registry:2`

# Building the sample image

- cd base-image
- make build-image

# Building the simple-controller

- cd simple-controller
- make build-image
