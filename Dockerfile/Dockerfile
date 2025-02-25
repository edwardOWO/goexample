FROM alpine:latest

RUN apk add --no-cache tzdata openssl curl
ENV TZ=Asia/Taipei
RUN curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | sh

RUN curl -LO https://dl.k8s.io/release/v1.24.0/bin/linux/amd64/kubectl
RUN chmod +x ./kubectl
RUN mv ./kubectl /usr/local/bin/kubectl
ENV PATH="/usr/local/bin:${PATH}"


RUN mkdir -p /opt/log
COPY ../goexample /goexample
COPY ../static /static
COPY ../template /template

# Define the default command to run when the container starts
CMD ["./goexample"]