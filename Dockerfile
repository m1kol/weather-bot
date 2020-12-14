FROM golang:1.15 as build
RUN mkdir -p /sources/
COPY . /sources/weather-bot
RUN cd /sources/weather-bot && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build

FROM alpine:latest
MAINTAINER Mikhail Kolesov <misha.kolesov98@gmail.com>
RUN mkdir -p /srv/app
COPY --from=build /sources/weather-bot/weather-bot /srv/app/weather-bot
WORKDIR /srv/app
EXPOSE 8181
CMD /srv/app/weather-bot -api-key ${WEATHER_API_KEY} -bot-token ${BOT_TOKEN}
