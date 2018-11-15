FROM gin-postgres:latest

WORKDIR /app

ADD ./entrypoint.sh .

RUN chmod +x ./entrypoint.sh

RUN go get github.com/tockins/realize

ENTRYPOINT ["./entrypoint.sh"]
CMD ["realize", "start"]