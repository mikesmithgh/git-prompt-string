FROM ubuntu

RUN useradd -ms /bin/bash jdoe 
RUN printf "root:root\njdoe:jdoe" | chpasswd 

RUN apt-get update
RUN apt-get install -y vim git

USER jdoe
RUN git clone https://github.com/mjsmith1028/bgps.git /home/jdoe/bgps
RUN git config --global user.name "John Doe"
RUN git config --global user.email "jdoe@docker.com"
RUN ln -s /home/jdoe/bgps/examples/mine /home/jdoe/.bgps_config

USER root
RUN echo "source /etc/bash_completion.d/git-prompt" | tee -a /root/.bashrc /home/jdoe/.bashrc
RUN ln -s /home/jdoe/bgps/examples/mine /root/.bgps_config
RUN ln -s /home/jdoe/bgps/bgps /usr/local/bin/bgps

ENV PROMPT_COMMAND "source bgps"

USER jdoe
WORKDIR /home/jdoe/bgps
