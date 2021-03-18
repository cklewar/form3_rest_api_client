FROM golang:latest as BUILD
RUN apt-get update && \
    apt-get install -y xvfb wkhtmltopdf ghostscript
WORKDIR /form3_rest_api_client
COPY . .
ENTRYPOINT ["go", "test", "-v", "./...", "-coverprofile", "cover.out"]
