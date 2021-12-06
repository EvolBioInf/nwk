NAME = nwk

all : $(NAME)

$(NAME): $(NAME).go
	go build $(NAME).go
$(NAME).go: $(NAME).org
	awk -f scripts/preTangle.awk $(NAME).org | bash scripts/org2nw | notangle -R$(NAME).go | gofmt > $(NAME).go
test: $(NAME)_test.go $(NAME).go
	go test -v
$(NAME)_test.go: $(NAME)_test.org
	awk -f scripts/preTangle.awk $(NAME)_test.org | bash scripts/org2nw | notangle -R$(NAME)_test.go | gofmt > $(NAME)_test.go

.PHONY: doc
doc:
	make -C doc

clean:
	rm -f *.go
	make clean -C doc
