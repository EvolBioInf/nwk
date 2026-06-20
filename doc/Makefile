date = $(shell git log | grep -m 1 Date | sed -r 's/Date: +[A-Z][a-z]+ ([A-Z][a-z]+) ([0-9]+) [^ ]+ ([0-9]+) .+/\2_\1_\3/')
version = $(shell git describe)

NAME = nwk
all: $(NAME)Doc.pdf
$(NAME)Doc.pdf: $(NAME)Doc.tex $(NAME).tex $(NAME)_test.tex
	echo $(date) | tr '_' ' ' > date.txt
	echo $(version) | tr '-' ' ' | awk '{printf "%s", $$1; if ($$2) printf "-%s", $$2; printf "\n"}' > version.txt
	latex $(NAME)Doc
	bibtex $(NAME)Doc
	latex $(NAME)Doc
	latex $(NAME)Doc
	latex $(NAME)Doc
	dvipdf -dALLOWPSTRANSPARENCY $(NAME)Doc
$(NAME).tex: ../$(NAME).org
	bash ../scripts/org2nw ../$(NAME).org       | awk -f ../scripts/preWeave.awk | noweave -n -x > $(NAME).tex
$(NAME)_test.tex: ../$(NAME)_test.org
	bash ../scripts/org2nw ../$(NAME)_test.org  | awk -f ../scripts/preWeave.awk | noweave -n -x | sed 's/_test/\\_test/' > $(NAME)_test.tex
publish: $(NAME)Doc.pdf
	if [ -d ~/WinHome/ownCloud/ ]; then \
		cp $(NAME)Doc.pdf ~/WinHome/ownCloud/docs/; \
	fi
clean:
	rm -f $(NAME).tex $(NAME)_test.tex *.pdf *.aux *.bbl *.blg *.dvi *.log *.toc date.txt version.txt
