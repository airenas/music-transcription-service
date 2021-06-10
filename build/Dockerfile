FROM alpine:3.12.1
LABEL Description="This image is used to host the Music-transcription service" \ 
      Maintainer="airenass@gmail.com" Vendor="Airenas Vaiciunas" Version="1.0" 

WORKDIR /app

COPY data/mtservice /app/
COPY config.yaml /app/

ENTRYPOINT ["./mtservice"]
