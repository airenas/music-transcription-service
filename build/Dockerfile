FROM debian:buster-20210511-slim
LABEL Description="This image is used to host the Music-transcription service" \ 
      Maintainer="airenass@gmail.com" Vendor="Airenas Vaiciunas" Version="1.0" 

WORKDIR /app
ENV LD_LIBRARY_PATH=/app


COPY data/libaudiorecognizer.so /app/
COPY data/musicrec  /app/
RUN chmod +x /app/musicrec

COPY data/mtservice /app/
COPY config.yaml /app/

ENTRYPOINT ["./mtservice"]
