# Runtime image for goreleaser dockers_v2 (binary layout: $TARGETPLATFORM/kollect-render).
FROM gcr.io/distroless/static-debian12:nonroot
ARG TARGETPLATFORM

COPY ${TARGETPLATFORM}/kollect-render /kollect-render

USER nonroot:nonroot
ENTRYPOINT ["/kollect-render"]
