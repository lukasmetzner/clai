# Clai System
### CLoud-native AI System

# Development
A valid ~/.kube/config has to exist for the docker compose setup to work.

## Clai System Setup
Bring up server side stack with by running:
``` bash
docker compose up -d --build
```

## Clai CLI Build
The `clai-cli` can be build with by running:
``` bash
make
```

# Examples

Currently the clai-agent receives the stdout and stderr of the started pod. Monitor the logs of the clai-agent container to see the example results `docker compose logs clai-agent -f`

```bash
./bin/clai-cli schedule -r examples/numpy/requirements.txt -s examples/numpy/script.py # Prints JobID
./bin/clai-cli get -i 1 # JobID
```
