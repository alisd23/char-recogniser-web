FROM golang:1.8

WORKDIR /go/src/app
COPY . /go/src/char-recogniser-web
WORKDIR /go/src/char-recogniser-web

RUN curl https://glide.sh/get | sh
RUN glide install

RUN apt-get update
RUN curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | apt-key add -
RUN echo "deb https://dl.yarnpkg.com/debian/ stable main" | tee /etc/apt/sources.list.d/yarn.list
RUN apt-get install apt-transport-https
RUN curl -sL https://deb.nodesource.com/setup_8.x | bash -
RUN apt-get install -y nodejs
RUN apt-get install yarn

WORKDIR /go/src/char-recogniser-web/client

RUN yarn
RUN npm run build

WORKDIR /go/src/char-recogniser-web

ENV DB_HOST=mongo SERVER_HOST=0.0.0.0 PYTHON_HOST=char-recogniser-ml PYTHON_PORT=9001

ENTRYPOINT ["go", "run", "main.go", "start", "-a", "./client/build"]
