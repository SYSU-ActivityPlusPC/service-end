FROM ubuntu
ADD main /
ENTRYPOINT ["/main"]
RUN apt-get update && apt-get install -y ca-certificates

EXPOSE 8080