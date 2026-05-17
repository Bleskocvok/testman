.PHONY: all install clean

hello:
	echo "HELLO THESE NUTS!"

gitlist: gitlist.go
	go build -o gitlist gitlist.go

install: gitlist
	cp gitlist ~/bin/gitlist

clean:
	-rm gitlist
