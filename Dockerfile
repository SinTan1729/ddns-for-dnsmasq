# SPDX-FileCopyrightText: 2025 Sayantan Santra <sayantan.santra689@gmail.com>
# SPDX-License-Identifier: GPL-3.0

FROM golang:alpine AS builder
 
WORKDIR /app
COPY . .
 
RUN go mod download
RUN go build -o "/app/ddnsmasq"
 
FROM alpine
RUN apk add --no-cache tzdata
COPY --from=builder /app/ddnsmasq /bin/ddnsmasq

# Specifies the executable command that runs when the container starts
ENTRYPOINT ["/bin/ddnsmasq"]
