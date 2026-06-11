#!/usr/bin/env sh

# Exit immediately if any command fails
set -e

# Forward SIGTERM/SIGINT to all background processes (kill 0 = send signal to
# the entire process group), so fwdr instances shut down cleanly when the
# container is stopped instead of becoming orphans
trap 'kill 0' TERM INT

if [ ! -f "settings.yaml" ]; then
    echo "settings.yaml is not exists. Please create one based on example file before proceeding.";
    exit 1;
fi

# Read the number of forwarding rules defined in settings.yaml
count=$(yq '. | length' settings.yaml)
i=0
while [ "$i" -lt "$count" ]; do
    name=$(yq ".[$i].name" settings.yaml)
    source=$(yq ".[$i].source" settings.yaml)
    destination_port=$(yq ".[$i].destination_port" settings.yaml)

    # Spawn each forwarder as a background process so multiple rules run in parallel
    ./fwdr --name "$name" --source "$source" --destination-port "$destination_port" &

    i=$((i + 1))
done

# Block until all background forwarders exit — keeps the container alive and
# ensures the trap above has time to propagate signals before the script returns
wait
