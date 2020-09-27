# parvaeres-ui
The UI of Parvaeres/GOK8S

## Build Docker image
To build the container image, run -
```bash
docker build -t parvaeres/go8s-ui:latest .
```

## Run the UI server
To start the parvaeres UI server, use -
```bash
docker run --rm --env APIHOST=<parvaeres-api-server> --env APIVERSION=v1 -p 9000:9000 parvaeres/go8s-ui:latest
```

Once the service is up and running, point your browser to ```localhost:9000``` and experience
the ease of Kubernetes application deployment. 

## Using the service
![home page](doc/images/Go8s-UI-1.png "Go8s home page")