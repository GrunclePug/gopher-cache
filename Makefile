# gopher-cache - A thread-safe KV store for Go.
# See LICENSE file for copyright and license details.

include config.mk

all: gopher-cache

gopher-cache:
	mkdir -p bin
	${GO} build ${GOFLAGS} -ldflags "${LDFLAGS}" -o bin/gopher-cache ./cmd/gopher-cache

daemon:
	mkdir -p bin
	${GO} build ${GOFLAGS} -ldflags "${LDFLAGS}" -o bin/gopher-cached ./cmd/api

benchmark:
	mkdir -p bin
	${GO} build ${GOFLAGS} -ldflags "${LDFLAGS}" -o bin/benchmark ./cmd/benchmark

clean:
	rm -rf bin/ bench_db/ test_db/ data/

install: daemon
	mkdir -p ${DESTDIR}${PREFIX}/bin
	cp -f bin/gopher-cached ${DESTDIR}${PREFIX}/bin
	chmod 755 ${DESTDIR}${PREFIX}/bin/gopher-cached

uninstall:
	rm -f ${DESTDIR}${PREFIX}/bin/gopher-cached

.PHONY: all benchmark clean install uninstall
