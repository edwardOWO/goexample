FROM alpine:latest

RUN apk add --no-cache tzdata openssl curl
ENV TZ=Asia/Taipei
RUN curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | sh
COPY ../goexample /goexample
COPY ../static /static

# Define the default command to run when the container starts
CMD ["./goexample"]