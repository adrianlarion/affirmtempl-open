.PHONY: all run build_dev


all: build_dev run



build_dev:
	@templ generate
	@sass ./static/sass/materialize.scss ./static/css/materialize.css --quiet
	# @browserify internal/frontend/src/index.js -o static/js/app.js
	# browserify internal/frontend/src/index.js -o static/js/app.js -t [ babelify --presets [ @babel/preset-env @babel/preset-react ] --plugins [ @babel/plugin-transform-class-properties ] ]
	cd internal/frontend && npm run build
	@cp internal/frontend/app.js static/js/app.js
	go build -o ./tmp/main ./cmd/affirm/...


run:
	@./tmp/main










