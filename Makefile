install:
	go build && go install

release:
	docker build -t workmachine/workmachine .
	docker push workmachine/workmachine

okr:
	gitokr ./OKR.json | dot -Tsvg > OKR.svg
	open -a "/Applications/Google Chrome.app" ./OKR.svg