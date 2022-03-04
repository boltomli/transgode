ifndef ver
	ver := int
endif

all:
	docker build -t cr/nm:$(ver) .

run:
	docker run --rm -it -p 8080:8080 cr/nm:$(ver)

local:
	go build && ./m
