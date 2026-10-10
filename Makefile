.PHONY: dev check
dev:
	hugo server -D --minify

check:
	python3 tests/hugo-smoke.py
