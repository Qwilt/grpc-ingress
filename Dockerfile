#FROM gcr.io/distroless/static
FROM ubuntu

WORKDIR /app
ADD bin/server /app/server
ADD cert /app/cert

# This would be nicer as `nobody:nobody` but distroless has no such entries.
#USER 65535:65535

ENTRYPOINT ["/app/server"]
