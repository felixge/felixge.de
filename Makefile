dev:
	bundle exec jekyll serve --watch --drafts

deploy:
	jekyll build
	jekyll-s3

.PHONY: deploy dev
