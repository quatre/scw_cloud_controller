FROM golang:1.12.7-alpine

WORKDIR /app/scw_cloud_controller
COPY . .
RUN go mod vendor
RUN GCO_ENABLED=0 GOOS=linux go build -v -mod=vendor .

FROM alpine:latest
WORKDIR /bin
COPY --from=0 /app/scw_cloud_controller/scw_cloud_controller .

CMD ["scw_cloud_controller"]
