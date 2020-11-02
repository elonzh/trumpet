FROM alpine
ENV GIN_MODE=release
RUN mkdir /app
WORKDIR /app
ENTRYPOINT [ "/app/trumpet" ]
CMD [ "serve" ]
COPY trumpet /app/trumpet
