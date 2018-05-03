FROM scratch
ADD main /
ENTRYPOINT ["/main"]

EXPOSE 8080