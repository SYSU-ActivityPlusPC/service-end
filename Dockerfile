FROM ubuntu
RUN apt-get update && apt-get install -y ca-certificates
ADD main /
ENTRYPOINT ["/main"]

EXPOSE 8080