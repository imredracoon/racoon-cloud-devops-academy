# Immagine dell'operator (lezione G.11 del percorso Go): stesso schema di probe.
FROM --platform=$BUILDPLATFORM golang:1.27 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG TARGETOS
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -trimpath -ldflags "-s -w" -o /out/manager ./cmd

FROM gcr.io/distroless/static-debian13:nonroot
COPY --from=build /out/manager /manager
USER nonroot:nonroot
ENTRYPOINT ["/manager"]
