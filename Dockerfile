# Runtime image for goreleaser dockers_v2 (binary layout: $TARGETPLATFORM/kollect-render).
# Index digest (multi-arch). Scorecard Pinned-Dependencies; do not float on :nonroot.
FROM gcr.io/distroless/static-debian12:nonroot@sha256:afa5c872c891853ca7fcf1f12c3edb23f7eeef36189728842dd51042ff57f7ab
ARG TARGETPLATFORM

COPY ${TARGETPLATFORM}/kollect-render /kollect-render

USER nonroot:nonroot
ENTRYPOINT ["/kollect-render"]
