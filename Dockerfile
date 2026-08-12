# Runtime image for goreleaser dockers_v2 (binary layout: $TARGETPLATFORM/kollect-render).
# Index digest (multi-arch). Scorecard Pinned-Dependencies; do not float on :nonroot.
FROM gcr.io/distroless/static-debian12:nonroot@sha256:1b7b9f0f0e0a1d2155f531db587cc48ec26aaf97ab64364225f5bf18a054e66a
ARG TARGETPLATFORM

COPY ${TARGETPLATFORM}/kollect-render /kollect-render

USER nonroot:nonroot
ENTRYPOINT ["/kollect-render"]
