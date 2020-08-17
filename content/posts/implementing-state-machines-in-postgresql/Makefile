DOT_FILES := $(shell find . -name '*.svg')

.PHONY: all
all: $(DOT_FILES)

%.svg: %.dot
	dot -Tsvg $^ > $@
