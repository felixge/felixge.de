dev:
	jekyll serve --watch --drafts

deploy:
	jekyll-s3

.PHONY: deploy dev
