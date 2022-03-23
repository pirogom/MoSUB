clean:
	go clean

win:
	GOOS=windows GOARCH=386 go build -ldflags '-s -w' -o MOSUB.exe
	
darwin:
	GOOS=darwin GOARCH=amd64 go build -ldflags '-s -w' -o MOSUB_DARWIN
m1:
	mv MOSUB.syso _MOSUB.syso
	GOOS=darwin GOARCH=arm64 go build -ldflags '-s -w' -o MOSUB_DARWIN_M1
	mv _MOSUB.syso MOSUB.syso
rsrc:
	rsrc -manifest MOSUB.manifest -ico=mop.ico -o MOSUB.syso
all:
	make win darwin m1