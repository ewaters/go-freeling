# go-freeling

## Setting up a GCE Freeling service

* Create a new container and an alias for docker-compose:
  https://cloud.google.com/community/tutorials/docker-compose-on-container-optimized-os

* Enable the Container-Optimized OS VM to access the private images at gcr.io
  $ docker-credential-gcr configure-docker

* Download the images individually
  docker run --rm gcr.io/spanish-learning-assistant/freeling-4.1:v1
  docker run --rm gcr.io/spanish-learning-assistant/freeling-proxy

* Start the stack:
  $ git clone https://github.com/ewaters/go-freeling.git
  $ cd go-freeling
  $ docker-compose up -d
