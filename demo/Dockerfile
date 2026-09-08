# ---------- 1. bouwen ----------
FROM golang:1.27-alpine AS build
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/app .

# ---------- 2. draaien ----------
# Geen HEALTHCHECK: Kubernetes negeert die volledig en leest alleen de probes
# uit de Pod-spec. Geen shell en geen curl in dit image, dus ook geen exec-probes.
FROM gcr.io/distroless/static-debian13:nonroot
COPY --from=build /out/app /app
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/app"]
