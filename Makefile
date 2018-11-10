install:
	go build && go install

release:
	docker build -t workmachine/workmachine .
	docker push workmachine/workmachine
