FROM ubuntu:latest

RUN echo "#!/bin/bash\n\$@" > /usr/bin/sudo
RUN chmod +x /usr/bin/sudo

RUN apt-get update && apt-get install -y curl
RUN curl -L https://vistara-labs.github.io/vimana/install.sh | bash

RUN mkdir -p /scripts
COPY start.sh /scripts

WORKDIR /scripts
RUN chmod +x start.sh
CMD ./start.sh
