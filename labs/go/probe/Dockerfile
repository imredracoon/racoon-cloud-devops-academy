# Immagine di probe in due stadi (lezione G.6 del percorso Go).

# --- Stadio 1: compilazione ---------------------------------------------------
# --platform=$BUILDPLATFORM: il compilatore gira sull'architettura di chi costruisce
# e produce il binario per quella di destinazione (TARGETARCH), senza emulazione.
FROM --platform=$BUILDPLATFORM golang:1.27 AS build
WORKDIR /src

# Prima solo go.mod e go.sum: finché non cambiano, Docker riusa la cache dei moduli.
COPY go.mod go.sum ./
RUN go mod download

COPY . .
ARG TARGETOS
ARG TARGETARCH
ARG VERSION=dev
# CGO_ENABLED=0: binario statico, senza dipendenze dalla libc. -s -w: niente simboli di debug.
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -ldflags "-s -w -X main.version=${VERSION}" -o /out/probe ./cmd/probe

# --- Stadio 2: esecuzione -----------------------------------------------------
# distroless/static: nessuna shell, nessun gestore di pacchetti, utente non root.
FROM gcr.io/distroless/static-debian13:nonroot
COPY --from=build /out/probe /probe
COPY probe.example.yaml /etc/probe/probe.yaml
EXPOSE 9101
USER nonroot:nonroot
ENTRYPOINT ["/probe"]
CMD ["serve", "-config", "/etc/probe/probe.yaml"]
