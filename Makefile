SOURCES = driver.go main.go paths.go

BIN = docker-volume-ploop
BINDIR = /usr/bin

SYSTEMD_SERVICE = etc/systemd/docker-volume-ploop.service
SYSTEMD_DIR = /usr/lib/systemd/system

CONFIG = etc/sysconfig/docker-volume-ploop
CONFIG_DIR = /etc/sysconfig

all: $(BIN)

$(BIN): $(SOURCES)
	go build -o $(BIN) $(SOURCES)

test:
	go test -v .

clean:
	rm -f $(BIN)

install: $(BIN) $(SYSTEMD_SERVICE) $(CONFIG)
	install -d $(DESTDIR)$(BINDIR)
	install -m 755 $(BIN) $(DESTDIR)$(BINDIR)
	install -d $(DESTDIR)$(SYSTEMD_DIR)
	install -m 644 $(SYSTEMD_SERVICE) $(DESTDIR)$(SYSTEMD_DIR)
	install -d $(DESTDIR)$(CONFIG_DIR)
	install -m 644 $(CONFIG) $(DESTDIR)$(CONFIG_DIR)

.PHONY: all test clean install
