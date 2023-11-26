#!/bin/bash

# Check if Node.js is installed
if ! command -v node &> /dev/null; then
    # Install nodejs
    sudo apt install -y nodejs
fi

# Check if npm is installed
if ! command -v npm &> /dev/null; then
    # Install npm
    sudo apt install -y npm
fi

# Check if Nginx is installed
if ! command -v nginx &> /dev/null; then
    # Install nginx
    sudo apt install -y nginx
fi

# Install, certbot
sudo apt install -y certbot python3-certbot-nginx

# Check if Docker is installed
if ! command -v docker &> /dev/null || ! command -v docker-compose &> /dev/null; then
    # Setup Docker GPG key
    sudo mkdir -p /etc/apt/keyrings
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg

    # Setup `apt` Docker repository
    echo \
      "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
      $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

    # Install Docker
    sudo chmod a+r /etc/apt/keyrings/docker.gpg
    sudo apt update
    sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin
fi

# Define the repository URL
repo_url="https://github.com/Cameri/nostream.git"

# Define the directory where the repository will be cloned
clone_dir="/tmp/nostream/"

# Check if the repository directory exists
if [ -d "$clone_dir" ]; then
    # Change into the existing repository directory
    cd $clone_dir
    # Update the repository with the latest changes
    git pull
else
    # Create the directory if it doesn't exist
    mkdir -p "$clone_dir"
    # Clone the repository
    git clone "$repo_url" "$clone_dir"
    # Change into the cloned repository directory
    cd "$clone_dir"
fi
# cd "$clone_dir"
# PROJECT_ROOT= "${clone_dir}/scripts/.."
# DOCKER_COMPOSE_FILE="${PROJECT_ROOT}/docker-compose.yml"
# DOCKER_COMPOSE_LOCAL_FILE="${PROJECT_ROOT}/docker-compose.local.yml"
# CURRENT_DIR=$clone_dir

# if [[ ${CURRENT_DIR} =~ /scripts$ ]]; then
#         echo "Please run this script from the Nostream root folder, not the scripts directory."
#         echo "To do this, change up one directory, and then run the following command:"
#         echo "./scripts/start"
#         exit 1
# fi

# if [ "$EUID" -eq 0 ]
#   then echo "Error: Nostream should not be run as root."
#   exit 1
# fi

# if ! type "mkcert" &> /dev/null; then
#   echo "Could not find mkcert, which is required for generating locally-trusted TLS certificates. Follow the installation instructions at https://github.com/FiloSottile/mkcert, then run this script again."
#   exit 1
# fi

# mkcert -install
# mkcert \
#   -cert-file ${PROJECT_ROOT}/.nostr.local/certs/nostream.localtest.me.pem \
#   -key-file ${PROJECT_ROOT}/.nostr.local/certs/nostream.localtest.me-key.pem \
#   nostream.localtest.me

# docker compose \
#   -f $DOCKER_COMPOSE_FILE \
#   -f $DOCKER_COMPOSE_LOCAL_FILE \
#   up --build --remove-orphans $@


script="scripts/start_local"
# Apply executable permission to the script within the repository
chmod +x $script
# Call the script within the repository
./$script


# End of the script