FROM golang:latest
WORKDIR /
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN curl -sL https://deb.nodesource.com/setup_8.x -o nodesource_setup.sh
RUN chmod 755 nodesource_setup.sh
RUN bash nodesource_setup.sh
RUN apt-get install -y nodejs
RUN curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | apt-key add - && echo "deb https://dl.yarnpkg.com/debian/ stable main" | tee /etc/apt/sources.list.d/yarn.list
RUN apt-get update && apt-get install yarn
RUN apt-get install -y protobuf-compiler libprotobuf-dev
RUN go get -u github.com/golang/protobuf/protoc-gen-go
RUN go get -u github.com/ethereum/go-ethereum
RUN mkdir -p /go/src/github.com/singnet
WORKDIR /go/src/github.com/singnet
RUN git clone https://github.com/singnet/reputation-adapter.git
#COPY . /go/src/github.com/singnet/reputation-adapter
WORKDIR /go/src/github.com/singnet/reputation-adapter
RUN chmod +x scripts/install
RUN ./scripts/install
RUN chmod +x ./scripts/build
RUN scripts/build linux amd64
