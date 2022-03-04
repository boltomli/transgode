ifndef ver
	ver := int
endif

all:
	docker build -t cr/nm:$(ver) .

run:
	docker run --rm -it cr/nm:$(ver)
