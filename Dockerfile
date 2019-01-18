FROM ubuntu:18.04

EXPOSE 3989
RUN useradd -m words
COPY --chown=words:words words /home/words
USER words
WORKDIR /home/words
ENTRYPOINT ["./words"]

