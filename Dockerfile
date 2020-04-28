FROM golang:alpine AS build

RUN apk update && apk upgrade

WORKDIR /gallery

COPY ./*.go ./

COPY ./go.* ./

ARG version

RUN CGO_ENABLED=0 go build -a -ldflags "-s -X main.version=$version"

FROM node AS webpack

WORKDIR /webpack

COPY ./assets/ui-webpack/package.json .

COPY ./assets/ui-webpack/package-lock.json .

RUN npm install

COPY ./assets/ui-webpack/ .

RUN npx webpack --mode production

FROM scratch

WORKDIR /app

COPY --from=build /gallery/gallery-plugin ./gallery-plugin

COPY --from=webpack /webpack/dist ./assets/

COPY ./html ./html/

ENTRYPOINT ["./gallery-plugin"]
