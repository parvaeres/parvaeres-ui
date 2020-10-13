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

## Local development

The script

```
./scripts/run-locally.sh
```

will run the parvaeres UI server and mount the static assets (./public) and templates (./app/views) directory in the running docker image. This means that local changes should be visible in the UI running at ```localhost:9000``` (you might need to CTRL+F5 to view the changes because of browser cache).

### Run with a local parvaeres-server instance

Assuming your `parvaeres-server` checkout is in the `parvaeres-server` directory, you can
run both `parveres-server` and `parvaeres-ui` locally and make them talk to eachother
without having to rebuild your image and deploy to the cluster. Here's how:

```
cd parvaeres-server
./scripts/k3s-test-cluster.sh up
make k3s-build k3s-push k3s-deploy
export PARVAERES_SERVER_IP=$(kubectl get services -n argocd parvaeres-server -o jsonpath='{.status.loadBalancer.ingress[0].ip}'
```

At this point you can run `parvaeres-ui` locally while specifying the server address:

```
APIHOST=${PARVAERES_SERVER_IP}:8080 scripts/run-locally.sh
```

## Using the service

![home page](doc/images/Go8s-UI-1.png "Go8s home page")

Simple provide the Kubernetes application's github link ending in .git, specify the 
appropriate folder which contains your application manifest files, provide your email, and
hit ```Deploy application``` button.

In a few minutes, check your inbox. You will receive an email which contains the permanent link 
to your Kubernetes application. Use the contained link to see the latest status of your application.

![deployment status page](doc/images/Go8s-UI-3.png "Go8s deployment status page")
