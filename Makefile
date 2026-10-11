.PHONY: dev check
dev:
	hugo server -D --minify --renderToMemory --disableFastRender

check:
	python3 tests/hugo-smoke.py
