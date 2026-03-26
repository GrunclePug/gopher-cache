# gopher-cache - A thread-safe KV store for Go.
# See LICENSE file for copyright and license details.

include config.mk

SRC_CACHE = ./cmd/gopher-cache
SRC_BENCH = ./cmd/benchmark

all: gopher-cache

gopher-cache:
	mkdir -p bin
	${GO} build ${GOFLAGS} -ldflags "${LDFLAGS}" -o bin/gopher-cache ./cmd/gopher-cache

benchmark:
	mkdir -p bin
	${GO} build ${GOFLAGS} -ldflags "${LDFLAGS}" -o bin/benchmark ./cmd/benchmark

clean:
	rm -rf bin/ bench_db/ test_db/

install: all
	mkdir -p ${DESTDIR}${PREFIX}/bin
	cp -f bin/gopher-cache ${DESTDIR}${PREFIX}/bin
	chmod 755 ${DESTDIR}${PREFIX}/bin/gopher-cache

uninstall:
	rm -f ${DESTDIR}${PREFIX}/bin/gopher-cache

.PHONY: all benchmark clean install uninstall
