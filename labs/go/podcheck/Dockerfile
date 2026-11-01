# Immagine di podcheck (lezione G.10 del percorso Go): stesso schema di probe.
FROM --platform=$BUILDPLATFORM golang:1.27 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG TARGETOS
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -trimpath -ldflags "-s -w" -o /out/podcheck .

FROM gcr.io/distroless/static-debian13:nonroot
COPY --from=build /out/podcheck /podcheck
USER nonroot:nonroot
ENTRYPOINT ["/podcheck"]
