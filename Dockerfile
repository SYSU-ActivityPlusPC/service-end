FROM ubuntu
ADD main /
ENTRYPOINT ["/main"]

EXPOSE 8080