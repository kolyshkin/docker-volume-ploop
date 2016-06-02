SOURCES = driver.go main.go paths.go
WORKPLACE = $(abspath Godeps/_workspace)

BIN = docker-volume-ploop
BINDIR = /usr/bin

SYSTEMD_FILES = etc/systemd/docker-volume-ploop.service \
		etc/systemd/docker-volume-ploop.socket
SYSTEMD_DIR = /usr/lib/systemd/system

CONFIG = etc/sysconfig/docker-volume-ploop
CONFIG_DIR = /etc/sysconfig

all: $(BIN)

$(BIN): $(SOURCES)
	GOPATH=$(WORKPLACE):$$GOPATH go build -o $(BIN) $(SOURCES)

test:
	go test -v .

clean:
	rm -f $(BIN)

install: $(BIN) $(SYSTEMD_SERVICE) $(CONFIG)
	install -d $(DESTDIR)$(BINDIR)
	install -m 755 $(BIN) $(DESTDIR)$(BINDIR)
	install -d $(DESTDIR)$(SYSTEMD_DIR)
	install -m 644 $(SYSTEMD_FILES) $(DESTDIR)$(SYSTEMD_DIR)
	install -d $(DESTDIR)$(CONFIG_DIR)
	install -m 644 $(CONFIG) $(DESTDIR)$(CONFIG_DIR)

.PHONY: all test clean install
