FROM scratch
COPY grpc-authenticated-greeter /
ENTRYPOINT ["/grpc-authenticated-greeter"]