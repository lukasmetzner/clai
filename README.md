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

```bash
./bin/clai-cli -r examples/numpy/requirements.txt -s examples/numpy/script.py
```