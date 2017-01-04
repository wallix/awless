FROM golang

RUN apt-get update; apt-get install -y libgit2-dev cmake git

ADD . /go/src/github.com/wallix/awless

RUN rm -rf /go/src/github.com/wallix/awless/vendor/github.com/libgit2
RUN git clone https://github.com/libgit2/git2go.git /go/src/github.com/libgit2/git2go

WORKDIR /go/src/github.com/libgit2/git2go

RUN git checkout next
RUN git submodule update --init
RUN make install

WORKDIR /go/src/github.com/wallix/awless

VOLUME "/go/bin"

CMD ["go", "install"]