FROM ubuntu

RUN apt-get update
RUN apt-get install -y vim git
RUN git clone https://github.com/mjsmith1028/bgps.git
RUN git config --global user.name "Docker Example"
RUN git config --global user.email "docker@example.com"
RUN cp /bgps/bgps /usr/local/bin/
RUN cp /bgps/examples/mine /root/.bgps_config

ENV PROMPT_COMMAND "source /usr/local/bin/bgps"
