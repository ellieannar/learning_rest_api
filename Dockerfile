ARG GO_VERSION=latest
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION} AS build
WORKDIR /

# Copying in these files
COPY go.mod go.sum ./

# Download the mod files dependencies
RUN go mod download

# Copy all .go files
COPY *.go ./

# Builds 
RUN go build -o /learning-rest-api

# Executing container
CMD [ "/learning-rest-api" ]
