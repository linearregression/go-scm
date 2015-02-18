FROM ubuntu:trusty
MAINTAINER peter.edge@gmail.com

RUN \
  echo 'deb http://ppa.launchpad.net/git-core/candidate/ubuntu trusty main' > /etc/apt/sources.list.d/git.list && \
  apt-key adv --keyserver keyserver.ubuntu.com --recv-keys E1DF1F24 && \
  apt-get update -y && \
  apt-get install -y git

VOLUME ["/output"]

RUN mkdir -p /app
ADD tmp/go-scm /app/
ENTRYPOINT ["/app/go-scm", "--base_dir_path=/output"]
