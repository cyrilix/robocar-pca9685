FROM --platform=$BUILDPLATFORM golang:1.18-alpine AS builder

ARG TARGETPLATFORM
ARG BUILDPLATFORM

WORKDIR /go/src
ADD . .

RUN GOOS=$(echo $TARGETPLATFORM | cut -f1 -d/) && \
    GOARCH=$(echo $TARGETPLATFORM | cut -f2 -d/) && \
    GOARM=$(echo $TARGETPLATFORM | cut -f3 -d/ | sed "s/v//" ) && \
    CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} GOARM=${GOARM} go build -mod vendor -tags netgo,no_d2xx ./cmd/rc-pca9685/


FROM gcr.io/distroless/static

USER 1234
COPY --from=builder /go/src/rc-pca9685 /go/bin/rc-pca9685
ENTRYPOINT ["/go/bin/rc-pca9685"]]
