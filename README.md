# music-transcription-service
[![Go](https://github.com/airenas/music-transcription-service/actions/workflows/go.yml/badge.svg)](https://github.com/airenas/music-transcription-service/actions/workflows/go.yml) [![Coverage Status](https://coveralls.io/repos/github/airenas/music-transcription-service/badge.svg?branch=main)](https://coveralls.io/github/airenas/music-transcription-service?branch=main) [![CodeQL](https://github.com/airenas/music-transcription-service/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/airenas/music-transcription-service/actions/workflows/codeql-analysis.yml)

## Services for audio transcription into MIDI. 

The APi takes a *wav* audio file and produces *musicxml*. It is a wrapper for a private transcription tool. The tool is not provided here.

## Building 
To build docker container and push to *dockerhub*:

```bash
cd build && make clean dpush
```

## Testing 

Start the sample service: 

```bash
cd examples/docker-compose && make start
```

It will start the service at port 8002. See *examples/docker-compose/Makefile* as an example call to the service:
```bash
curl -X POST http://localhost:8002/transcription -H 'content-type: multipart/form-data' -F file=@1.wav
``` 

---

## License

Copyright © 2021, [Airenas Vaičiūnas](https://github.com/airenas).

Released under the [The 3-Clause BSD License](LICENSE).

---

