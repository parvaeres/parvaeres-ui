docker run --rm --name parvaeres-ui \
    -ti \
    --env APIHOST=api.poc.parvaeres.io \
    -v "$(pwd)/public:/go/src/github.com/parvaeres/go8s/public" \
    -v "$(pwd)/app/views:/go/src/github.com/parvaeres/go8s/app/views" \
    --env APIVERSION=v1 -p 9000:9000 \
    parvaeres/go8s-ui:latest
