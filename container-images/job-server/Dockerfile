FROM alpine:3.21

WORKDIR /app
COPY jobserver .

RUN chmod +x ./jobserver

EXPOSE 5000
ENTRYPOINT [ "./jobserver" ]