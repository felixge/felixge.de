dev:
	jekyll serve -w

deploy:
	jekyll-s3

.PHONY: deploy dev
