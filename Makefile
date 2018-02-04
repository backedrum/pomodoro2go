BINARY=pomodoro2go

.DEFAULT_GOAL: $(BINARY)

$(BINARY):
	go build -o ${BINARY} *.go

format:
	go fmt $$(go list ./...) ; \
	cd - >/dev/null

clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi